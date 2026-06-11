# OCR.space 模块设计文档

## 1. 背景与目标

### 1.1 背景

项目需要对接第三方 OCR 能力，用于识别图片、PDF 中的文字内容。选用 [OCR.space API](https://ocr.space/ocrapi) 作为云服务提供商，其支持：

- 图片（PNG、JPG、GIF、TIF、BMP）与 PDF 识别
- 多语言（含中文简体/繁体）
- 三种 OCR 引擎（速度/精度权衡）
- 文字坐标 overlay、可搜索 PDF 等高级能力

### 1.2 设计目标

| 目标 | 说明 |
|------|------|
| 独立模块 | 不依赖现有业务逻辑，不修改项目架构 |
| 接口清晰 | 对外暴露少量、语义明确的方法 |
| 配置外置 | API Key 及端点通过 YAML / 环境变量配置 |
| 可测试 | 核心逻辑可通过 `httptest` 模拟，无需真实 API Key |
| 可扩展 | 后续接入业务层时仅需 import 并调用 |

### 1.3 非目标（当前版本不做）

- 不接入 controller / logic / service 层
- 不写入主配置文件 `config.yaml`
- 不做 OCR 结果持久化、队列、重试策略
- 不封装 overlay 可视化绘制

---

## 2. 模块位置与依赖

```
server/internal/library/ocrspace/
├── client.go              # HTTP 客户端与 API 调用
├── client_test.go         # 单元测试
├── config.go              # 配置加载与校验
├── config.yaml.example    # 配置示例
├── types.go               # 类型定义与响应辅助方法
├── DESIGN.md              # 本文档
└── USAGE.md               # 使用文档
```

**外部依赖：**

- 标准库：`net/http`、`mime/multipart`、`encoding/json`、`context`
- 第三方：`gopkg.in/yaml.v3`（项目已有）

**与现有模块关系：**

- `platformocr`：本地 Tesseract / mock，用于合规流水截图，**与本模块无关**
- 本模块为全新独立包，命名 `ocrspace`，避免与现有 OCR 逻辑混淆

---

## 3. 架构设计

### 3.1 分层结构

```
┌─────────────────────────────────────────┐
│  业务层（未来接入，当前不存在）            │
│  logic / service / controller           │
└──────────────────┬──────────────────────┘
                   │ import
┌──────────────────▼──────────────────────┐
│  ocrspace.Client（对外入口）             │
│  - ParseURL / ParseFile / ParseBase64   │
│  - ExtractTextFrom*（便捷方法）          │
└──────────────────┬──────────────────────┘
                   │
        ┌──────────┴──────────┐
        ▼                     ▼
┌───────────────┐     ┌───────────────────┐
│  Config       │     │  types.ParseResponse │
│  加载/校验     │     │  结果解析/辅助方法    │
└───────────────┘     └───────────────────┘
                   │
                   ▼
        ┌──────────────────────┐
        │  OCR.space HTTP API   │
        │  POST /parse/image    │
        │  GET  /parse/imageurl │
        └──────────────────────┘
```

### 3.2 核心类型

| 类型 | 职责 |
|------|------|
| `Config` | 存储 API Key、端点、超时、默认语言/引擎 |
| `Client` | 持有配置与 `http.Client`，执行 OCR 请求 |
| `ParseOptions` | 单次请求的可选参数（语言、引擎、overlay 等） |
| `ParseResponse` | OCR.space 原始 JSON 的结构化映射 |
| `ParsedResult` | 单页/单图识别结果 |
| `TextOverlay` | 文字坐标信息（可选返回） |

### 3.3 配置加载优先级

```
DefaultConfig() 默认值
    ↓ 合并
YAML 配置文件（ocrspace.yaml）
    ↓ 覆盖
环境变量 OCRSPACE_API_KEY
    ↓ 校验
Config.Validate()
```

配置文件查找顺序（`LoadConfig("")`）：

1. 函数参数 `path`（若非空）
2. 环境变量 `OCRSPACE_CONFIG`
3. 默认路径 `manifest/config/ocrspace.yaml`

---

## 4. API 映射设计

### 4.1 输入方式与 Client 方法对应

OCR.space 支持三种 POST 输入方式，模块分别封装为：

| OCR.space 参数 | Client 方法 | 说明 |
|----------------|-------------|------|
| `url` | `ParseURL` | 远程图片/PDF URL |
| `file` | `ParseFile` / `ParseBytes` / `ParseReader` | 本地或内存文件上传 |
| `base64Image` | `ParseBase64` | Base64 字符串（需 data URI 前缀） |
| GET `url` | `ParseURLGet` | 仅 URL，参数较少，API Key 在 query 中 |

### 4.2 请求构造

- **认证**：API Key 通过 HTTP Header `apikey` 传递（POST）
- **Content-Type**：`multipart/form-data`
- **可选参数**：由 `ParseOptions` 映射为 form field，布尔值为 `"true"` 字符串
- **默认值**：`Language`、`OCREngine` 未指定时使用 `Config` 中的默认值

### 4.3 响应处理

```
HTTP 响应
  ├─ StatusCode != 2xx → 返回 transport 层 error
  └─ JSON → ParseResponse
        ├─ OCRExitCode 1/2 → IsSuccess() = true
        ├─ CombinedText() → 合并多页文本
        └─ Err() → 业务层 error
```

**OCRExitCode 语义（来自 OCR.space 文档）：**

| 值 | 含义 |
|----|------|
| 1 | 全部解析成功 |
| 2 | 部分成功（如多页 PDF 部分页失败） |
| 3 | 全部失败 |
| 4 | 解析过程发生致命错误 |

---

## 5. OCR 引擎与参数选型

### 5.1 引擎对比

| 引擎 | 常量 | 特点 | 适用场景 |
|------|------|------|----------|
| Engine 1 | `Engine1` | 最快，支持亚洲语言 | 清晰扫描件、大批量 |
| Engine 2 | `Engine2` | 速度/质量均衡（**默认**） | 通用场景，首选 |
| Engine 3 | `Engine3` | 精度最高，支持手写/表格 | 复杂版式、低质量图片 |

### 5.2 常用 ParseOptions

| 字段 | 用途 |
|------|------|
| `Language` | 识别语言，`chs` 简体、`cht` 繁体、`auto` 自动检测（Engine 2/3） |
| `Scale` | 内部放大，低分辨率扫描件建议开启 |
| `IsTable` | 表格/收据/发票，按行返回 |
| `DetectOrientation` | 自动旋转校正 |
| `IsOverlayRequired` | 返回文字坐标 |
| `IsCreateSearchablePdf` | 生成可搜索 PDF（Engine 3 暂不支持） |
| `Filetype` | 覆盖 Content-Type 自动检测 |

---

## 6. 错误处理策略

模块区分两类错误：

1. **Transport 错误**：网络超时、HTTP 非 2xx、JSON 解析失败 → 由 `Parse*` 方法直接返回 `error`
2. **OCR 业务错误**：HTTP 200 但 `OCRExitCode` 表示失败 → 由 `ParseResponse.Err()` 返回

便捷方法 `ExtractTextFrom*` 会在两类错误时都返回 `error`，成功时仅返回合并文本。

错误信息统一前缀 `ocrspace:`，便于日志过滤。

---

## 7. 安全与运维

### 7.1 凭证管理

- API Key **不得**硬编码在源码中
- 生产环境推荐：配置文件 + 环境变量覆盖，或将 Key 注入容器环境
- `ocrspace.yaml` 应加入 `.gitignore`（若含真实 Key）

### 7.2 限流与配额

Free 计划限制（参考 OCR.space 文档，可能变更）：

- 约 25,000 次/月（Engine 1/2）
- Engine 3 单独配额
- 单文件大小、PDF 页数有限制

模块本身不做限流，业务接入时需自行控制调用频率。

### 7.3 超时

默认 `60s`，可通过配置 `timeout: 90s` 调整。大 PDF 或 Engine 3 可适当增大。

---

## 8. 测试设计

| 测试 | 覆盖点 |
|------|--------|
| `TestLoadConfigFromBytes` | YAML 解析、默认值合并 |
| `TestLoadConfigFromBytesMissingKey` | 缺少 apiKey 校验 |
| `TestParseResponseHelpers` | `IsSuccess` / `CombinedText` / `Err` |
| `TestClientParseURL` | POST multipart + Header apikey |
| `TestClientParseFile` | 文件上传 |
| `TestClientParseBase64` | Base64 提交 |
| `TestEnvAPIKeyOverride` | 环境变量覆盖 |

测试使用 `httptest.NewServer` 模拟 OCR.space，无需外网与真实 Key。

运行：

```bash
cd server
go test ./internal/library/ocrspace/... -v
```

---

## 9. 未来接入建议

当需要接入业务层时，建议：

1. 在 `internal/logic/` 新建独立 logic 包，注入 `*ocrspace.Client`
2. 通过 GoFrame 的 `g.Cfg()` 或继续沿用独立 `ocrspace.yaml`，二选一即可
3. 对上传文件先做大小/格式校验，再调用 `ParseBytes`
4. 对 Engine 3 表格结果，可直接使用 `ParsedText`（含 Markdown 格式）
5. 如需重试，在 logic 层实现（如仅对 timeout 重试 1 次）

**不建议**将 OCR 逻辑直接写入现有 `platformocr` 包，保持职责分离。

---

## 10. 参考

- OCR.space API 文档：https://ocr.space/ocrapi
- 免费 API Key 注册：https://ocr.space/ocrapi/freekey
- API 状态页：https://status.ocr.space
