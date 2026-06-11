<!-- P-10 账套与利润表 -->
<template>
  <section class="member-panel">
    <div class="member-panel__head">
      <div>
        <h2 class="member-panel__title">利润报表</h2>
        <p class="member-panel__desc">自动生成账套摘要与利润表</p>
      </div>
      <div class="member-tabs">
        <button
          v-for="t in periodTypes"
          :key="t.value"
          type="button"
          class="member-tab"
          :class="{ 'is-active': periodType === t.value }"
          @click="periodType = t.value; loadData()"
        >
          {{ t.label }}
        </button>
      </div>
    </div>

    <div class="mb-6 member-field">
      <ElDatePicker
        v-if="periodType === 'month'"
        v-model="selectedPeriod"
        type="month"
        value-format="YYYY-MM"
        @change="loadData"
      />
      <ElDatePicker
        v-else-if="periodType === 'quarter'"
        v-model="selectedPeriod"
        type="month"
        value-format="YYYY-MM"
        placeholder="选择季度内任一月"
        @change="loadData"
      />
      <ElDatePicker
        v-else
        v-model="selectedPeriod"
        type="year"
        value-format="YYYY"
        @change="loadData"
      />
    </div>

    <div v-if="loading" class="member-empty">
      <ArtSvgIcon icon="ri:loader-4-line" class="text-3xl text-clay-accent animate-spin mx-auto" />
    </div>

    <div v-else-if="profit" class="space-y-6">
      <!-- 利润表简版 -->
      <div class="p-6 rounded-lg border border-[#e8edf3] bg-[#f8fafc] space-y-4">
        <h3 class="text-sm font-black text-clay-muted uppercase tracking-widest">利润表（简版）</h3>
        <div class="flex justify-between text-sm">
          <span class="text-clay-muted font-medium">营业收入（不含税）</span>
          <span class="font-black text-clay-foreground">¥{{ formatMoney(profit.revenue) }}</span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-clay-muted font-medium">营业成本</span>
          <span class="font-black text-clay-foreground">¥{{ formatMoney(profit.cost) }}</span>
        </div>
        <div class="border-t border-gray-300/50 pt-4 flex justify-between">
          <span class="font-black text-clay-foreground">利润总额</span>
          <span class="font-black text-xl" :class="profit.profit >= 0 ? 'text-clay-success' : 'text-red-500'">
            ¥{{ formatMoney(profit.profit) }}
          </span>
        </div>
        <div class="flex justify-between text-sm">
          <span class="text-clay-muted font-medium">累计年度利润</span>
          <span class="font-bold text-clay-accent">¥{{ formatMoney(profit.cumulativeProfit) }}</span>
        </div>
      </div>

      <!-- 会计分录 -->
      <div>
        <button
          type="button"
          class="flex items-center gap-2 text-sm font-bold text-clay-accent hover:underline mb-4"
          @click="showVouchers = !showVouchers"
        >
          <ArtSvgIcon :icon="showVouchers ? 'ri:arrow-up-s-line' : 'ri:arrow-down-s-line'" />
          {{ showVouchers ? '收起' : '展开' }}本月会计分录（只读）
        </button>

        <div v-if="showVouchers" class="member-table-wrap">
          <table v-if="vouchers.length" class="member-table">
            <thead>
              <tr>
                <th class="py-3 px-4">日期</th>
                <th class="py-3 px-4">摘要</th>
                <th class="py-3 px-4">借方</th>
                <th class="py-3 px-4">贷方</th>
                <th class="py-3 px-4">金额</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="v in vouchers" :key="v.id" class="border-b border-gray-100">
                <td class="py-3 px-4">{{ v.occurredAt }}</td>
                <td class="py-3 px-4">{{ v.summary }}</td>
                <td class="py-3 px-4">{{ v.debitAccount }}</td>
                <td class="py-3 px-4">{{ v.creditAccount }}</td>
                <td class="py-3 px-4 font-bold">¥{{ formatMoney(v.amount) }}</td>
              </tr>
            </tbody>
          </table>
          <p v-else class="text-sm text-clay-muted py-8 text-center">暂无分录记录</p>
        </div>
      </div>
    </div>

    <div v-else class="py-16 text-center text-clay-muted font-medium space-y-2">
      <p>暂无 {{ selectedPeriod }} 的利润数据</p>
      <p class="text-sm">请先在「收入台账」或「费用台账」录入<strong>发生日期落在该月</strong>的记录</p>
    </div>
  </section>
</template>

<script setup lang="ts">
import {
  getProfitSummary,
  getLedgerVouchers,
  formatMoney,
  type ProfitSummary,
  type ProfitPeriodType,
  type LedgerVoucher
} from '@/api/frontend/compliance/member'

defineOptions({ name: 'ComplianceLedger' })

const loading = ref(false)
const profit = ref<ProfitSummary | null>(null)
const vouchers = ref<LedgerVoucher[]>([])
const showVouchers = ref(false)

const periodTypes: { value: ProfitPeriodType; label: string }[] = [
  { value: 'month', label: '月' },
  { value: 'quarter', label: '季' },
  { value: 'year', label: '年' }
]
const periodType = ref<ProfitPeriodType>('month')

const now = new Date()
const selectedPeriod = ref(`${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`)

function voucherPeriod(): string {
  if (periodType.value === 'year') return selectedPeriod.value
  return selectedPeriod.value.slice(0, 7)
}

function isEmptyProfit(data: ProfitSummary) {
  return !Number(data.revenue) && !Number(data.cost) && !Number(data.profit)
}

async function loadData() {
  loading.value = true
  try {
    const data = await getProfitSummary({
      period: selectedPeriod.value,
      periodType: periodType.value
    })
    profit.value = isEmptyProfit(data) ? null : data
    const vRes = await getLedgerVouchers({ period: voucherPeriod() })
    vouchers.value = vRes.list || []
  } catch {
    profit.value = null
    vouchers.value = []
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>

<style lang="scss" scoped>
.text-clay-foreground { color: #32325d; }
.text-clay-muted { color: #8898aa; }
.text-clay-accent { color: #2563eb; }
.text-clay-success { color: #16a34a; }
</style>
