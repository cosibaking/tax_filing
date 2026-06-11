# OCR.space 模块使用文档

## 1. 快速开始

### 1.1 获取 API Key

在 [OCR.space 免费注册页](https://ocr.space/ocrapi/freekey) 申请 API Key。测试可用官方示例 Key `helloworld`（有配额限制，仅供调试）。

### 1.2 创建配置文件

```bash
cd server
cp internal/library/ocrspace/config.yaml.example manifest/config/ocrspace.yaml
```

编辑 `manifest/config/ocrspace.yaml`：

```yaml
ocrspace:
  apiKey: "你的-api-key"
  defaultLanguage: "chs"
  defaultEngine: 2
  timeout: 60s
```

> 含真实 Key 的配置文件请勿提交到 Git。

### 1.3 最简调用

```go
package main

import (
    "context"
    "fmt"
    "log"

    "xygo/internal/library/ocrspace"
)

func main() {
    client, err := ocrspace.NewClientFromFile("manifest/config/ocrspace.yaml")
    if err != nil {
        log.Fatal(err)
    }

    text, err := client.ExtractTextFromFile(context.Background(), "invoice.pdf", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(text)
}
```

---

## 2. 配置说明

### 2.1 配置文件字段

| 字段 | 必填 | 默认值 | 说明 |
|------|------|--------|------|
| `apiKey` | 是 | — | OCR.space API Key |
| `endpoint` | 否 | `https://api.ocr.space/parse/image` | POST 端点 |
| `endpointGet` | 否 | `https://api.ocr.space/parse/imageurl` | GET 端点 |
| `timeout` | 否 | `60s` | HTTP 请求超时 |
| `defaultLanguage` | 否 | `eng` | 默认语言代码 |
| `defaultEngine` | 否 | `2` | 默认 OCR 引擎 |

### 2.2 环境变量

| 变量 | 作用 |
|------|------|
| `OCRSPACE_API_KEY` | 覆盖配置文件中的 `apiKey`（生产推荐） |
| `OCRSPACE_CONFIG` | 指定配置文件路径 |

示例（PowerShell）：

```powershell
$env:OCRSPACE_API_KEY = "your-key"
$env:OCRSPACE_CONFIG = "D:\configs\ocrspace.yaml"
```

### 2.3 代码中直接构造配置

```go
cfg := ocrspace.DefaultConfig()
cfg.APIKey = "your-key"
cfg.DefaultLanguage = ocrspace.LanguageChineseSimplified
cfg.DefaultEngine = ocrspace.Engine2

client, err := ocrspace.NewClient(cfg)
```

---

## 3. API 方法一览

### 3.1 完整响应方法

返回 `*ParseResponse`，包含原始结构化数据：

| 方法 | 输入 | 说明 |
|------|------|------|
| `ParseURL(ctx, url, opts)` | 远程 URL | POST，推荐 |
| `ParseURLGet(ctx, url, opts)` | 远程 URL | GET，仅 URL |
| `ParseFile(ctx, path, opts)` | 本地文件路径 | 自动读取并上传 |
| `ParseBytes(ctx, data, filename, opts)` | `[]byte` + 文件名 | 内存上传 |
| `ParseReader(ctx, r, filename, opts)` | `io.Reader` + 文件名 | 流式上传 |
| `ParseBase64(ctx, base64, opts)` | Base64 字符串 | 需 data URI 前缀 |

### 3.2 便捷文本方法

识别成功返回合并文本，失败返回 `error`：

| 方法 | 对应完整方法 |
|------|-------------|
| `ExtractTextFromURL` | `ParseURL` |
| `ExtractTextFromFile` | `ParseFile` |
| `ExtractTextFromBytes` | `ParseBytes` |

---

## 4. 使用示例

### 4.1 识别远程图片

```go
opts := &ocrspace.ParseOptions{
    Language: ocrspace.LanguageChineseSimplified,
    Scale:    true,
}

resp, err := client.ParseURL(ctx, "https://example.com/receipt.jpg", opts)
if err != nil {
    return err
}
if err = resp.Err(); err != nil {
    return err
}
fmt.Println(resp.CombinedText())
```

### 4.2 识别上传的字节内容

适用于 HTTP 接口已读取 multipart 文件到内存的场景：

```go
text, err := client.ExtractTextFromBytes(ctx, fileBytes, "upload.png", &ocrspace.ParseOptions{
    IsTable: true, // 表格/收据
})
```

### 4.3 识别 Base64 图片

Base64 必须带 data URI 前缀：

```go
b64 := "data:image/jpeg;base64,/9j/4AAQSkZJRg..."
resp, err := client.ParseBase64(ctx, b64, nil)
```

若仅有裸 Base64 数据，需手动拼接前缀：

```go
b64 := "data:image/png;base64," + rawBase64String
```

### 4.4 获取文字坐标（Overlay）

```go
resp, err := client.ParseFile(ctx, "scan.png", &ocrspace.ParseOptions{
    IsOverlayRequired: true,
})
if err != nil || resp.Err() != nil {
    return err
}

for _, page := range resp.ParsedResults {
    if page.TextOverlay == nil {
        continue
    }
    for _, line := range page.TextOverlay.Lines {
        for _, word := range line.Words {
            fmt.Printf("%s @ (%d,%d)\n", word.WordText, word.Left, word.Top)
        }
    }
}
```

### 4.5 生成可搜索 PDF

```go
resp, err := client.ParseFile(ctx, "scan.pdf", &ocrspace.ParseOptions{
    IsCreateSearchablePdf:        true,
    IsSearchablePdfHideTextLayer: true, // 隐藏文字层
})
if err != nil || resp.Err() != nil {
    return err
}
if resp.SearchablePDFURL != nil {
    fmt.Println("下载链接（1 小时有效）:", *resp.SearchablePDFURL)
}
```

> Engine 3 暂不支持可搜索 PDF，请使用 Engine 1 或 2。

### 4.6 指定 OCR 引擎

```go
// 高精度识别（手写、复杂表格）
opts := &ocrspace.ParseOptions{
    OCREngine: ocrspace.Engine3,
    Language:  ocrspace.LanguageAuto,
}

text, err := client.ExtractTextFromFile(ctx, "handwriting.jpg", opts)
```

---

## 5. ParseOptions 参数参考

| 字段 | 类型 | 默认 | 说明 |
|------|------|------|------|
| `Language` | `Language` | 配置默认 | 语言代码，如 `eng`、`chs`、`auto` |
| `IsOverlayRequired` | `bool` | `false` | 返回文字坐标 |
| `Filetype` | `FileType` | 自动检测 | 强制指定 `PDF`/`PNG`/`JPG` 等 |
| `DetectOrientation` | `bool` | `false` | 自动旋转校正 |
| `IsCreateSearchablePdf` | `bool` | `false` | 生成可搜索 PDF |
| `IsSearchablePdfHideTextLayer` | `bool` | `false` | 隐藏 PDF 文字层 |
| `Scale` | `bool` | `false` | 内部放大，低清扫描建议开启 |
| `IsTable` | `bool` | `false` | 表格/收据模式，按行返回 |
| `OCREngine` | `OCREngine` | 配置默认 | `1` / `2` / `3` |

### 常用语言代码

| 代码 | 语言 |
|------|------|
| `eng` | 英语 |
| `chs` | 简体中文 |
| `cht` | 繁体中文 |
| `jpn` | 日语 |
| `kor` | 韩语 |
| `auto` | 自动检测（Engine 2/3） |

完整列表见 [OCR.space 文档](https://ocr.space/ocrapi)。

---

## 6. 响应处理

### 6.1 判断成功与否

```go
resp, err := client.ParseURL(ctx, imageURL, nil)
if err != nil {
    // 网络/HTTP/JSON 错误
    return err
}
if !resp.IsSuccess() {
    return resp.Err()
}
```

`IsSuccess()` 在 `OCRExitCode` 为 1（全成功）或 2（部分成功）时返回 `true`。

### 6.2 获取文本

```go
// 所有页合并，页间空一行
text := resp.CombinedText()

// 逐页访问
for i, page := range resp.ParsedResults {
    fmt.Printf("第 %d 页: %s\n", i+1, page.ParsedText)
}
```

### 6.3 响应字段

| 字段 | 说明 |
|------|------|
| `ParsedResults` | 每页/每图结果数组 |
| `OCRExitCode` | 整体退出码 |
| `IsErroredOnProcessing` | 是否处理出错 |
| `ErrorMessage` / `ErrorDetails` | 错误信息 |
| `SearchablePDFURL` | 可搜索 PDF 下载链接 |
| `ProcessingTimeInMilliseconds` | 处理耗时 |

---

## 7. 常见问题

### Q: 报错 `apiKey is required`

配置文件中未设置 `apiKey`，且环境变量 `OCRSPACE_API_KEY` 也为空。请检查配置文件路径是否正确（默认 `manifest/config/ocrspace.yaml`，相对于**进程工作目录**）。

### Q: 报错 `Not a valid base64 image`

Base64 字符串缺少 `data:image/jpeg;base64,` 等前缀，或末尾有多余换行。

### Q: 中文识别效果差

尝试：

1. 设置 `Language: ocrspace.LanguageChineseSimplified`
2. 开启 `Scale: true`
3. 换用 `Engine2` 或 `Engine3`

### Q: PDF 多页只识别部分

检查 Free 计划 PDF 页数限制（通常 3 页）。`OCRExitCode = 2` 表示部分成功，可逐页检查 `ParsedResults[i].ErrorMessage`。

### Q: 文件太大上传失败

Free 计划单文件约 1 MB 限制。可升级 PRO 计划，或在业务层压缩/裁剪后再调用。

### Q: 如何在 GoFrame 项目中使用

模块不依赖 GoFrame，可在任意 logic 函数中：

```go
import "xygo/internal/library/ocrspace"

func (s *sSomeLogic) Recognize(ctx context.Context, fileBytes []byte, filename string) (string, error) {
    client, err := ocrspace.NewClientFromFile("")
    if err != nil {
        return "", err
    }
    return client.ExtractTextFromBytes(ctx, fileBytes, filename, nil)
}
```

---

## 8. 运行测试

```bash
cd server
go test ./internal/library/ocrspace/... -v
```

测试不依赖外网和真实 API Key。

---

## 9. 文件清单

| 文件 | 说明 |
|------|------|
| `client.go` | 客户端实现 |
| `config.go` | 配置加载 |
| `types.go` | 类型与响应辅助 |
| `config.yaml.example` | 配置模板 |
| `DESIGN.md` | 设计文档 |
| `USAGE.md` | 本文档 |
