<!-- 关于我们 -->
<template>
  <div class="about-page">
    <header class="about-page__header">
      <h1 class="about-page__title">{{ aboutIntro.title }}</h1>
      <p class="about-page__subtitle">{{ aboutIntro.subtitle }}</p>
    </header>

    <section class="about-card about-card--intro">
      <p class="about-card__text">{{ aboutIntro.description }}</p>
      <div class="about-mission">
        <div class="about-mission__item">
          <span class="about-mission__label">使命</span>
          <p>{{ aboutMission.mission }}</p>
        </div>
        <div class="about-mission__item">
          <span class="about-mission__label">愿景</span>
          <p>{{ aboutMission.vision }}</p>
        </div>
      </div>
    </section>

    <section class="about-section">
      <h2 class="about-section__title">我们的优势</h2>
      <div class="about-values">
        <article v-for="item in aboutValues" :key="item.title" class="about-value">
          <ArtSvgIcon :icon="item.icon" class="about-value__icon" />
          <h3>{{ item.title }}</h3>
          <p>{{ item.desc }}</p>
        </article>
      </div>
    </section>

    <section class="about-section">
      <h2 class="about-section__title">服务内容</h2>
      <ul class="about-services">
        <li v-for="service in aboutServices" :key="service">
          <ArtSvgIcon icon="ri:check-line" class="about-services__icon" />
          <span>{{ service }}</span>
        </li>
      </ul>
    </section>

    <section class="about-section">
      <h2 class="about-section__title">发展历程</h2>
      <div class="about-timeline">
        <article v-for="item in aboutMilestones" :key="item.title" class="about-timeline__item">
          <span class="about-timeline__year">{{ item.year }}</span>
          <h3>{{ item.title }}</h3>
          <p>{{ item.desc }}</p>
        </article>
      </div>
    </section>

    <section class="about-card about-card--contact">
      <h2 class="about-section__title">联系我们</h2>
      <div class="about-contact">
        <p><strong>服务时间：</strong>{{ aboutContact.serviceTime }}</p>
        <p><strong>联系邮箱：</strong>{{ aboutContact.email }}</p>
        <p class="about-contact__note">{{ aboutContact.note }}</p>
      </div>
    </section>

    <div class="about-page__cta">
      <RouterLink to="/diagnosis" class="opc-btn opc-btn--primary">免费合规诊断</RouterLink>
      <RouterLink to="/pricing" class="opc-btn opc-btn--secondary">查看服务套餐</RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  aboutIntro,
  aboutMission,
  aboutValues,
  aboutServices,
  aboutMilestones,
  aboutContact
} from '@/data/frontend/about'

defineOptions({ name: 'FrontendAbout' })
</script>

<style lang="scss" scoped>
.about-page {
  max-width: 960px;
  margin: 0 auto;
  padding: 40px 24px 80px;
}

.about-page__header {
  text-align: center;
  margin-bottom: 36px;
}

.about-page__title {
  margin: 0 0 12px;
  font-size: 32px;
  font-weight: 800;
  color: #1a1f36;
}

.about-page__subtitle {
  margin: 0 auto;
  max-width: 640px;
  font-size: 15px;
  line-height: 1.75;
  color: #6b7c93;
}

.about-card {
  padding: 28px;
  margin-bottom: 32px;
  background: #fff;
  border: 1px solid #e8edf3;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);
}

.about-card__text {
  margin: 0 0 24px;
  font-size: 15px;
  line-height: 1.85;
  color: #475569;
}

.about-mission {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.about-mission__item {
  padding: 18px 20px;
  border-radius: 10px;
  background: #f8fafc;
  border: 1px solid #e8edf3;

  p {
    margin: 8px 0 0;
    font-size: 14px;
    line-height: 1.65;
    color: #334155;
    font-weight: 600;
  }
}

.about-mission__label {
  font-size: 12px;
  font-weight: 800;
  color: #2563eb;
  letter-spacing: 0.04em;
}

.about-section {
  margin-bottom: 32px;
}

.about-section__title {
  margin: 0 0 20px;
  font-size: 20px;
  font-weight: 700;
  color: #1a1f36;
}

.about-values {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.about-value {
  padding: 22px;
  background: #fff;
  border: 1px solid #e8edf3;
  border-radius: 12px;

  h3 {
    margin: 12px 0 8px;
    font-size: 16px;
    font-weight: 700;
    color: #1a1f36;
  }

  p {
    margin: 0;
    font-size: 14px;
    line-height: 1.65;
    color: #64748b;
  }
}

.about-value__icon {
  font-size: 24px;
  color: #2563eb;
}

.about-services {
  margin: 0;
  padding: 0;
  list-style: none;
  display: grid;
  gap: 12px;
}

.about-services li {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 14px 18px;
  background: #fff;
  border: 1px solid #e8edf3;
  border-radius: 10px;
  font-size: 14px;
  color: #334155;
  font-weight: 600;
}

.about-services__icon {
  margin-top: 2px;
  font-size: 16px;
  color: #16a34a;
  flex-shrink: 0;
}

.about-timeline {
  display: grid;
  gap: 14px;
}

.about-timeline__item {
  padding: 20px 22px;
  background: #fff;
  border-left: 3px solid #2563eb;
  border-radius: 0 10px 10px 0;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);

  h3 {
    margin: 6px 0 8px;
    font-size: 16px;
    font-weight: 700;
    color: #1a1f36;
  }

  p {
    margin: 0;
    font-size: 14px;
    line-height: 1.65;
    color: #64748b;
  }
}

.about-timeline__year {
  font-size: 12px;
  font-weight: 800;
  color: #2563eb;
}

.about-contact {
  p {
    margin: 0 0 10px;
    font-size: 14px;
    line-height: 1.7;
    color: #475569;
  }
}

.about-contact__note {
  margin-top: 14px !important;
  padding-top: 14px;
  border-top: 1px solid #e8edf3;
  color: #94a3b8 !important;
  font-size: 13px !important;
}

.about-page__cta {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 16px;
  margin-top: 12px;
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
    }
  }

  &--secondary {
    color: #334155;
    background: #fff;
    border: 1px solid #d8dee9;

    &:hover {
      background: #f8fafc;
    }
  }
}

@media (max-width: 768px) {
  .about-mission,
  .about-values {
    grid-template-columns: 1fr;
  }

  .about-page__cta {
    flex-direction: column;
    align-items: stretch;
  }

  .opc-btn {
    width: 100%;
  }
}
</style>
