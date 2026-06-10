<!-- 服务套餐 -->
<template>
  <div class="pricing-page">
    <header class="pricing-page__header">
      <span class="pricing-page__badge">税务合规服务</span>
      <h1 class="pricing-page__title">选择适合您的 <span class="is-accent">服务套餐</span></h1>
      <p class="pricing-page__subtitle">
        从主体设立到日常记账申报，全程代办，让您专注经营本身
      </p>
    </header>

    <div v-if="loading" class="pricing-page__loading">加载套餐信息...</div>

    <div v-else class="pricing-page__grid">
      <article
        v-for="plan in displayPlans"
        :key="plan.tier"
        class="pricing-card"
        :class="{ 'pricing-card--featured': plan.recommended }"
      >
        <span v-if="plan.recommended" class="pricing-card__tag">最受欢迎</span>

        <h3 class="pricing-card__name">{{ plan.name }}</h3>
        <div class="pricing-card__price">
          <span class="pricing-card__amount">{{ plan.priceLabel }}</span>
          <span v-if="plan.monthlyPrice" class="pricing-card__unit">/月起</span>
        </div>

        <ul class="pricing-card__features">
          <li v-for="(feat, idx) in plan.features" :key="idx">
            <ArtSvgIcon icon="ri:check-line" class="pricing-card__check" />
            <span>{{ feat }}</span>
          </li>
        </ul>

        <button
          type="button"
          class="pricing-card__btn"
          :class="{ 'pricing-card__btn--primary': plan.recommended }"
          @click="handleSelectPlan(plan.tier)"
        >
          选择此套餐
        </button>
      </article>
    </div>

    <section class="pricing-page__sla">
      <ComplianceSlaPanel
        title="服务时效承诺（SLA）"
        subtitle="签约后各环节交付时效，减少等待焦虑"
        :items="SERVICE_SLA_ITEMS"
      />
    </section>

    <div class="pricing-page__cta">
      <p>还不确定选哪个？先做免费合规诊断</p>
      <RouterLink to="/diagnosis" class="opc-btn opc-btn--secondary">
        <ArtSvgIcon icon="ri:shield-check-line" />
        免费诊断
      </RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { fetchServicePlans, type ServicePlan } from '@/api/frontend/compliance/diagnosis'
import { requireLogin } from '@/utils/auth/requireLogin'
import ComplianceSlaPanel from '@/components/frontend/ComplianceSlaPanel.vue'
import { SERVICE_SLA_ITEMS } from '@/data/frontend/complianceSla'

defineOptions({ name: 'CompliancePricing' })

const router = useRouter()

const loading = ref(true)
const plans = ref<ServicePlan[]>([])

const FALLBACK_PLANS: ServicePlan[] = [
  {
    id: 1,
    name: '基础版',
    tier: 'basic',
    monthlyPrice: 299,
    priceLabel: '¥299',
    features: ['主体注册代办', '月度记账', '季度申报', '年度汇算清缴']
  },
  {
    id: 2,
    name: '进阶版',
    tier: 'advanced',
    monthlyPrice: 799,
    priceLabel: '¥799',
    recommended: true,
    features: ['基础版全部服务', '税务筹划建议', '发票管理', '专属顾问答疑']
  },
  {
    id: 3,
    name: '尊享版',
    tier: 'premium',
    monthlyPrice: null,
    priceLabel: '面议',
    features: ['进阶版全部服务', '股权架构设计', '稽查应对支持', '一对一专属服务']
  }
]

const displayPlans = computed(() => (plans.value.length ? plans.value : FALLBACK_PLANS))

async function loadPlans() {
  loading.value = true
  try {
    const res = await fetchServicePlans()
    plans.value = res?.list?.length ? res.list : FALLBACK_PLANS
  } catch {
    plans.value = FALLBACK_PLANS
  } finally {
    loading.value = false
  }
}

function handleSelectPlan(tier: string) {
  const target = `/user/compliance/plan?plan=${tier}`
  if (!requireLogin({ redirect: target })) return
  router.push({ path: '/user/compliance/plan', query: { plan: tier } })
}

onMounted(loadPlans)
</script>

<style lang="scss" scoped>
.pricing-page {
  max-width: 960px;
  margin: 0 auto;
  padding: 40px 24px 80px;
}

.pricing-page__header {
  text-align: center;
  margin-bottom: 36px;
}

.pricing-page__badge {
  display: inline-block;
  margin-bottom: 16px;
  padding: 4px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 700;
  color: #2563eb;
  background: #eff6ff;
}

.pricing-page__title {
  margin: 0 0 12px;
  font-size: clamp(26px, 4vw, 36px);
  font-weight: 800;
  color: #1a1f36;
  letter-spacing: -0.02em;

  .is-accent {
    color: #2563eb;
  }
}

.pricing-page__subtitle {
  margin: 0 auto;
  max-width: 560px;
  font-size: 15px;
  line-height: 1.75;
  color: #6b7c93;
}

.pricing-page__loading {
  padding: 48px 16px;
  text-align: center;
  font-size: 14px;
  font-weight: 600;
  color: #94a3b8;
}

.pricing-page__grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

.pricing-card {
  position: relative;
  display: flex;
  flex-direction: column;
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

  &--featured {
    border-color: #bfdbfe;
    box-shadow: 0 4px 16px rgba(37, 99, 235, 0.12);

    &:hover {
      box-shadow: 0 8px 28px rgba(37, 99, 235, 0.16);
    }
  }
}

.pricing-card__tag {
  position: absolute;
  top: -10px;
  left: 50%;
  transform: translateX(-50%);
  padding: 4px 12px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 700;
  color: #fff;
  background: #2563eb;
  white-space: nowrap;
}

.pricing-card__name {
  margin: 8px 0 8px;
  font-size: 20px;
  font-weight: 700;
  color: #1a1f36;
}

.pricing-card__price {
  margin-bottom: 24px;
}

.pricing-card__amount {
  font-size: 32px;
  font-weight: 800;
  color: #1a1f36;
}

.pricing-card__unit {
  margin-left: 4px;
  font-size: 14px;
  font-weight: 600;
  color: #94a3b8;
}

.pricing-card__features {
  margin: 0 0 28px;
  padding: 0;
  list-style: none;
  flex: 1;

  li {
    display: flex;
    align-items: flex-start;
    gap: 10px;
    margin-bottom: 12px;
    font-size: 14px;
    line-height: 1.5;
    color: #475569;
  }
}

.pricing-card__check {
  flex-shrink: 0;
  margin-top: 2px;
  font-size: 16px;
  color: #16a34a;
}

.pricing-card__btn {
  width: 100%;
  padding: 12px 20px;
  border: 1px solid #d8dee9;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 600;
  color: #334155;
  background: #fff;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    background: #f8fafc;
    border-color: #94a3b8;
  }

  &--primary {
    color: #fff;
    background: #2563eb;
    border-color: #2563eb;
    box-shadow: 0 4px 14px rgba(37, 99, 235, 0.3);

    &:hover {
      background: #1d4ed8;
      border-color: #1d4ed8;
    }
  }
}

.pricing-page__sla {
  margin-bottom: 48px;
  padding: 28px;
  background: #fff;
  border: 1px solid #e8edf3;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);
}

.pricing-page__cta {
  margin-top: 40px;
  text-align: center;

  p {
    margin: 0 0 16px;
    font-size: 14px;
    color: #6b7c93;
    font-weight: 600;
  }
}

.opc-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 28px;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.2s ease;

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

@media (max-width: 768px) {
  .pricing-page {
    padding: 32px 16px 64px;
  }

  .pricing-page__grid {
    grid-template-columns: 1fr;
  }
}
</style>
