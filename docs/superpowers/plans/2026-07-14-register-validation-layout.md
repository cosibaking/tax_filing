# 创建账号页校验提示布局修复 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 修复创建账号页全部表单校验提示与输入框、协议文字及注册按钮重叠的问题。

**Architecture:** 保持现有 Element Plus 表单和校验规则不变，仅在注册页使用专用布局类管理表单项间距，并把错误提示从绝对定位调整为正常文档流。该修复通过 scoped 样式隔离，不影响登录页和其他表单。

**Tech Stack:** Vue 3、Element Plus、SCSS、Tailwind CSS、pnpm、Vue TypeScript

---

## 文件结构

- Modify: `web/src/views/frontend/member/register.vue`：调整注册表单的结构类与局部校验提示样式。
- Modify: `docs/superpowers/plans/2026-07-14-register-validation-layout.md`：执行过程中勾选任务状态。

### Task 1: 建立失败复现断言

**Files:**
- Inspect: `web/src/views/frontend/member/register.vue`

- [ ] **Step 1: 运行当前结构断言，证明页面仍使用会导致重叠的零底边距**

Run:

```bash
cd web
node -e "const s=require('fs').readFileSync('src/views/frontend/member/register.vue','utf8'); if(!s.includes('class=\"!mb-0\"')) process.exit(1); console.log('FAIL reproduction: zero-margin form items found')"
```

Expected: 输出 `FAIL reproduction: zero-margin form items found`，证明当前页面仍存在根因结构。

- [ ] **Step 2: 运行修复目标断言并确认当前失败**

Run:

```bash
cd web
node -e "const s=require('fs').readFileSync('src/views/frontend/member/register.vue','utf8'); const ok=s.includes('register-form')&&s.includes('register-form-item')&&s.includes('position: static'); if(!ok){console.error('EXPECTED FAIL: flow-based validation layout missing');process.exit(1)}"
```

Expected: 退出码为 1，并输出 `EXPECTED FAIL: flow-based validation layout missing`。

### Task 2: 实现正常文档流的校验提示布局

**Files:**
- Modify: `web/src/views/frontend/member/register.vue:38-107`
- Modify: `web/src/views/frontend/member/register.vue:202-225`

- [ ] **Step 1: 替换表单和表单项布局类**

将表单根节点改为：

```vue
<ElForm
  ref="formRef"
  :model="formData"
  :rules="rules"
  class="register-form"
  @keyup.enter="handleSubmit"
>
```

四个输入项统一使用：

```vue
<ElFormItem prop="字段名" class="register-form-item">
```

协议容器和协议项使用：

```vue
<div class="register-agreements">
  <ElFormItem prop="agreeTerms" class="register-agreement-item">
  <ElFormItem prop="agreePrivacy" class="register-agreement-item">
</div>
```

提交按钮移除 `mt-4`，添加 `register-submit`，由表单布局统一控制垂直间距。

- [ ] **Step 2: 增加页面隔离的流式校验样式**

在 scoped SCSS 中加入：

```scss
.register-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.register-form-item,
.register-agreement-item {
  margin-bottom: 0;
}

.register-agreements {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 2px 4px 0;
}

.register-submit {
  margin-top: 2px;
}

:deep(.register-form-item .el-form-item__content),
:deep(.register-agreement-item .el-form-item__content) {
  display: flex;
  flex-direction: column;
  align-items: stretch;
}

:deep(.register-form .el-form-item__error) {
  position: static;
  width: 100%;
  padding-top: 6px;
  line-height: 1.45;
}
```

- [ ] **Step 3: 运行修复目标断言并确认通过**

Run:

```bash
cd web
node -e "const s=require('fs').readFileSync('src/views/frontend/member/register.vue','utf8'); const ok=s.includes('register-form')&&s.includes('register-form-item')&&s.includes('position: static')&&!s.includes('class=\"!mb-0\"'); if(!ok)process.exit(1); console.log('PASS: validation messages participate in layout')"
```

Expected: 输出 `PASS: validation messages participate in layout`。

- [ ] **Step 4: 格式化并检查单文件 lint**

Run:

```bash
cd web
pnpm exec prettier --write src/views/frontend/member/register.vue
pnpm exec eslint src/views/frontend/member/register.vue
pnpm exec stylelint src/views/frontend/member/register.vue
```

Expected: 三条命令退出码均为 0。

- [ ] **Step 5: 提交最小修复**

```bash
git add web/src/views/frontend/member/register.vue docs/superpowers/plans/2026-07-14-register-validation-layout.md
git commit -m "fix: prevent register validation messages overlapping"
```

### Task 3: 完整构建与视觉验收

**Files:**
- Verify: `web/src/views/frontend/member/register.vue`

- [ ] **Step 1: 执行前端类型检查与生产构建**

Run:

```bash
cd web
pnpm build
```

Expected: `vue-tsc --noEmit` 和 `vite build` 成功，退出码为 0。

- [ ] **Step 2: 执行项目 lint**

Run:

```bash
cd web
pnpm lint
```

Expected: 退出码为 0；若存在与本次文件无关的历史问题，记录具体文件，并确保注册页单文件 lint 已通过。

- [ ] **Step 3: 在 375px 宽度验证全部错误同时出现**

启动前端后打开 `/user/register`，保持所有字段为空并点击「立即注册」。确认：

```text
用户名提示位于用户名输入框下方
密码提示位于密码输入框下方
确认密码提示位于确认密码输入框下方
手机号提示位于手机号输入框下方
两条协议提示分别位于对应复选框下方
注册按钮位于全部协议提示下方
页面可自然滚动，无横向滚动条
```

- [ ] **Step 4: 在桌面宽度验证错误清除后的布局**

依次填写合法用户名、密码、确认密码和手机号，并勾选协议。确认所有错误消失，字段间距一致，无异常大块空白，提交按钮和返回登录区域位置正常。

- [ ] **Step 5: 检查最终差异**

Run:

```bash
git diff --check HEAD~1..HEAD
git status --short --branch
```

Expected: `git diff --check` 无输出；工作区没有遗漏的非计划改动。

## 回滚

若上线后出现布局回归，回滚实现提交即可：

```bash
git revert <implementation-commit>
```

该回滚不涉及数据库、接口或生产配置。
