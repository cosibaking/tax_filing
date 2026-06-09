<!-- 法律文档通用页面（隐私政策 / 用户协议） -->
<template>
  <main class="legal-page">
    <div class="legal-page__inner">
      <section class="legal-panel">
        <header class="legal-panel__head">
          <div>
            <h1 class="legal-panel__title">{{ title }}</h1>
            <p class="legal-panel__subtitle">{{ subtitle }}</p>
          </div>
          <div class="legal-panel__meta">
            <span><ArtSvgIcon :icon="icon" /> {{ version }}</span>
            <span><ArtSvgIcon icon="ri:time-line" /> 更新：{{ updatedAt }}</span>
          </div>
        </header>

        <nav class="legal-panel__nav">
          <RouterLink
            to="/legal/privacy"
            class="legal-nav-link"
            :class="{ 'is-active': docType === 'privacy' }"
          >
            隐私政策
          </RouterLink>
          <RouterLink
            to="/legal/terms"
            class="legal-nav-link"
            :class="{ 'is-active': docType === 'terms' }"
          >
            用户协议
          </RouterLink>
        </nav>

        <div class="legal-markdown">
          <MdPreview
            :model-value="content"
            :editor-id="previewId"
            preview-theme="github"
            :no-mermaid="true"
            :no-katex="true"
          />
        </div>

        <footer class="legal-panel__footer">
          <div class="legal-panel__contact">
            <p>运营主体：{{ companyName }}</p>
            <p>平台名称：{{ platformName }}</p>
            <p>联系邮箱：<a :href="`mailto:${contactEmail}`" class="legal-footer-link">{{ contactEmail }}</a></p>
          </div>
          <div class="legal-panel__footer-links">
            <RouterLink to="/" class="legal-footer-link">返回首页</RouterLink>
            <RouterLink
              :to="docType === 'privacy' ? '/legal/terms' : '/legal/privacy'"
              class="legal-footer-link"
            >
              {{ docType === 'privacy' ? '查看用户协议' : '查看隐私政策' }}
            </RouterLink>
          </div>
        </footer>
      </section>
    </div>
  </main>
</template>

<script setup lang="ts">
import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

defineOptions({ name: 'LegalDocumentPage' })

withDefaults(defineProps<{
  title: string
  subtitle?: string
  content: string
  version: string
  updatedAt: string
  companyName: string
  platformName: string
  contactEmail: string
  previewId: string
  docType: 'privacy' | 'terms'
  icon?: string
}>(), {
  subtitle: '请仔细阅读以下条款，了解您的权利与义务',
  icon: 'ri:file-shield-2-line',
})
</script>

<style lang="scss" scoped>
.legal-page {
  padding: 32px 24px 64px;
}

.legal-page__inner {
  max-width: 800px;
  margin: 0 auto;
}

.legal-panel {
  padding: 28px 28px 24px;
  background: #fff;
  border: 1px solid #e8edf3;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);
}

.legal-panel__head {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  padding-bottom: 20px;
  margin-bottom: 20px;
  border-bottom: 1px solid #e8edf3;
}

.legal-panel__title {
  margin: 0 0 6px;
  font-size: 24px;
  font-weight: 800;
  color: #1a1f36;
}

.legal-panel__subtitle {
  margin: 0;
  font-size: 14px;
  color: #6b7c93;
  line-height: 1.5;
}

.legal-panel__meta {
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: 12px;
  font-weight: 600;
  color: #94a3b8;

  span {
    display: inline-flex;
    align-items: center;
    gap: 6px;
  }
}

.legal-panel__nav {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
}

.legal-nav-link {
  padding: 8px 16px;
  border-radius: 8px;
  border: 1px solid #e8edf3;
  background: #f8fafc;
  font-size: 13px;
  font-weight: 600;
  color: #475569;
  text-decoration: none;
  transition: all 0.15s ease;

  &:hover {
    border-color: #bfdbfe;
    color: #2563eb;
  }

  &.is-active {
    background: #2563eb;
    border-color: #2563eb;
    color: #fff;
  }
}

.legal-panel__footer {
  margin-top: 32px;
  padding-top: 20px;
  border-top: 1px solid #e8edf3;
}

.legal-panel__contact {
  margin-bottom: 16px;

  p {
    margin: 0 0 6px;
    font-size: 13px;
    color: #94a3b8;
  }
}

.legal-panel__footer-links {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.legal-footer-link {
  font-size: 13px;
  font-weight: 600;
  color: #2563eb;
  text-decoration: none;

  &:hover {
    text-decoration: underline;
  }
}

:deep(.legal-markdown) {
  .md-editor-preview-wrapper {
    padding: 0;
  }

  .md-editor-preview {
    color: #475569;
    font-size: 14px;
    line-height: 1.8;
    background: transparent;

    h2 {
      margin: 28px 0 12px;
      padding-bottom: 8px;
      border-bottom: 1px solid #f1f5f9;
      font-size: 17px;
      font-weight: 700;
      color: #1a1f36;

      &:first-child {
        margin-top: 0;
      }
    }

    h3 {
      margin: 20px 0 8px;
      font-size: 15px;
      font-weight: 700;
      color: #334155;
    }

    p {
      margin: 0 0 12px;
      color: #475569;
    }

    strong {
      color: #1a1f36;
      font-weight: 700;
    }

    ul, ol {
      margin: 0 0 12px;
      padding-left: 1.4em;
    }

    li {
      margin-bottom: 6px;
      color: #475569;
    }

    table {
      width: 100%;
      margin: 12px 0 16px;
      border-collapse: collapse;
      font-size: 13px;
    }

    th, td {
      padding: 10px 12px;
      border: 1px solid #e8edf3;
      text-align: left;
    }

    th {
      background: #f8fafc;
      font-weight: 700;
      color: #334155;
    }

    hr {
      margin: 24px 0;
      border: none;
      border-top: 1px solid #e8edf3;
    }
  }
}

@media (max-width: 640px) {
  .legal-panel {
    padding: 20px 16px;
  }

  .legal-panel__title {
    font-size: 20px;
  }
}
</style>
