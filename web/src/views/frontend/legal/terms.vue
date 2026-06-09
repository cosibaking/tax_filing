<!-- +----------------------------------------------------------------------
  | XYGo Admin — 用户协议 P-14
  +---------------------------------------------------------------------- -->
<template>
  <main class="pt-8 pb-16 px-6">
    <div class="max-w-3xl mx-auto">
      <div class="bg-white/70 backdrop-blur-2xl rounded-[48px] shadow-clay-deep border border-[#d1d9e6]/40 p-8 md:p-12">
        <header class="mb-8 pb-6 border-b border-[#d1d9e6]/40">
          <h1 class="font-heading font-black text-3xl text-clay-foreground mb-3">用户服务协议</h1>
          <div class="flex flex-wrap items-center gap-4 text-xs font-bold text-clay-muted">
            <span class="flex items-center gap-1">
              <ArtSvgIcon icon="ri:file-text-line" class="text-sm" />
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
            :model-value="termsContent"
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

defineOptions({ name: 'LegalTerms' })

const previewId = 'legal-terms-preview'

const termsContent = `## 一、协议范围

本协议是您与金税管家平台（以下简称「本平台」）之间关于使用本平台服务所订立的协议。注册或使用本平台服务，即表示您已阅读并同意本协议。

## 二、服务内容

本平台为个人及小微企业提供以下服务：

1. 合规诊断与税负对比分析
2. 经营主体设立代办
3. 日常记账与财务台账管理
4. 税务申报代办与提醒
5. 合规咨询与月度对账单

具体服务内容以您选择的服务套餐为准。

## 三、用户注册与账号

1. 您应提供真实、准确、完整的注册信息
2. 您应妥善保管账号密码，对账号下的行为负责
3. 禁止转让、出借账号给他人使用
4. 发现账号异常请立即通知本平台

## 四、服务费用

1. 服务费用以套餐公示价格为准
2. 费用支付方式及周期以签约时约定为准
3. 本平台有权调整价格，调整前已签约订单不受影响

## 五、用户义务

您承诺：

1. 提供真实、完整的收入与成本信息
2. 不利用本平台从事逃税、虚开发票等违法行为
3. 配合提供主体设立、申报所需的合法资料
4. 及时响应平台的资料补充与确认请求
5. 遵守国家税收法律法规

## 六、平台义务与免责

1. 本平台按约定提供合规服务，尽力确保服务质量
2. 本平台不承诺「包不被税务稽查」
3. 税负估算仅供参考，实际税负以税务机关核定为准
4. 因您提供虚假信息导致的后果由您自行承担
5. 因不可抗力、政策变化导致的服务中断，本平台不承担责任

## 七、知识产权

本平台上的文字、图片、软件、界面设计等知识产权归本平台或权利人所有。未经授权，不得复制、传播或用于商业用途。

## 八、协议变更与终止

1. 本平台有权修订本协议，修订后继续使用视为接受
2. 您可随时停止使用服务并申请注销账号
3. 如您严重违反本协议，本平台有权暂停或终止服务

## 九、争议解决

本协议适用中华人民共和国法律。因本协议引起的争议，双方应友好协商；协商不成的，提交本平台所在地有管辖权的人民法院诉讼解决。

## 十、联系方式

如对本协议有任何疑问，请通过平台客服联系我们。`
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
