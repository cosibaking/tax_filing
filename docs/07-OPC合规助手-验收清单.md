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

## 1.5.7 结构化报告迁移

本项目对应的第三个合规助手增量迁移在仓库中命名为 `1.5.7_monthly_report_structured`（即部署清单中的 migration 000003）。它只在 `xy_compliance_monthly_report` 增加可空 `structured_json`，不改写历史报告。

### 执行前

```bash
cd server
# 使用与目标环境一致的 .env/配置，先只查看状态
go run tools.go migrate status

# MySQL 8.0+：确认目标表和字段状态
mysql --defaults-extra-file=/path/to/client.cnf -e \
  "SHOW CREATE TABLE xy_compliance_monthly_report\\G; SHOW COLUMNS FROM xy_compliance_monthly_report LIKE 'structured_json';"
```

生产执行前必须备份表结构和当前发布包，并在维护窗口评估 `ALTER TABLE` 持锁时间。密码使用配置或 `--defaults-extra-file`，禁止写入文档和命令历史。

### 执行与验证

```bash
cd server
go run tools.go migrate up
go run tools.go migrate status

mysql --defaults-extra-file=/path/to/client.cnf -e \
  "SHOW COLUMNS FROM xy_compliance_monthly_report LIKE 'structured_json';"
```

也可由发布系统单独执行 `cmd_tools/migrate/1.5.7_monthly_report_structured.mysql.sql`。MySQL 脚本是幂等的：字段已存在时仅执行 `SELECT 1`。PostgreSQL 环境使用同名 `.pgsql.sql`。

### 回滚

默认仅回滚二进制和前端资源，保留可空字段；旧代码不依赖 `structured_json`。不得因代码回滚删除报告数据。

只有确认新版本已停止写入、已识别字段内数据的保留要求，并单独获得生产数据库变更批准后，才可手工执行：

```sql
ALTER TABLE xy_compliance_monthly_report DROP COLUMN structured_json;
```

当前迁移框架不提供自动 down 脚本，因此不要将上述删列 SQL 加入常规回滚流程。

## PDF 中文字体与发布验证

月度报告导出按以下相对路径加载仓库内置字体：

```text
server/resource/captcha/fonts/SourceHanSansCN-Normal.ttf
```

构建和发布包必须保留 `resource/captcha/fonts/SourceHanSansCN-Normal.ttf`；服务工作目录应为 `server`，或保证运行时相对路径可见。当前月度报告不读取 `compliance.statementPdfFont`，该配置仅供原有对账单 PDF 使用。

```bash
cd server
test -s resource/captcha/fonts/SourceHanSansCN-Normal.ttf
file resource/captcha/fonts/SourceHanSansCN-Normal.ttf
```

发布前构建：

```bash
cd server
GOTOOLCHAIN=local go test ./internal/logic/compliance/report ./internal/library/reportpdf \
  ./internal/controller/member ./api/member ./internal/service -v
GOTOOLCHAIN=local go vet ./internal/logic/compliance/report ./internal/library/reportpdf
GOTOOLCHAIN=local go build ./...

cd ../web
export PATH="$HOME/.nvm/versions/node/v22.22.3/bin:$PATH"
corepack pnpm exec vue-tsc --noEmit
corepack pnpm exec eslint src/views/frontend/compliance/reports.vue \
  src/views/frontend/compliance/reports/*.vue src/api/frontend/compliance/assistant.ts
corepack pnpm exec stylelint src/views/frontend/compliance/reports.vue \
  src/views/frontend/compliance/reports/*.vue
corepack pnpm build
```

发布顺序：备份表结构和线上前端 `dist` 及二进制 → 执行 1.5.7 迁移 → 发布包含字体的新二进制和前端资源 → 重启服务 → 验证健康检查和以下冒烟路径。

```bash
# 使用测试会员的会话，不要在命令中硬编码凭据
curl -fsS -b /secure/path/member.cookie \
  'https://example.com/compliance/reports?periodKey=2026-07'

curl -fsS -D /tmp/monthly-report.headers -o /tmp/monthly-report.pdf \
  -b /secure/path/member.cookie \
  'https://example.com/compliance/reports/REPORT_ID/pdf'

grep -i '^Content-Type: application/pdf' /tmp/monthly-report.headers
grep -i '^Content-Disposition: attachment;' /tmp/monthly-report.headers
head -c 4 /tmp/monthly-report.pdf | grep -q '%PDF'
```

同时验证：创建草稿、查询列表、发布、草稿/已发布 PDF 导出；使用另一会员请求同一报告 ID 应返回“报告不存在或无权访问”。不要在生产企业上构造虚假流水或申报事实。

## 业务验收

- [ ] 开关关闭时，会员接口拒绝访问且旧模块正常。
- [ ] 保存企业画像后生成新版本；内容不变时不重复建版本。
- [ ] 相同规则、期间和触发键不会重复生成任务。
- [ ] 任务只能按定义状态机流转，跨会员访问被拒绝。
- [ ] 资料原附件保留；分类结果必须由用户确认后使用。
- [ ] 十条风险规则分别覆盖命中、边界和数据缺失情况。
- [ ] 同期间多类风险可共存；处置记录按会员隔离。
- [ ] 数据不足的报告明确显示“数据不足”。
- [ ] 空数据生成草稿时，总结显示数据不足，六个分类均可读且不输出确定性税务结论。
- [ ] 有异常报告按高/中/低风险展示，事实、依据、影响、建议、所需资料、截止日期、人工复核和规则版本完整可见。
- [ ] 无异常报告明确显示“未识别到异常”，不显示空白异常区。
- [ ] 无结构化快照的历史报告显示 `content` 摘要，PDF 导出也使用历史内容降级。
- [ ] 草稿和已发布版本均能导出 PDF，响应以 `%PDF` 开头且文件名包含期间和版本。
- [ ] 其他会员请求相同报告 ID 时导出被拒绝，返回“报告不存在或无权访问”。
- [ ] 375px 宽度下摘要、分类、异常长文和操作按钮不重叠、不溢出，全部文字可换行。
- [ ] 桌面宽度下摘要和六类检查布局对齐，不因长标题或空字段造成错位。
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
