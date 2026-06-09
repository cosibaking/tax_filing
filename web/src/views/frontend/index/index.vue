<!-- +----------------------------------------------------------------------
  | P-01 首页 — 主播 OPC 合规营销落地页
  +---------------------------------------------------------------------- -->
<template>
  <div class="opc-home">
    <!-- Hero -->
    <section class="opc-hero">
      <h1 class="opc-hero__title">个人税务合规，从诊断到申报一站搞定</h1>
      <p class="opc-hero__subtitle">
        税务诊断、收入费用台账、申报提醒与月度对账单，让合规可执行、可追踪。
      </p>
      <div class="opc-hero__actions">
        <RouterLink to="/diagnosis" class="opc-btn opc-btn--primary">免费合规诊断</RouterLink>
        <RouterLink to="/pricing" class="opc-btn opc-btn--secondary">查看服务套餐</RouterLink>
      </div>
    </section>

    <!-- 三步流程 -->
    <section class="opc-steps">
      <h2 class="opc-steps__title">三步开启合规服务</h2>
      <div class="opc-steps__grid">
        <div v-for="(step, index) in steps" :key="step.title" class="opc-step-card">
          <div class="opc-step-card__badge">步骤 {{ index + 1 }}</div>
          <h3 class="opc-step-card__title">{{ step.title }}</h3>
          <p class="opc-step-card__desc">{{ step.desc }}</p>
        </div>
      </div>
    </section>

    <!-- 应用案例 -->
    <section id="cases" class="opc-cases">
      <div class="opc-cases__header">
        <h2 class="opc-cases__title">应用案例</h2>
        <p class="opc-cases__subtitle">
          不同收入类型用户，如何通过金税管家实现合规诊断、台账管理与按时申报。
        </p>
      </div>
      <div class="opc-cases__grid">
        <article v-for="item in previewCases" :key="item.id" class="opc-case-card">
          <span class="opc-case-card__tag">{{ item.tag }}</span>
          <h3 class="opc-case-card__title">{{ item.title }}</h3>
          <p class="opc-case-card__profile">{{ item.profile }}</p>
          <p class="opc-case-card__outcome">{{ item.outcome }}</p>
        </article>
      </div>
      <div class="opc-cases__more">
        <RouterLink to="/cases" class="opc-btn opc-btn--secondary">查看全部案例</RouterLink>
      </div>
    </section>

    <!-- 关于我们 -->
    <section id="about" class="opc-about">
      <div class="opc-about__header">
        <h2 class="opc-about__title">关于金税管家</h2>
        <p class="opc-about__subtitle">{{ aboutIntro.subtitle }}</p>
      </div>
      <p class="opc-about__desc">{{ aboutIntro.description }}</p>
      <div class="opc-about__values">
        <article v-for="item in aboutValues" :key="item.title" class="opc-about-card">
          <ArtSvgIcon :icon="item.icon" class="opc-about-card__icon" />
          <h3>{{ item.title }}</h3>
          <p>{{ item.desc }}</p>
        </article>
      </div>
      <div class="opc-about__more">
        <RouterLink to="/about" class="opc-btn opc-btn--secondary">了解更多</RouterLink>
      </div>
    </section>

    <!-- 重要告知 -->
    <section class="opc-notice">
      <div class="opc-notice__box">
        <h3 class="opc-notice__title">重要告知</h3>
        <p class="opc-notice__text">
          本服务为合法合规税务优化方案，非逃税避税方案；不承诺「零风险」「包不被查」。税务筹划基于现行法律法规及优惠政策实施，最终执行结果以税务机关认定为准。请您确保提供资料真实、完整，并在签约前仔细阅读《服务协议》及《风险告知书》。
        </p>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { appCases } from '@/data/frontend/cases'
import { aboutIntro, aboutValues } from '@/data/frontend/about'

defineOptions({ name: 'FrontendHome' })

const previewCases = appCases.slice(0, 3)

const steps = [
  {
    title: '免费诊断',
    desc: '填写问卷，了解税负对比与合规风险'
  },
  {
    title: '签约服务',
    desc: '选择套餐，完成风险告知与电子签约'
  },
  {
    title: '记账申报',
    desc: '录入流水，顾问协助申报与对账'
  }
]
</script>

<style lang="scss" scoped>
.opc-home {
  max-width: 960px;
  margin: 0 auto;
  padding: 48px 24px 80px;
}

.opc-hero {
  text-align: center;
  padding: 24px 0 56px;
}

.opc-hero__title {
  margin: 0 0 20px;
  font-size: clamp(28px, 5vw, 40px);
  font-weight: 800;
  line-height: 1.35;
  color: #1a1f36;
  letter-spacing: -0.02em;
}

.opc-hero__subtitle {
  margin: 0 auto 36px;
  max-width: 640px;
  font-size: 16px;
  line-height: 1.75;
  color: #6b7c93;
}

.opc-hero__actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 16px;
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

  &--secondary {
    color: #334155;
    background: #fff;
    border: 1px solid #d8dee9;

    &:hover {
      border-color: #94a3b8;
      background: #f8fafc;
    }
  }
}

.opc-steps {
  padding: 16px 0 48px;
}

.opc-steps__title {
  margin: 0 0 28px;
  text-align: center;
  font-size: 22px;
  font-weight: 700;
  color: #1a1f36;
}

.opc-steps__grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

.opc-step-card {
  padding: 28px 24px;
  background: #fff;
  border: 1px solid #e8edf3;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);
  transition: box-shadow 0.2s ease, transform 0.2s ease;

  &:hover {
    box-shadow: 0 8px 24px rgba(15, 23, 42, 0.08);
    transform: translateY(-2px);
  }
}

.opc-step-card__badge {
  display: inline-block;
  margin-bottom: 14px;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 700;
  color: #2563eb;
  background: #eff6ff;
}

.opc-step-card__title {
  margin: 0 0 10px;
  font-size: 18px;
  font-weight: 700;
  color: #1a1f36;
}

.opc-step-card__desc {
  margin: 0;
  font-size: 14px;
  line-height: 1.65;
  color: #6b7c93;
}

.opc-cases {
  padding: 16px 0 48px;
}

.opc-cases__header {
  text-align: center;
  margin-bottom: 28px;
}

.opc-cases__title {
  margin: 0 0 10px;
  font-size: 22px;
  font-weight: 700;
  color: #1a1f36;
}

.opc-cases__subtitle {
  margin: 0 auto;
  max-width: 560px;
  font-size: 14px;
  line-height: 1.7;
  color: #6b7c93;
}

.opc-cases__grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

.opc-case-card {
  padding: 24px;
  background: #fff;
  border: 1px solid #e8edf3;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);
  transition: box-shadow 0.2s ease, transform 0.2s ease;

  &:hover {
    box-shadow: 0 8px 24px rgba(15, 23, 42, 0.08);
    transform: translateY(-2px);
  }
}

.opc-case-card__tag {
  display: inline-block;
  margin-bottom: 12px;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 700;
  color: #2563eb;
  background: #eff6ff;
}

.opc-case-card__title {
  margin: 0 0 8px;
  font-size: 16px;
  font-weight: 700;
  line-height: 1.45;
  color: #1a1f36;
}

.opc-case-card__profile {
  margin: 0 0 12px;
  font-size: 13px;
  color: #94a3b8;
}

.opc-case-card__outcome {
  margin: 0;
  font-size: 14px;
  line-height: 1.65;
  color: #64748b;
}

.opc-cases__more {
  margin-top: 28px;
  text-align: center;
}

.opc-about {
  padding: 16px 0 48px;
}

.opc-about__header {
  text-align: center;
  margin-bottom: 20px;
}

.opc-about__title {
  margin: 0 0 10px;
  font-size: 22px;
  font-weight: 700;
  color: #1a1f36;
}

.opc-about__subtitle {
  margin: 0 auto;
  max-width: 560px;
  font-size: 14px;
  line-height: 1.7;
  color: #6b7c93;
}

.opc-about__desc {
  margin: 0 auto 28px;
  max-width: 720px;
  font-size: 15px;
  line-height: 1.85;
  text-align: center;
  color: #475569;
}

.opc-about__values {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.opc-about-card {
  padding: 22px 18px;
  background: #fff;
  border: 1px solid #e8edf3;
  border-radius: 12px;
  text-align: center;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);

  h3 {
    margin: 12px 0 8px;
    font-size: 15px;
    font-weight: 700;
    color: #1a1f36;
  }

  p {
    margin: 0;
    font-size: 13px;
    line-height: 1.6;
    color: #64748b;
  }
}

.opc-about-card__icon {
  font-size: 26px;
  color: #2563eb;
}

.opc-about__more {
  margin-top: 28px;
  text-align: center;
}

.opc-notice {
  padding-top: 8px;
}

.opc-notice__box {
  padding: 24px 28px;
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: 12px;
}

.opc-notice__title {
  margin: 0 0 14px;
  font-size: 16px;
  font-weight: 700;
  color: #92400e;
}

.opc-notice__text {
  margin: 0;
  color: #78350f;
  font-size: 14px;
  line-height: 1.85;
}

@media (max-width: 768px) {
  .opc-home {
    padding: 32px 16px 64px;
  }

  .opc-steps__grid,
  .opc-cases__grid,
  .opc-about__values {
    grid-template-columns: 1fr;
  }

  .opc-hero__actions {
    flex-direction: column;
    align-items: stretch;
  }

  .opc-btn {
    width: 100%;
  }
}
</style>
