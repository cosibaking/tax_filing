# OPC 公司经营合规助手：接口与部署

## 功能开关

生产配置默认关闭，确认数据库迁移、后台规则和服务机构配置后再逐项开启：

```yaml
complianceAssistant:
  enabled: false
  taskCronEnabled: false
  riskCronEnabled: false
  reportCronEnabled: false
  aiNarrationEnabled: false
```

关闭 `enabled` 后，会员端助手接口返回“尚未开放”；数据与历史记录不会删除。AI 关闭或调用失败时使用确定性模板报告。

## 会员接口

| 方法 | 路径 | 用途 |
|---|---|---|
| GET/PUT | `/compliance/profile` | 查询、保存企业画像 |
| GET | `/compliance/tasks` | 查询合规任务 |
| GET | `/compliance/tasks/{id}` | 任务详情 |
| POST | `/compliance/tasks/{id}/action` | 任务流转 |
| GET/POST | `/compliance/documents` | 资料列表、登记附件 |
| PUT | `/compliance/documents/{id}/confirm` | 确认资料分类 |
| GET | `/compliance/risks` | 风险列表 |
| POST | `/compliance/risks/scan` | 按确认事实执行扫描 |
| POST | `/compliance/risks/{id}/action` | 确认、排除或解决风险 |
| GET/POST | `/compliance/reports` | 报告列表、生成草稿 |
| POST | `/compliance/reports/{id}/publish` | 发布不可变报告版本 |
| GET | `/compliance/reports/{id}/pdf` | 导出月度体检 PDF |
| GET/POST | `/compliance/tickets` | 工单列表、提交工单 |
| POST | `/compliance/tickets/{id}/action` | 工单流转 |

统一响应由项目现有 HTTP 封装处理；业务数据位于 `data`。所有资源按当前会员隔离。

### 示例：执行风险扫描

```json
{
  "periodKey": "2026-07",
  "facts": {
    "revenue": 35600,
    "salesDocumentAmount": 21000,
    "zeroFiled": false,
    "personalTransferCount": 3,
    "personalTransferAmount": 1268
  }
}
```

### 示例：提交人工服务工单

```json
{
  "ticketType": "compliance_review",
  "title": "第二季度资料人工复核",
  "description": "请核对银行流水和销项记录",
  "riskEventId": 12
}
```

`tax_filing`、`bookkeeping` 类型必须填写 `providerName`，以明确有资质服务机构。

## 月度体检报告接口

报告仅能由当前登录会员访问。会员身份和 OPC 主体均由服务端根据登录上下文确定，客户端不得指定或覆盖。创建会保存服务端在生成时组装的结构化快照，发布后该版本不可修改；需要更新时应生成新版本。

### 生成报告草稿

`POST /compliance/reports`

请求体：

```json
{
  "periodKey": "2026-07"
}
```

`periodKey` 是唯一业务入参，必须为真实有效的 `YYYY-MM`。服务端使用当前会员与 OPC 主体，从可信的经营资料、流水、资料完整度和合规风险记录组装报告。

为兼容旧客户端，服务端仍可从旧的 `data.periodKey` 中读取期间；但旧请求内携带的 `data.statistics`、`data.completeness` 和 `data.risks` 一律忽略，不会进入报告快照。这些数据必须由服务端按当前会员范围查询，避免客户端伪造风险数量、资料完整度或他人经营数据。

如果指定期间的可信数据不足，服务端仍可生成草稿，但必须在总体结论、分类状态和资料提示中明确标记“数据不足”，不得因客户端补传快照而得出确定性结论。

成功响应的 `data.report` 与列表项结构一致，初始 `status` 为 `draft`。

### 查询报告列表

`GET /compliance/reports?periodKey=2026-07`

`periodKey` 可选；不传时返回当前会员的全部报告，按期间和版本倒序排列。结构化报告示例：

```json
{
  "code": 0,
  "message": "",
  "data": {
    "list": [
      {
        "id": 27,
        "opcEntityId": 6,
        "memberId": 18,
        "periodKey": "2026-07",
        "version": 1,
        "status": "draft",
        "content": "2026-07 月度经营合规体检\n……",
        "aiModel": "",
        "knowledgeVersion": "",
        "structuredReport": {
          "schemaVersion": 1,
          "summary": {
            "conclusion": "urgent",
            "completenessRate": 81.5,
            "highCount": 1,
            "mediumCount": 0,
            "lowCount": 0,
            "dataNotice": "资料尚未完全齐备，请根据缺失清单及时补充。"
          },
          "categories": [
            {
              "code": "tax",
              "name": "税务与申报",
              "status": "urgent",
              "summary": "识别到 1 项待处理风险，请按异常清单核实。",
              "checks": [
                {
                  "code": "tax-overview",
                  "name": "税务与申报检查",
                  "status": "urgent",
                  "message": "识别到 1 项待处理风险，请按异常清单核实。"
                }
              ]
            }
          ],
          "anomalies": [
            {
              "code": "unmatched-revenue",
              "categoryCode": "tax",
              "title": "部分收入未匹配发票",
              "severity": "high",
              "facts": "已确认回款 35,600 元，其中 14,600 元尚未匹配销项发票。",
              "basis": "基于本期已确认的银行流水与销项资料比对。",
              "impact": "可能影响申报数据的完整性。",
              "recommendation": "核对合同、回款和开票记录，交由专业人员复核。",
              "requiredMaterials": ["银行流水", "业务合同", "销项发票"],
              "dueDate": "2026-08-10",
              "requiresManualReview": true,
              "ruleVersion": "tax-check-v3"
            }
          ]
        },
        "legacy": false,
        "publishedAt": 0,
        "createdAt": 1783958400
      }
    ]
  }
}
```

`summary.conclusion` 可为 `normal`、`attention`、`urgent`；分类 `status` 还可为 `insufficient`；异常 `severity` 可为 `low`、`medium`、`high`。异常中的“事实、依据、影响、建议、所需资料、截止日期、是否人工复核、规则版本”用于给出可追溯的详细说明。

历史报告没有 `structured_json` 时，返回 `legacy: true`，且不返回 `structuredReport`；前端和 PDF 导出必须降级使用 `content`，不得将历史报告显示为空白。结构化快照损坏时，列表会降级为历史内容并记录警告日志。

### 发布报告

`POST /compliance/reports/{id}/publish`

| 参数 | 位置 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| `id` | path | uint64 | 是 | 报告 ID，必须大于 0 |

成功后 `data.report.status` 为 `published`，并返回 `publishedAt`。重复发布同一版本返回“已发布报告不可修改，请创建新版本”。

### 导出 PDF

`GET /compliance/reports/{id}/pdf`

| 参数 | 位置 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| `id` | path | uint64 | 是 | 报告 ID，必须大于 0 |

草稿和已发布版本都可导出。成功响应不使用 JSON 信封：

```http
HTTP/1.1 200 OK
Content-Type: application/pdf
Content-Disposition: attachment; filename="monthly-checkup-2026-07-v1.pdf"
Content-Length: 12345
Cache-Control: private, no-store

%PDF-...
```

PDF 包含企业、期间、版本、状态、总体结论、六类检查、完整异常说明和免责声明。历史报告降级输出 `content` 摘要。

越权或不存在的 ID 不泄露资源是否存在，返回统一错误信封：

```json
{
  "code": -1,
  "message": "报告不存在或无权访问",
  "traceId": "f7f9c265d6104c72"
}
```

参数错误示例：

```json
{
  "code": -1,
  "message": "请指定报告",
  "traceId": "f7f9c265d6104c72"
}
```

PDF 加载字体或生成失败时也返回 JSON 错误信封，客户端必须在 `responseType: blob` 下检查 `Content-Type`，不要把 JSON 错误 Blob 作为 PDF 保存。

> 审慎说明：体检报告是基于已确认资料和规则的自动整理结果，不是税务申报结论，也不替代会计、税务或法律专业意见。高风险、数据不足、规则版本变化及重要结论必须结合原始凭证人工复核。

## 管理接口

| 方法 | 路径 | 用途 |
|---|---|---|
| POST | `/admin/compliance/rule-versions/{id}/simulate` | 规则试算 |
| POST | `/admin/compliance/rule-versions/{id}/transition` | 提审、发布、退役 |
| GET | `/admin/compliance/assistant/workspace` | rules/tasks/documents/risks/reports/tickets 工作台 |

规则创建人与发布审核人必须不同。发布版本不可直接修改，应创建新版本。

## 定时任务

在后台定时任务中注册并设置 Cron：

- `compliance_assistant_reminder`：建议每 10 分钟；处理 7/3/1 天提醒，失败三次转人工关注。
- `compliance_assistant_risk_scan`：建议每天凌晨；当前仅在开关开启后进入就绪态，扫描需有已确认期间事实。
- `compliance_assistant_monthly_report`：建议每月 2 日；当前仅在开关开启后进入就绪态，报告需有已确认数据。

后两项不在数据不足时自动下结论，避免将缺失资料误判为税务事实。

## 部署顺序

1. 备份数据库并部署代码，保持所有助手开关关闭。
2. 执行 `server/cmd_tools/migrate/1.5.6_compliance_assistant_core.mysql.sql`。
3. 启动后端并确认健康检查、管理端工作台和数据库表。
4. 配置已审核规则、定时任务及服务机构。
5. 先对测试企业开启 `enabled`，验证画像、任务、资料、风险、报告、工单闭环。
6. 再逐项开启任务提醒；风险与报告批处理在确认数据来源后开启。

## 回滚

将 `complianceAssistant.enabled` 和所有子开关设为 `false` 并重启服务。不要删除、清空新表；旧业务模块不依赖这些表，可继续运行。
