<!-- P-11 申报日历 -->
<template>
  <div>
    <EmploymentStatusCard only-when-unknown class="mb-6" @saved="loadChecklist" />

    <section class="member-panel">
    <div class="member-panel__head">
      <div>
        <h2 class="member-panel__title">申报管理</h2>
        <p class="member-panel__desc">申报截止日与任务状态</p>
      </div>
      <div v-if="calendar?.nextDueDate" class="member-alert member-alert--info">
        下次截止：{{ calendar.nextDueDate }}
        <span v-if="calendar.daysUntilDue != null">（{{ calendar.daysUntilDue }} 天）</span>
      </div>
    </div>

    <!-- 月历导航 -->
    <div class="flex items-center justify-between mb-6">
      <button type="button" class="member-btn member-btn--ghost" @click="prevMonth">
        <ArtSvgIcon icon="ri:arrow-left-s-line" class="text-xl text-clay-accent" />
      </button>
      <span class="font-bold text-lg text-clay-foreground">{{ viewYear }}年{{ viewMonth }}月</span>
      <button type="button" class="member-btn member-btn--ghost" @click="nextMonth">
        <ArtSvgIcon icon="ri:arrow-right-s-line" class="text-xl text-clay-accent" />
      </button>
    </div>

    <div v-if="loading" class="member-empty">
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
        <div class="mb-4">
          <div>
            <h3 class="text-sm font-black text-clay-muted uppercase tracking-widest">本年度任务列表</h3>
            <p class="text-xs text-clay-muted mt-1">可按申报状态筛选；未到申报期任务默认折叠。</p>
          </div>
          <div class="tax-task-filter" role="tablist" aria-label="申报任务筛选">
            <button
              v-for="option in taskFilterOptions"
              :key="option.value"
              type="button"
              class="tax-task-filter__btn"
              :class="{ 'is-active': taskFilter === option.value }"
              @click="taskFilter = option.value"
            >
              <span>{{ option.label }}</span>
              <span class="tax-task-filter__count">{{ option.count }}</span>
            </button>
          </div>
        </div>

        <div v-if="tasks.length === 0" class="text-sm text-clay-muted py-4">本年暂无申报任务</div>
        <div v-else-if="visibleRegularTasks.length === 0 && !showUpcomingSection" class="text-sm text-clay-muted py-4">当前筛选下暂无申报任务</div>
        <div v-else class="space-y-3">
          <div
            v-for="task in visibleRegularTasks"
            :key="task.id"
            class="p-4 rounded-lg border border-[#e8edf3] bg-[#f8fafc] flex flex-wrap items-center justify-between gap-3"
          >
            <div>
              <span class="font-bold text-clay-foreground">{{ formatTaskTitle(task) }}</span>
              <span class="ml-3 text-xs font-bold px-2 py-0.5 rounded-full" :class="statusClass(task.status)">
                {{ taskStatusLabel(task) }}
              </span>
              <p class="text-xs text-clay-muted mt-1">截止 {{ task.dueDate || '—' }} · 预估 ¥{{ formatMoney(task.calculatedAmount) }}</p>
            </div>
            <button
              type="button"
              class="text-sm font-bold text-clay-accent hover:underline"
              @click="openTaskDetail(task)"
            >
              {{ task.status === 'filed' ? '回执' : '详情' }}
            </button>
          </div>

          <div v-if="showUpcomingSection" class="rounded-lg border border-[#e8edf3] bg-white overflow-hidden">
            <button
              type="button"
              class="w-full px-4 py-3 flex items-center justify-between text-left hover:bg-[#f8fafc] transition-colors"
              @click="toggleUpcomingCollapsed"
            >
              <span class="font-bold text-clay-foreground">未到申报期任务</span>
              <span class="inline-flex items-center gap-2 text-xs font-bold text-clay-muted">
                {{ upcomingTasks.length }} 项
                <ArtSvgIcon
                  icon="ri:arrow-down-s-line"
                  class="text-base transition-transform"
                  :class="upcomingSectionCollapsed ? '' : 'rotate-180'"
                />
              </span>
            </button>
            <div v-if="!upcomingSectionCollapsed" class="border-t border-[#e8edf3] divide-y divide-[#e8edf3]">
              <div
                v-for="task in upcomingTasks"
                :key="task.id"
                class="p-4 flex flex-wrap items-center justify-between gap-3 bg-[#fbfdff]"
              >
                <div>
                  <span class="font-bold text-clay-foreground">{{ formatTaskTitle(task) }}</span>
                  <span class="ml-3 text-xs font-bold px-2 py-0.5 rounded-full" :class="statusClass(task.status)">
                    {{ taskStatusLabel(task) }}
                  </span>
                  <p class="text-xs text-clay-muted mt-1">截止 {{ task.dueDate || '—' }} · 预估 ¥{{ formatMoney(task.calculatedAmount) }}</p>
                </div>
                <button
                  type="button"
                  class="text-sm font-bold text-clay-accent hover:underline"
                  @click="openTaskDetail(task)"
                >
                  详情
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 自查清单 -->
      <div class="p-6 rounded-lg border border-[#e8edf3] bg-[#f8fafc]">
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
type TaskFilter = 'all' | 'pending' | 'filed' | 'overdue' | 'upcoming'
const taskFilter = ref<TaskFilter>('all')
const upcomingCollapsed = ref(true)

const todayStart = computed(() => {
  const date = new Date()
  date.setHours(0, 0, 0, 0)
  return date
})

const taskFilterOptions = computed<Array<{ value: TaskFilter; label: string; count: number }>>(() => [
  { value: 'all', label: '全部', count: tasks.value.length },
  { value: 'pending', label: '待申报', count: tasks.value.filter(task => getTaskGroup(task) === 'pending').length },
  { value: 'filed', label: '已完成', count: tasks.value.filter(task => getTaskGroup(task) === 'filed').length },
  { value: 'overdue', label: '已逾期', count: tasks.value.filter(task => getTaskGroup(task) === 'overdue').length },
  { value: 'upcoming', label: '未到申报期', count: upcomingTasks.value.length }
])

const upcomingTasks = computed(() => tasks.value.filter(task => getTaskGroup(task) === 'upcoming'))
const visibleRegularTasks = computed(() => {
  if (taskFilter.value === 'all') return tasks.value.filter(task => getTaskGroup(task) !== 'upcoming')
  if (taskFilter.value === 'upcoming') return []
  return tasks.value.filter(task => getTaskGroup(task) === taskFilter.value)
})
const showUpcomingSection = computed(() => (taskFilter.value === 'all' || taskFilter.value === 'upcoming') && upcomingTasks.value.length > 0)
const upcomingSectionCollapsed = computed(() => taskFilter.value === 'all' && upcomingCollapsed.value)

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

function getTaskGroup(task: TaxTask): TaskFilter {
  if (task.status === 'filed' || task.status === 'overdue') return task.status
  if (isUpcomingTask(task)) return 'upcoming'
  return 'pending'
}

function isUpcomingTask(task: TaxTask) {
  if (task.status !== 'pending') return false
  const dueDate = parseDate(task.dueDate)
  return !!dueDate && dueDate > todayStart.value
}

function parseDate(value?: string) {
  if (!value) return null
  const normalized = value.replace(/-/g, '/')
  const date = new Date(normalized)
  return Number.isNaN(date.getTime()) ? null : date
}

function taskStatusLabel(task: TaxTask) {
  if (getTaskGroup(task) === 'upcoming') return '未到申报期'
  return TAX_STATUS_LABELS[task.status]
}

function toggleUpcomingCollapsed() {
  if (taskFilter.value === 'upcoming') return
  upcomingCollapsed.value = !upcomingCollapsed.value
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
    upcomingCollapsed.value = true
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
.text-clay-accent { color: #2563eb; }

.tax-task-filter {
  display: flex;
  flex-wrap: nowrap;
  gap: 8px;
  margin-top: 14px;
  padding-bottom: 2px;
  overflow-x: auto;
  white-space: nowrap;
  scrollbar-width: thin;
}

.tax-task-filter__btn {
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  gap: 6px;
  padding: 7px 12px;
  color: #334155;
  font-size: 13px;
  font-weight: 700;
  background: #fff;
  border: 1px solid #d8dee9;
  border-radius: 8px;
  transition: all 0.2s ease;

  &:hover {
    background: #f8fafc;
    border-color: #94a3b8;
  }

  &.is-active {
    color: #fff;
    background: #2563eb;
    border-color: #2563eb;
  }
}

.tax-task-filter__count {
  min-width: 18px;
  padding: 1px 5px;
  color: inherit;
  font-size: 12px;
  line-height: 1.35;
  text-align: center;
  background: rgba(148, 163, 184, 0.18);
  border-radius: 999px;
}
</style>
