# OPC 公司经营合规助手：验收清单

## 自动验证

```bash
cd server
GOTOOLCHAIN=local go test ./...
GOTOOLCHAIN=local go test -race ./internal/logic/compliance/risk ./internal/logic/compliance/task
GOTOOLCHAIN=local go build ./...

cd ../web
pnpm exec eslint src/api/frontend/compliance/assistant.ts src/api/backend/compliance/assistant.ts \
  src/views/frontend/compliance/{profile,calendar,task-detail,documents,risks,reports,tickets}.vue \
  src/views/backend/compliance/assistant/index.vue
pnpm build
```

仓库现有 `go vet ./...` 会被历史后台控制器中的未命名结构体字段告警拦截，相关文件不属于本功能改动；上线前可单独安排技术债修复。

## 业务验收

- [ ] 开关关闭时，会员接口拒绝访问且旧模块正常。
- [ ] 保存企业画像后生成新版本；内容不变时不重复建版本。
- [ ] 相同规则、期间和触发键不会重复生成任务。
- [ ] 任务只能按定义状态机流转，跨会员访问被拒绝。
- [ ] 资料原附件保留；分类结果必须由用户确认后使用。
- [ ] 十条风险规则分别覆盖命中、边界和数据缺失情况。
- [ ] 同期间多类风险可共存；处置记录按会员隔离。
- [ ] 数据不足的报告明确显示“数据不足”。
- [ ] AI 不可用时仍能生成确定性报告，敏感号码不会进入提示文本。
- [ ] 已发布报告不可修改，只能生成新版本。
- [ ] 申报/代理记账工单未指定服务机构时禁止提交。
- [ ] 提醒失败三次后进入 `manual_attention`。
- [ ] 管理员可试算规则，创建人与审核人相同时禁止发布。

## 真实环境冒烟路径

```text
开启灰度开关
→ 登录会员并确认 OPC 主体
→ /user/compliance/profile 保存画像
→ /user/compliance/calendar 查看任务
→ /user/compliance/documents 登记并确认资料
→ /user/compliance/risks 查看和处置疑点
→ /user/compliance/reports 生成并发布报告
→ /user/compliance/tickets 提交人工复核
→ 管理端 compliance/assistant 查看全链路记录
```

测试数据使用专用企业与附件；不得在生产企业上构造虚假流水或申报事实。
