<!-- 隐私政策 -->
<template>
  <LegalDocumentPage
    title="隐私政策"
    subtitle="了解金税管家如何收集、使用与保护您的个人信息"
    :content="PRIVACY_CONTENT"
    :version="LEGAL_META.privacyVersion"
    :updated-at="LEGAL_META.updatedAt"
    :company-name="LEGAL_META.companyName"
    :platform-name="LEGAL_META.platformName"
    :contact-email="LEGAL_META.contactEmail"
    preview-id="legal-privacy-preview"
    doc-type="privacy"
    icon="ri:file-shield-2-line"
  />
</template>

<script setup lang="ts">
import LegalDocumentPage from '@/components/frontend/LegalDocumentPage.vue'
import { LEGAL_META, PRIVACY_CONTENT } from '@/data/frontend/legal'

defineOptions({ name: 'LegalPrivacy' })
</script>
