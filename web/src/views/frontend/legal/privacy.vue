<!-- +----------------------------------------------------------------------
  | XYGo Admin — 隐私政策 P-14
  +---------------------------------------------------------------------- -->
<template>
  <main class="pt-8 pb-16 px-6">
    <div class="max-w-3xl mx-auto">
      <div class="bg-white/70 backdrop-blur-2xl rounded-[48px] shadow-clay-deep border border-[#d1d9e6]/40 p-8 md:p-12">
        <header class="mb-8 pb-6 border-b border-[#d1d9e6]/40">
          <h1 class="font-heading font-black text-3xl text-clay-foreground mb-3">隐私政策</h1>
          <div class="flex flex-wrap items-center gap-4 text-xs font-bold text-clay-muted">
            <span class="flex items-center gap-1">
              <ArtSvgIcon icon="ri:file-shield-2-line" class="text-sm" />
              版本 v1.0
            </span>
            <span class="flex items-center gap-1">
              <ArtSvgIcon icon="ri:time-line" class="text-sm" />
              更新日期：2026-06-08
            </span>
          </div>
        </header>

        <div class="legal-markdown">
          <MdPreview
            :model-value="privacyContent"
            :editor-id="previewId"
            preview-theme="github"
            :noMermaid="true"
            :noKatex="true"
          />
        </div>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

defineOptions({ name: 'LegalPrivacy' })

const previewId = 'legal-privacy-preview'

const privacyContent = `## 一、引言

欢迎使用 OPC 合规服务平台（以下简称「本平台」）。我们深知个人信息对您的重要性，将按照法律法规要求，采取相应安全保护措施，尽力保护您的个人信息安全可控。

## 二、我们收集的信息

在您使用本平台服务时，我们可能收集以下信息：

1. **账号信息**：用户名、手机号、邮箱、密码（加密存储）
2. **诊断信息**：平台类型、收入区间、经营主体、报税情况等合规诊断问卷数据
3. **经营主体信息**：公司名称、统一社会信用代码、经营范围等 OPC 注册所需资料
4. **财务信息**：收入台账、费用台账、银行流水等记账申报相关数据
5. **设备与日志信息**：IP 地址、浏览器类型、操作日志等

## 三、信息使用目的

我们收集的信息将用于：

- 提供合规诊断、税负计算、方案推荐服务
- 办理 OPC 设立、记账、申报等合规服务
- 履行法律法规规定的义务
- 改善产品体验与服务质量
- 发送服务通知与重要提醒

## 四、信息共享与披露

我们不会向第三方出售您的个人信息。仅在以下情形共享：

- 经您明确同意
- 为完成 OPC 注册、银行开户等必要服务，向合作机构提供（签署保密协议）
- 法律法规要求或政府机关依法要求

## 五、信息存储与保护

- 数据存储于中华人民共和国境内
- 采用加密传输与存储、访问控制等安全措施
- 敏感信息（如身份证号、银行账号）进行脱敏或加密处理
- 留存期限：服务期间及法律法规要求的最短保存期限

## 六、您的权利

您有权：

- 查询、更正您的个人信息
- 删除诊断记录（法律法规另有规定除外）
- 撤回授权同意
- 注销账号

如需行使上述权利，请通过平台客服联系我们。

## 七、Cookie 与同类技术

我们使用 Cookie 维持登录状态、记住偏好设置。您可通过浏览器设置管理 Cookie，但可能影响部分功能使用。

## 八、未成年人保护

本平台服务面向年满 18 周岁的成年人。如您为未成年人，请在监护人陪同下阅读本政策并使用服务。

## 九、政策更新

我们可能适时修订本政策。重大变更将通过平台公告或站内通知告知您。继续使用服务即视为接受更新后的政策。

## 十、联系我们

如对本隐私政策有任何疑问，请通过平台客服或注册邮箱联系我们。`
</script>

<style lang="scss" scoped>
.text-clay-foreground { color: #32325d; }
.text-clay-muted { color: #8898aa; }
.font-heading { font-family: 'Nunito', 'PingFang SC', sans-serif; }

.shadow-clay-deep {
  box-shadow: 30px 30px 60px #d1d9e6, -30px -30px 60px #ffffff,
    inset 10px 10px 20px rgba(90, 141, 238, 0.05), inset -10px -10px 20px rgba(255, 255, 255, 0.8);
}

:deep(.legal-markdown) {
  .md-editor-preview-wrapper { padding: 0; }
  .md-editor-preview {
    color: #525f7f;
    font-size: 15px;
    line-height: 1.8;
    h2 {
      font-family: 'Nunito', 'PingFang SC', sans-serif;
      font-weight: 800;
      color: #32325d;
      font-size: 1.25rem;
      margin-top: 2rem;
      margin-bottom: 0.75rem;
    }
    p, li { color: #525f7f; }
    strong { color: #32325d; }
  }
}
</style>
