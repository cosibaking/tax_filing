# Documents Upload Layout Fix Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 修复经营资料库上传区控件覆盖和错位，使桌面端三项同一行、小屏幕单列显示。

**Architecture:** 只修改经营资料库页面，通过页面级语义容器和 scoped CSS 约束三列网格。使用 Node 内置测试读取 Vue 单文件组件，先验证旧布局缺少边界约束，再验证修复后的结构和响应式规则。

**Tech Stack:** Vue 3、Element Plus、Scoped CSS、Node.js test runner、pnpm/Vite

---

### Task 1: 添加布局回归测试

**Files:**
- Create: `web/src/views/frontend/compliance/documents.layout.test.mjs`

- [ ] **Step 1: 写入失败测试**

```js
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import test from 'node:test'

const source = readFileSync(new URL('./documents.vue', import.meta.url), 'utf8')

test('documents upload row constrains all three desktop columns', () => {
  assert.match(source, /class="period-field"/)
  assert.match(source, /class="upload-field"/)
  assert.match(source, /grid-template-columns:\s*180px minmax\(0, 1fr\) auto/)
  assert.match(source, /\.upload-field\s*\{[^}]*min-width:\s*0/s)
})

test('documents upload row becomes full-width single column on mobile', () => {
  assert.match(source, /@media \(max-width: 760px\)/)
  assert.match(source, /\.period-field\s*[^}]*width:\s*100%/s)
})
```

- [ ] **Step 2: 验证测试因旧布局而失败**

运行：`node --test src/views/frontend/compliance/documents.layout.test.mjs`

预期：FAIL，提示缺少 `period-field` 或 `minmax(0, 1fr)`。

### Task 2: 实施页面局部布局修复

**Files:**
- Modify: `web/src/views/frontend/compliance/documents.vue`

- [ ] **Step 1: 增加三列语义容器**

```vue
<ElDatePicker class="period-field" ... />
<div class="upload-field">
  <MemberFileUpload ... />
</div>
<ElButton class="register-button" ...>登记资料</ElButton>
```

- [ ] **Step 2: 收紧网格和响应式规则**

```css
.upload-row {
  display: grid;
  grid-template-columns: 180px minmax(0, 1fr) auto;
  gap: 16px;
  align-items: start;
}
.period-field { width: 100%; }
.upload-field { min-width: 0; }
.register-button { align-self: start; }
@media (max-width: 760px) {
  .upload-row { grid-template-columns: minmax(0, 1fr); }
  .period-field, .upload-field, .register-button { width: 100%; }
}
```

- [ ] **Step 3: 验证回归测试转绿**

运行：`node --test src/views/frontend/compliance/documents.layout.test.mjs`

预期：2 tests passed。

- [ ] **Step 4: 执行前端质量检查**

运行：

```bash
pnpm exec prettier --check src/views/frontend/compliance/documents.vue
pnpm exec eslint src/views/frontend/compliance/documents.vue
pnpm build
```

预期：全部退出码为 0；Vite 可能保留既有的大分块警告。

- [ ] **Step 5: 提交修复**

```bash
git add web/src/views/frontend/compliance/documents.vue \
  web/src/views/frontend/compliance/documents.layout.test.mjs
git commit -m "fix: align documents upload controls"
git push
```

### Task 3: 发布并验证 82 前端

**Files:**
- Deploy artifact: `server/resource/public/dist`

- [ ] **Step 1: 打包并校验前端产物**

运行：`tar -C server/resource/public -czf /tmp/tax-filing-dist-layout-fix.tar.gz dist`

预期：归档包含 `dist/index.html`。

- [ ] **Step 2: 备份并替换服务器前端**

在 82 上将当前 `dist` 移入时间戳备份目录，再解压新产物；不修改后端配置和数据库。

- [ ] **Step 3: 验证线上状态**

检查：

```text
systemctl is-active tax-bridge.service = active
http://127.0.0.1:4096/ = 200
http://82.156.54.232/tax_filing/ = 200
```

并确认线上 `index.html` 引用本次构建的新资源哈希。
