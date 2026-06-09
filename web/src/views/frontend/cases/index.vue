<!-- 应用案例 -->
<template>
  <div class="cases-page">
    <header class="cases-page__header">
      <h1 class="cases-page__title">应用案例</h1>
      <p class="cases-page__subtitle">
        以下案例为金税管家服务场景示例，展示不同收入类型用户如何通过诊断、台账与申报协同实现合规管理。
      </p>
    </header>

    <div class="cases-page__list">
      <article v-for="item in appCases" :key="item.id" class="case-detail">
        <div class="case-detail__head">
          <span class="case-detail__tag">{{ item.tag }}</span>
          <span class="case-detail__highlight">{{ item.highlight }}</span>
        </div>
        <h2 class="case-detail__title">{{ item.title }}</h2>
        <p class="case-detail__profile">{{ item.profile }} · {{ item.income }}</p>

        <div class="case-detail__grid">
          <div class="case-detail__block">
            <h3>面临挑战</h3>
            <p>{{ item.pain }}</p>
          </div>
          <div class="case-detail__block">
            <h3>金税管家方案</h3>
            <p>{{ item.solution }}</p>
          </div>
          <div class="case-detail__block case-detail__block--outcome">
            <h3>服务成效</h3>
            <p>{{ item.outcome }}</p>
          </div>
        </div>
      </article>
    </div>

    <div class="cases-page__cta">
      <p>想了解您的税负结构与合规路径？</p>
      <RouterLink to="/diagnosis" class="opc-btn opc-btn--primary">免费合规诊断</RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { appCases } from '@/data/frontend/cases'

defineOptions({ name: 'FrontendCases' })
</script>

<style lang="scss" scoped>
.cases-page {
  max-width: 960px;
  margin: 0 auto;
  padding: 40px 24px 80px;
}

.cases-page__header {
  text-align: center;
  margin-bottom: 40px;
}

.cases-page__title {
  margin: 0 0 12px;
  font-size: 32px;
  font-weight: 800;
  color: #1a1f36;
}

.cases-page__subtitle {
  margin: 0 auto;
  max-width: 640px;
  font-size: 15px;
  line-height: 1.75;
  color: #6b7c93;
}

.cases-page__list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.case-detail {
  padding: 28px;
  background: #fff;
  border: 1px solid #e8edf3;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);
}

.case-detail__head {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}

.case-detail__tag {
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 700;
  color: #2563eb;
  background: #eff6ff;
}

.case-detail__highlight {
  font-size: 12px;
  color: #64748b;
}

.case-detail__title {
  margin: 0 0 8px;
  font-size: 20px;
  font-weight: 700;
  color: #1a1f36;
}

.case-detail__profile {
  margin: 0 0 20px;
  font-size: 14px;
  color: #6b7c93;
}

.case-detail__grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.case-detail__block {
  padding: 16px;
  border-radius: 10px;
  background: #f8fafc;

  h3 {
    margin: 0 0 8px;
    font-size: 13px;
    font-weight: 700;
    color: #334155;
  }

  p {
    margin: 0;
    font-size: 14px;
    line-height: 1.7;
    color: #64748b;
  }

  &--outcome {
    background: #f0fdf4;

    h3 {
      color: #166534;
    }

    p {
      color: #15803d;
    }
  }
}

.cases-page__cta {
  margin-top: 48px;
  text-align: center;

  p {
    margin: 0 0 16px;
    font-size: 15px;
    color: #6b7c93;
  }
}

.opc-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 160px;
  padding: 12px 28px;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.2s ease;

  &--primary {
    color: #fff;
    background: #2563eb;
    box-shadow: 0 4px 14px rgba(37, 99, 235, 0.35);

    &:hover {
      background: #1d4ed8;
      transform: translateY(-1px);
    }
  }
}

@media (max-width: 768px) {
  .cases-page {
    padding: 32px 16px 64px;
  }

  .case-detail__grid {
    grid-template-columns: 1fr;
  }
}
</style>
