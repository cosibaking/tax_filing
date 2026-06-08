<!-- P-12 月度对账单 -->
<template>
  <section class="bg-white/70 backdrop-blur-xl rounded-[48px] shadow-clay-card border border-[#d1d9e6]/40 p-8 md:p-10">
    <div class="mb-8">
      <h2 class="font-heading font-black text-2xl text-clay-foreground">月度对账单</h2>
      <p class="text-xs text-clay-muted mt-1">查看历史对账单与服务完成状态</p>
    </div>

    <div v-if="loading" class="py-16 text-center">
      <ArtSvgIcon icon="ri:loader-4-line" class="text-3xl text-clay-accent animate-spin mx-auto" />
    </div>

    <div v-else-if="list.length === 0" class="py-16 text-center text-clay-muted font-medium">
      暂无对账单，每月 16 日后由顾问生成
    </div>

    <template v-else>
      <div class="overflow-x-auto rounded-2xl mb-8">
        <table class="w-full text-sm">
          <thead>
            <tr class="text-left text-[10px] font-black text-clay-muted uppercase tracking-widest border-b">
              <th class="py-3 px-4">月份</th>
              <th class="py-3 px-4">收入</th>
              <th class="py-3 px-4">成本</th>
              <th class="py-3 px-4">利润</th>
              <th class="py-3 px-4">预缴税额</th>
              <th class="py-3 px-4">状态</th>
              <th class="py-3 px-4">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="row in list"
              :key="row.id"
              class="border-b border-gray-100 cursor-pointer transition-colors"
              :class="selected?.id === row.id ? 'bg-blue-50' : 'hover:bg-white/50'"
              @click="selected = row"
            >
              <td class="py-3 px-4 font-bold text-clay-foreground">{{ row.period }}</td>
              <td class="py-3 px-4">¥{{ formatMoney(row.revenue) }}</td>
              <td class="py-3 px-4">¥{{ formatMoney(row.cost) }}</td>
              <td class="py-3 px-4 font-bold" :class="row.profit >= 0 ? 'text-clay-success' : 'text-red-500'">
                ¥{{ formatMoney(row.profit) }}
              </td>
              <td class="py-3 px-4">¥{{ formatMoney(row.prepaidTax) }}</td>
              <td class="py-3 px-4">
                <span class="text-xs font-bold px-2 py-0.5 rounded-full" :class="row.status === 'completed' ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-clay-muted'">
                  {{ row.status === 'completed' ? '✅ 已完成' : '处理中' }}
                </span>
              </td>
              <td class="py-3 px-4">
                <button
                  type="button"
                  class="text-sm font-bold text-clay-accent hover:underline flex items-center gap-1"
                  @click.stop="downloadPdf(row)"
                >
                  <ArtSvgIcon icon="ri:file-pdf-line" />PDF
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 详情卡片 -->
      <div v-if="selected" class="p-6 rounded-[32px] bg-[#f0f3f8] shadow-clay-pressed">
        <h3 class="font-heading font-black text-lg text-clay-foreground mb-4">{{ selected.period }} 服务完成摘要</h3>
        <div class="grid sm:grid-cols-2 gap-4 text-sm">
          <div class="p-4 rounded-2xl bg-white shadow-clay-btn">
            <span class="block text-[10px] font-black text-clay-muted uppercase mb-1">本月收入</span>
            <span class="font-black text-xl text-clay-foreground">¥{{ formatMoney(selected.revenue) }}</span>
          </div>
          <div class="p-4 rounded-2xl bg-white shadow-clay-btn">
            <span class="block text-[10px] font-black text-clay-muted uppercase mb-1">本月成本</span>
            <span class="font-black text-xl text-clay-foreground">¥{{ formatMoney(selected.cost) }}</span>
          </div>
          <div class="p-4 rounded-2xl bg-white shadow-clay-btn">
            <span class="block text-[10px] font-black text-clay-muted uppercase mb-1">本月利润</span>
            <span class="font-black text-xl" :class="selected.profit >= 0 ? 'text-clay-success' : 'text-red-500'">
              ¥{{ formatMoney(selected.profit) }}
            </span>
          </div>
          <div class="p-4 rounded-2xl bg-white shadow-clay-btn">
            <span class="block text-[10px] font-black text-clay-muted uppercase mb-1">本月预缴税额</span>
            <span class="font-black text-xl text-clay-accent">¥{{ formatMoney(selected.prepaidTax) }}</span>
          </div>
        </div>
        <div class="mt-4 pt-4 border-t border-gray-200/50 flex flex-wrap gap-4 text-sm font-medium text-clay-muted">
          <span>累计年度利润：¥{{ formatMoney(selected.cumulativeProfit) }}</span>
          <span>申报状态：{{ selected.filingStatus || '已申报完成' }}</span>
        </div>
      </div>
    </template>
  </section>
</template>

<script setup lang="ts">
import {
  getStatementList,
  getStatementPdfUrl,
  formatMoney,
  type StatementItem
} from '@/api/frontend/compliance/member'
import { ElMessage } from 'element-plus'

defineOptions({ name: 'ComplianceStatement' })

const loading = ref(false)
const list = ref<StatementItem[]>([])
const selected = ref<StatementItem | null>(null)

async function loadData() {
  loading.value = true
  try {
    const res = await getStatementList({ page: 1, pageSize: 24 })
    list.value = res.list || []
    if (list.value.length && !selected.value) {
      selected.value = list.value[0]
    }
  } catch {
    list.value = []
  } finally {
    loading.value = false
  }
}

async function downloadPdf(row: StatementItem) {
  if (row.pdfUrl) {
    window.open(row.pdfUrl, '_blank')
    return
  }
  try {
    const res = await getStatementPdfUrl(row.id)
    if (res.url) window.open(res.url, '_blank')
    else ElMessage.warning('PDF 尚未生成')
  } catch {
    ElMessage.error('下载失败，请稍后重试')
  }
}

onMounted(loadData)
</script>

<style lang="scss" scoped>
.text-clay-foreground { color: #32325d; }
.text-clay-muted { color: #8898aa; }
.text-clay-accent { color: #5a8dee; }
.text-clay-success { color: #71dd37; }
.font-heading { font-family: 'Nunito', 'PingFang SC', sans-serif; }
.shadow-clay-card { box-shadow: 16px 16px 32px rgba(165, 175, 190, 0.3), -10px -10px 24px rgba(255, 255, 255, 0.9); }
.shadow-clay-btn { box-shadow: 12px 12px 24px rgba(90, 141, 238, 0.3), -8px -8px 16px rgba(255, 255, 255, 0.4); }
.shadow-clay-pressed { box-shadow: inset 10px 10px 20px #e0e5ec, inset -10px -10px 20px #ffffff; }
</style>
