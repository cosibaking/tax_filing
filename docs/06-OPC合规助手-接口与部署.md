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
