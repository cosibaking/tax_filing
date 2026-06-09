<!-- +----------------------------------------------------------------------
  | XYGo Admin — 会员合规诊断历史
  +---------------------------------------------------------------------- -->
<template>
  <section class="bg-white/70 backdrop-blur-xl rounded-[48px] shadow-clay-card border border-[#d1d9e6]/40 p-8 md:p-10 animate-in">
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
      <div>
        <h2 class="font-heading font-black text-2xl text-clay-foreground">诊断历史</h2>
        <p class="text-xs text-clay-muted mt-1">查看您的合规诊断记录与推荐方案</p>
      </div>
      <RouterLink
        to="/diagnosis"
        class="px-6 py-3 rounded-2xl bg-gradient-to-br from-blue-400 to-blue-600 text-white shadow-clay-btn font-bold text-sm active:scale-95 transition-all flex items-center gap-2"
      >
        <ArtSvgIcon icon="ri:add-line" class="text-lg" />
        新建诊断
      </RouterLink>
    </div>

    <div v-if="loading && list.length === 0" class="py-20 text-center">
      <ArtSvgIcon icon="ri:loader-4-line" class="text-4xl text-clay-muted animate-spin mx-auto mb-4" />
      <p class="text-clay-muted font-medium">加载中...</p>
    </div>

    <div v-else-if="list.length === 0" class="py-20 text-center">
      <div class="w-24 h-24 rounded-full bg-white shadow-clay-btn flex items-center justify-center mb-6 mx-auto">
        <ArtSvgIcon icon="ri:file-search-line" class="text-[40px] text-clay-muted opacity-50" />
      </div>
      <p class="text-clay-muted font-bold text-lg mb-4">暂无诊断记录</p>
      <RouterLink to="/diagnosis" class="text-clay-accent font-bold hover:underline">立即开始免费诊断</RouterLink>
    </div>

    <div v-else class="space-y-4">
      <div
        v-for="item in list"
        :key="item.id"
        class="p-6 rounded-[32px] bg-[#f0f3f8] shadow-clay-pressed hover:bg-white hover:shadow-clay-card transition-all duration-300"
      >
        <div class="flex flex-wrap items-center gap-3 mb-4">
          <span class="px-3 py-1 rounded-lg bg-white shadow-clay-btn text-[10px] font-black text-clay-accent uppercase">
            {{ planLabel(item.recommendedPlan) }}
          </span>
          <span class="text-xs text-clay-muted font-medium">{{ item.createdAt }}</span>
        </div>
        <div class="grid sm:grid-cols-2 md:grid-cols-4 gap-4 text-sm">
          <div>
            <span class="block text-[10px] font-black text-clay-muted uppercase tracking-widest mb-1">月收入区间</span>
            <span class="font-bold text-clay-foreground">{{ item.monthlyIncomeRange }}</span>
          </div>
          <div>
            <span class="block text-[10px] font-black text-clay-muted uppercase tracking-widest mb-1">年成本估算</span>
            <span class="font-bold text-clay-foreground">{{ formatMoney(item.annualCostEstimate) }} 元</span>
          </div>
          <div>
            <span class="block text-[10px] font-black text-clay-muted uppercase tracking-widest mb-1">现有主体</span>
            <span class="font-bold text-clay-foreground">{{ entityLabel(item.existingEntity) }}</span>
          </div>
          <div>
            <span class="block text-[10px] font-black text-clay-muted uppercase tracking-widest mb-1">收入渠道</span>
            <span class="font-bold text-clay-foreground">{{ item.platforms?.join('、') || '-' }}</span>
          </div>
        </div>
        <div v-if="item.taxComparison?.items?.length" class="mt-4 pt-4 border-t border-gray-200/50">
          <span class="text-[10px] font-black text-clay-muted uppercase tracking-widest">推荐方案预估年税负</span>
          <div class="flex flex-wrap gap-3 mt-2">
            <span
              v-for="tax in item.taxComparison.items.filter(t => t.recommended)"
              :key="tax.plan"
              class="px-3 py-1.5 rounded-xl bg-white shadow-clay-btn text-xs font-bold text-clay-success"
            >
              {{ tax.label }}：{{ formatMoney(tax.annualTax) }} 元
            </span>
          </div>
        </div>
        <div class="mt-4 flex justify-end">
          <RouterLink
            :to="`/diagnosis/result?id=${item.id}`"
            class="text-sm font-bold text-clay-accent hover:underline flex items-center gap-1"
          >
            查看详情
            <ArtSvgIcon icon="ri:arrow-right-s-line" />
          </RouterLink>
        </div>
      </div>

      <div v-if="total > pageSize" class="pt-4 flex justify-center">
        <button
          v-if="page * pageSize < total"
          class="flex items-center gap-2 text-sm font-bold text-clay-muted hover:text-clay-accent transition-colors"
          :disabled="loading"
          @click="loadMore"
        >
          <ArtSvgIcon v-if="loading" icon="ri:loader-4-line" class="animate-spin" />
          加载更多 <ArtSvgIcon icon="ri:arrow-down-s-line" class="text-base" />
        </button>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { getDiagnosisHistory, type DiagnosisHistoryItem } from '@/api/frontend/compliance/member'
import type { RecommendedPlan } from '@/api/frontend/compliance/diagnosis'
import { syncGuestDiagnosisAfterLogin } from '@/utils/compliance/syncGuestDiagnosis'

defineOptions({ name: 'ComplianceDiagnosisHistory' })

const list = ref<DiagnosisHistoryItem[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = 20
const total = ref(0)

const planLabel = (plan: RecommendedPlan | string) => {
  const map: Record<string, string> = {
    opc: 'OPC 方案',
    individual: '个体户',
    labor: '劳务报酬',
    transitional: '过渡方案',
    none: '暂无推荐'
  }
  return map[plan] || plan
}

const entityLabel = (entity: string) => {
  const map: Record<string, string> = {
    none: '无',
    individual: '个体户',
    company: '公司',
    other: '其他'
  }
  return map[entity] || entity
}

const formatMoney = (v: number) => (Number(v) || 0).toLocaleString('zh-CN', { minimumFractionDigits: 0, maximumFractionDigits: 2 })

async function loadHistory(p = 1) {
  loading.value = true
  try {
    const res = await getDiagnosisHistory({ page: p, pageSize })
    if (res) {
      if (p === 1) {
        list.value = res.list || []
      } else {
        list.value.push(...(res.list || []))
      }
      page.value = res.page
      total.value = res.total
    }
  } catch { /* 拦截器已处理 */ } finally {
    loading.value = false
  }
}

function loadMore() {
  loadHistory(page.value + 1)
}

onMounted(async () => {
  await syncGuestDiagnosisAfterLogin()
  loadHistory(1)
})
</script>

<style lang="scss" scoped>
.text-clay-foreground { color: #32325d; }
.text-clay-muted { color: #8898aa; }
.text-clay-accent { color: #5a8dee; }
.text-clay-success { color: #71dd37; }
.font-heading { font-family: 'Nunito', 'PingFang SC', sans-serif; }

.shadow-clay-card {
  box-shadow: 16px 16px 32px rgba(165, 175, 190, 0.3), -10px -10px 24px rgba(255, 255, 255, 0.9),
    inset 6px 6px 12px rgba(90, 141, 238, 0.03), inset -6px -6px 12px rgba(255, 255, 255, 1);
}
.shadow-clay-btn {
  box-shadow: 12px 12px 24px rgba(90, 141, 238, 0.3), -8px -8px 16px rgba(255, 255, 255, 0.4),
    inset 4px 4px 8px rgba(255, 255, 255, 0.4), inset -4px -4px 8px rgba(0, 0, 0, 0.05);
}
.shadow-clay-pressed {
  box-shadow: inset 10px 10px 20px #e0e5ec, inset -10px -10px 20px #ffffff;
}

.animate-in {
  animation: slideIn 0.4s ease-out;
}
@keyframes slideIn {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
