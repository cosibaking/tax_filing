<!-- +----------------------------------------------------------------------
  | XYGo Admin — 服务套餐 P-04
  +---------------------------------------------------------------------- -->
<template>
  <main class="pt-8 pb-16 px-6">
    <div class="max-w-6xl mx-auto">
      <!-- Header -->
      <div class="text-center mb-16">
        <div class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-white/50 border border-white/50 shadow-sm mb-6">
          <ArtSvgIcon icon="ri:vip-crown-line" class="text-lg text-clay-accent" />
          <span class="text-sm font-bold text-clay-muted">OPC 合规服务</span>
        </div>
        <h1 class="font-heading font-black text-4xl md:text-5xl text-clay-foreground mb-4">
          选择适合您的 <span class="text-clay-accent">服务套餐</span>
        </h1>
        <p class="text-lg text-clay-muted max-w-2xl mx-auto">
          从 OPC 设立到日常记账申报，全程代办，让您专注内容创作
        </p>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="text-center py-20">
        <ArtSvgIcon icon="ri:loader-4-line" class="text-4xl text-clay-accent animate-spin mx-auto mb-4" />
        <p class="text-clay-muted font-bold">加载套餐信息...</p>
      </div>

      <!-- Pricing cards -->
      <div v-else class="grid md:grid-cols-3 gap-8">
        <div
          v-for="plan in displayPlans"
          :key="plan.tier"
          class="relative rounded-[48px] p-8 md:p-10 transition-all duration-500 hover:-translate-y-2 border"
          :class="plan.recommended
            ? 'bg-white/90 border-blue-200 shadow-clay-deep ring-2 ring-blue-200/40'
            : 'bg-white/70 border-[#d1d9e6]/40 shadow-clay-card hover:shadow-clay-card-hover'"
        >
          <div v-if="plan.recommended" class="absolute -top-4 left-1/2 -translate-x-1/2 px-5 py-1.5 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 text-white text-xs font-black shadow-clay-btn whitespace-nowrap">
            最受欢迎
          </div>

          <h3 class="font-heading font-black text-2xl text-clay-foreground mb-2 mt-2">{{ plan.name }}</h3>
          <div class="mb-6">
            <span class="font-heading font-black text-4xl text-clay-foreground">{{ plan.priceLabel }}</span>
            <span v-if="plan.monthlyPrice" class="text-clay-muted font-bold ml-1">/月起</span>
          </div>

          <ul class="space-y-4 mb-10">
            <li v-for="(feat, idx) in plan.features" :key="idx" class="flex items-start gap-3">
              <div class="w-5 h-5 rounded-full bg-clay-success/20 text-clay-success flex items-center justify-center shrink-0 mt-0.5">
                <ArtSvgIcon icon="ri:check-line" class="text-xs" />
              </div>
              <span class="text-sm font-medium text-clay-foreground">{{ feat }}</span>
            </li>
          </ul>

          <button
            type="button"
            class="w-full py-4 rounded-2xl font-black text-base transition-all duration-300 active:scale-95"
            :class="plan.recommended
              ? 'bg-gradient-to-br from-blue-400 to-blue-600 text-white shadow-clay-btn hover:shadow-clay-btn-hover hover:-translate-y-1'
              : 'bg-white text-clay-foreground shadow-clay-btn hover:shadow-clay-btn-hover'"
            @click="handleSelectPlan(plan.tier)"
          >
            选择此套餐
          </button>
        </div>
      </div>

      <!-- Free diagnosis CTA -->
      <div class="mt-16 text-center">
        <p class="text-clay-muted font-medium mb-4">还不确定选哪个？先做免费合规诊断</p>
        <RouterLink
          to="/diagnosis"
          class="inline-flex items-center gap-2 px-8 py-3 rounded-2xl bg-white shadow-clay-btn hover:shadow-clay-btn-hover font-bold text-clay-accent transition-all"
        >
          <ArtSvgIcon icon="ri:shield-check-line" class="text-lg" />
          免费诊断
        </RouterLink>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { fetchServicePlans, type ServicePlan } from '@/api/frontend/compliance/diagnosis'
import { useMemberStore } from '@/store/modules/member'

defineOptions({ name: 'CompliancePricing' })

const router = useRouter()
const memberStore = useMemberStore()

const loading = ref(true)
const plans = ref<ServicePlan[]>([])

const FALLBACK_PLANS: ServicePlan[] = [
  {
    id: 1,
    name: '基础版',
    tier: 'basic',
    monthlyPrice: 299,
    priceLabel: '¥299',
    features: ['OPC 注册代办', '月度记账', '季度申报', '年度汇算清缴']
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
  const query = { plan: tier }
  if (memberStore.isLogin) {
    router.push({ path: '/user/compliance/plan', query })
  } else {
    router.push({
      path: '/user/login',
      query: { redirect: `/user/compliance/plan?plan=${tier}` }
    })
  }
}

onMounted(loadPlans)
</script>

<style lang="scss" scoped>
.text-clay-foreground { color: #32325d; }
.text-clay-muted { color: #8898aa; }
.text-clay-accent { color: #5a8dee; }
.text-clay-success { color: #71dd37; }
.font-heading { font-family: 'Nunito', 'PingFang SC', sans-serif; }

.shadow-clay-deep {
  box-shadow: 30px 30px 60px #d1d9e6, -30px -30px 60px #ffffff,
    inset 10px 10px 20px rgba(90, 141, 238, 0.05), inset -10px -10px 20px rgba(255, 255, 255, 0.8);
}
.shadow-clay-card {
  box-shadow: 16px 16px 32px rgba(165, 175, 190, 0.3), -10px -10px 24px rgba(255, 255, 255, 0.9),
    inset 6px 6px 12px rgba(90, 141, 238, 0.03), inset -6px -6px 12px rgba(255, 255, 255, 1);
}
.shadow-clay-card-hover {
  box-shadow: 20px 20px 40px rgba(165, 175, 190, 0.35), -12px -12px 30px rgba(255, 255, 255, 0.95),
    inset 6px 6px 12px rgba(90, 141, 238, 0.03), inset -6px -6px 12px rgba(255, 255, 255, 1);
}
.shadow-clay-btn {
  box-shadow: 12px 12px 24px rgba(90, 141, 238, 0.3), -8px -8px 16px rgba(255, 255, 255, 0.4),
    inset 4px 4px 8px rgba(255, 255, 255, 0.4), inset -4px -4px 8px rgba(0, 0, 0, 0.05);
}
.shadow-clay-btn-hover {
  box-shadow: 16px 16px 32px rgba(90, 141, 238, 0.4), -10px -10px 20px rgba(255, 255, 255, 0.5),
    inset 4px 4px 8px rgba(255, 255, 255, 0.4), inset -4px -4px 8px rgba(0, 0, 0, 0.05);
}
</style>
