<!-- 用户服务协议 -->
<template>
  <LegalDocumentPage
    title="用户服务协议"
    subtitle="使用金税管家服务前，请阅读并理解以下条款"
    :content="TERMS_CONTENT"
    :version="LEGAL_META.termsVersion"
    :updated-at="LEGAL_META.updatedAt"
    :company-name="LEGAL_META.companyName"
    :platform-name="LEGAL_META.platformName"
    :contact-email="LEGAL_META.contactEmail"
    preview-id="legal-terms-preview"
    doc-type="terms"
    icon="ri:file-text-line"
  />
</template>

<script setup lang="ts">
import LegalDocumentPage from '@/components/frontend/LegalDocumentPage.vue'
import { LEGAL_META, TERMS_CONTENT } from '@/data/frontend/legal'

defineOptions({ name: 'LegalTerms' })
</script>
