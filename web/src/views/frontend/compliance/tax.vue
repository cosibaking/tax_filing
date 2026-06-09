<!-- P-11 申报日历 -->
<template>
  <div>
    <EmploymentStatusCard only-when-unknown class="mb-6" @saved="loadChecklist" />

    <section class="bg-white/70 backdrop-blur-xl rounded-[48px] shadow-clay-card border border-[#d1d9e6]/40 p-8 md:p-10">
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
      <div>
        <h2 class="font-heading font-black text-2xl text-clay-foreground">申报管理</h2>
        <p class="text-xs text-clay-muted mt-1">申报截止日与任务状态</p>
      </div>
      <div v-if="calendar?.nextDueDate" class="px-4 py-2 rounded-2xl bg-orange-50 border border-orange-200 text-sm font-bold text-orange-700">
        下次截止：{{ calendar.nextDueDate }}
        <span v-if="calendar.daysUntilDue != null">（{{ calendar.daysUntilDue }} 天）</span>
      </div>
    </div>

    <!-- 月历导航 -->
    <div class="flex items-center justify-between mb-6">
      <button type="button" class="p-2 rounded-xl bg-[#f0f3f8] shadow-clay-pressed" @click="prevMonth">
        <ArtSvgIcon icon="ri:arrow-left-s-line" class="text-xl text-clay-accent" />
      </button>
      <span class="font-heading font-black text-lg text-clay-foreground">{{ viewYear }}年{{ viewMonth }}月</span>
      <button type="button" class="p-2 rounded-xl bg-[#f0f3f8] shadow-clay-pressed" @click="nextMonth">
        <ArtSvgIcon icon="ri:arrow-right-s-line" class="text-xl text-clay-accent" />
      </button>
    </div>

    <div v-if="loading" class="py-16 text-center">
      <ArtSvgIcon icon="ri:loader-4-line" class="text-3xl text-clay-accent animate-spin mx-auto" />
    </div>

    <template v-else>
      <!-- 月历格子 -->
      <div class="grid grid-cols-7 gap-1 mb-8 text-center text-xs">
        <div v-for="d in weekDays" :key="d" class="py-2 font-black text-clay-muted">{{ d }}</div>
        <div v-for="blank in firstDayOffset" :key="'b' + blank" />
        <div
          v-for="day in daysInMonth"
          :key="day"
          class="py-2 rounded-lg font-medium relative"
          :class="isDueDay(day) ? 'bg-red-50 text-red-600 font-black' : 'text-clay-foreground'"
        >
          {{ day }}
          <span v-if="isDueDay(day)" class="absolute bottom-0.5 left-1/2 -translate-x-1/2 w-1.5 h-1.5 rounded-full bg-red-500" />
        </div>
      </div>

      <!-- 本年任务列表 -->
      <div class="mb-8">
        <h3 class="text-sm font-black text-clay-muted uppercase tracking-widest mb-4">本年任务列表</h3>
        <div v-if="tasks.length === 0" class="text-sm text-clay-muted py-4">本年暂无申报任务</div>
        <div v-else class="space-y-3">
          <div
            v-for="task in tasks"
            :key="task.id"
            class="p-4 rounded-2xl bg-[#f0f3f8] shadow-clay-pressed flex flex-wrap items-center justify-between gap-3"
          >
            <div>
              <span class="font-bold text-clay-foreground">{{ formatTaskTitle(task) }}</span>
              <span class="ml-3 text-xs font-bold px-2 py-0.5 rounded-full" :class="statusClass(task.status)">
                {{ TAX_STATUS_LABELS[task.status] }}
              </span>
              <p class="text-xs text-clay-muted mt-1">预估 ¥{{ formatMoney(task.calculatedAmount) }}</p>
            </div>
            <button
              type="button"
              class="text-sm font-bold text-clay-accent hover:underline"
              @click="openTaskDetail(task)"
            >
              {{ task.status === 'filed' ? '回执' : '详情' }}
            </button>
          </div>
        </div>
      </div>

      <!-- 自查清单 -->
      <div class="p-6 rounded-[32px] bg-[#f0f3f8] shadow-clay-pressed">
        <h3 class="text-sm font-black text-clay-muted uppercase tracking-widest mb-4">报税前自查清单</h3>
        <div class="space-y-3">
          <label
            v-for="item in checklist"
            :key="item.key"
            class="flex items-start gap-3 text-sm font-medium"
            :class="item.na ? 'text-clay-muted' : 'text-clay-foreground'"
          >
            <ElCheckbox :model-value="item.checked" disabled class="mt-0.5" />
            <span>
              {{ item.label }}
              <span v-if="item.na" class="text-xs text-clay-muted">（不适用）</span>
              <span v-else-if="item.hint" class="block text-xs text-clay-muted mt-0.5">{{ item.hint }}</span>
            </span>
          </label>
        </div>
        <p class="text-xs text-clay-muted mt-4">由顾问在申报前勾选确认，此处为只读展示</p>
      </div>
    </template>

    <!-- 任务详情 Drawer -->
    <ElDrawer v-model="drawerVisible" :title="drawerTask?.taxTypeLabel || '申报详情'" size="400px">
      <template v-if="drawerTask">
        <div class="space-y-4 text-sm">
          <div class="flex justify-between">
            <span class="text-clay-muted">税种</span>
            <span class="font-bold">{{ drawerTask.taxTypeLabel || TAX_TYPE_LABELS[drawerTask.taxType] }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-clay-muted">申报周期</span>
            <span class="font-bold">{{ drawerTask.period }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-clay-muted">截止日</span>
            <span class="font-bold">{{ drawerTask.dueDate }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-clay-muted">状态</span>
            <span class="font-bold">{{ TAX_STATUS_LABELS[drawerTask.status] }}</span>
          </div>

          <div v-if="drawerTask.detail" class="p-4 rounded-xl bg-gray-50 space-y-2">
            <h4 class="font-black text-clay-muted text-xs uppercase">计算明细</h4>
            <div class="flex justify-between"><span>不含税收入</span><span>¥{{ formatMoney(drawerTask.detail.revenueExTax) }}</span></div>
            <div class="flex justify-between"><span>增值税</span><span>¥{{ formatMoney(drawerTask.detail.vat) }}</span></div>
            <div class="flex justify-between"><span>附加税</span><span>¥{{ formatMoney(drawerTask.detail.surcharge) }}</span></div>
            <div class="flex justify-between"><span>企税</span><span>¥{{ formatMoney(drawerTask.detail.cit) }}</span></div>
            <div class="flex justify-between font-bold border-t pt-2"><span>合计</span><span>¥{{ formatMoney(drawerTask.detail.total) }}</span></div>
          </div>

          <a
            v-if="drawerTask.status === 'filed' && drawerTask.receiptUrl"
            :href="drawerTask.receiptUrl"
            target="_blank"
            class="inline-flex items-center gap-2 text-clay-accent font-bold hover:underline"
          >
            <ArtSvgIcon icon="ri:file-pdf-line" />下载回执 PDF
          </a>
        </div>
      </template>
    </ElDrawer>
    </section>
  </div>
</template>

<script setup lang="ts">
import EmploymentStatusCard from '@/components/member/EmploymentStatusCard.vue'
import {
  getTaxCalendar,
  getTaxChecklist,
  getTaxTaskDetail,
  TAX_CHECKLIST_ITEMS,
  TAX_TYPE_LABELS,
  TAX_STATUS_LABELS,
  formatMoney,
  type TaxTask,
  type TaxChecklistItem
} from '@/api/frontend/compliance/member'

defineOptions({ name: 'ComplianceTax' })

const loading = ref(false)
const calendar = ref<Awaited<ReturnType<typeof getTaxCalendar>> | null>(null)
const tasks = ref<TaxTask[]>([])
const checklist = ref<TaxChecklistItem[]>([])

const now = new Date()
const viewYear = ref(now.getFullYear())
const viewMonth = ref(now.getMonth() + 1)

const weekDays = ['日', '一', '二', '三', '四', '五', '六']
const daysInMonth = computed(() => new Date(viewYear.value, viewMonth.value, 0).getDate())
const firstDayOffset = computed(() => new Date(viewYear.value, viewMonth.value - 1, 1).getDay())
const dueDateSet = computed(() => new Set((calendar.value?.dueDates || []).map(d => parseInt(d.split('-')[2] || '0', 10))))

const drawerVisible = ref(false)
const drawerTask = ref<TaxTask | null>(null)

function isDueDay(day: number) {
  return dueDateSet.value.has(day)
}

function formatTaskTitle(task: TaxTask) {
  const label = task.taxTypeLabel || TAX_TYPE_LABELS[task.taxType]
  return task.period ? `${task.period}${label}` : label
}

function statusClass(status: string) {
  if (status === 'filed') return 'bg-green-100 text-green-700'
  if (status === 'overdue') return 'bg-red-100 text-red-700'
  return 'bg-blue-100 text-clay-accent'
}

function prevMonth() {
  if (viewMonth.value === 1) { viewYear.value--; viewMonth.value = 12 }
  else viewMonth.value--
  loadData()
}

function nextMonth() {
  if (viewMonth.value === 12) { viewYear.value++; viewMonth.value = 1 }
  else viewMonth.value++
  loadData()
}

async function loadChecklist() {
  try {
    const period = `${viewYear.value}-${String(viewMonth.value).padStart(2, '0')}`
    const cl = await getTaxChecklist({ period })
    checklist.value = cl.items?.length
      ? cl.items
      : TAX_CHECKLIST_ITEMS.map(i => ({ ...i, checked: false }))
  } catch {
    checklist.value = TAX_CHECKLIST_ITEMS.map(i => ({ ...i, checked: false }))
  }
}

async function loadData() {
  loading.value = true
  try {
    calendar.value = await getTaxCalendar({ year: viewYear.value, month: viewMonth.value })
    tasks.value = calendar.value?.tasks || []
    await loadChecklist()
  } catch {
    calendar.value = null
    tasks.value = []
    checklist.value = TAX_CHECKLIST_ITEMS.map(i => ({ ...i, checked: false }))
  } finally {
    loading.value = false
  }
}

async function openTaskDetail(task: TaxTask) {
  try {
    drawerTask.value = await getTaxTaskDetail(task.id)
  } catch {
    drawerTask.value = task
  }
  drawerVisible.value = true
}

onMounted(loadData)
</script>

<style lang="scss" scoped>
.text-clay-foreground { color: #32325d; }
.text-clay-muted { color: #8898aa; }
.text-clay-accent { color: #5a8dee; }
.font-heading { font-family: 'Nunito', 'PingFang SC', sans-serif; }
.shadow-clay-card { box-shadow: 16px 16px 32px rgba(165, 175, 190, 0.3), -10px -10px 24px rgba(255, 255, 255, 0.9); }
.shadow-clay-pressed { box-shadow: inset 10px 10px 20px #e0e5ec, inset -10px -10px 20px #ffffff; }
</style>
