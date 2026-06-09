<!-- +----------------------------------------------------------------------
  | XYGo Admin — 会员合规诊断历史
  +---------------------------------------------------------------------- -->
<template>
  <div class="diag-page">
    <section class="overview-panel">
      <div class="overview-panel__head">
        <div>
          <h2 class="overview-panel__section-title">诊断历史</h2>
          <p class="overview-panel__section-desc">查看您的合规诊断记录与推荐方案</p>
        </div>
        <RouterLink to="/diagnosis" class="overview-btn overview-btn--primary">
          <ArtSvgIcon icon="ri:add-line" />
          新建诊断
        </RouterLink>
      </div>

      <div v-if="loading && list.length === 0" class="overview-empty">
        <ArtSvgIcon icon="ri:loader-4-line" class="overview-empty__icon diag-loading" />
        <p>加载中...</p>
      </div>

      <div v-else-if="list.length === 0" class="overview-empty">
        <ArtSvgIcon icon="ri:file-search-line" class="overview-empty__icon" />
        <p>暂无诊断记录</p>
        <div class="overview-empty__actions">
          <RouterLink to="/diagnosis" class="overview-btn overview-btn--primary">立即开始免费诊断</RouterLink>
        </div>
      </div>

      <ul v-else class="diag-list">
        <li v-for="item in list" :key="item.id" class="diag-item">
          <div class="diag-item__head">
            <span class="diag-tag">{{ planLabel(item.recommendedPlan) }}</span>
            <span class="diag-item__time">{{ item.createdAt }}</span>
          </div>
          <div class="diag-item__grid">
            <div class="diag-field">
              <span class="diag-field__label">月收入区间</span>
              <span class="diag-field__value">{{ item.monthlyIncomeRange }}</span>
            </div>
            <div class="diag-field">
              <span class="diag-field__label">年成本估算</span>
              <span class="diag-field__value">{{ formatMoney(item.annualCostEstimate) }} 元</span>
            </div>
            <div class="diag-field">
              <span class="diag-field__label">现有主体</span>
              <span class="diag-field__value">{{ entityLabel(item.existingEntity) }}</span>
            </div>
            <div class="diag-field">
              <span class="diag-field__label">收入渠道</span>
              <span class="diag-field__value">{{ item.platforms?.join('、') || '-' }}</span>
            </div>
          </div>
          <div v-if="item.taxComparison?.items?.length" class="diag-item__tax">
            <span class="diag-field__label">推荐方案预估年税负</span>
            <div class="diag-tax-tags">
              <span
                v-for="tax in item.taxComparison.items.filter(t => t.recommended)"
                :key="tax.plan"
                class="diag-tax-tag"
              >
                {{ tax.label }}：{{ formatMoney(tax.annualTax) }} 元
              </span>
            </div>
          </div>
          <div class="diag-item__footer">
            <RouterLink :to="`/diagnosis/result?id=${item.id}`" class="diag-link">
              查看详情
              <ArtSvgIcon icon="ri:arrow-right-s-line" />
            </RouterLink>
          </div>
        </li>
      </ul>

      <div v-if="total > pageSize && page * pageSize < total" class="diag-load-more">
        <button type="button" class="overview-btn overview-btn--ghost" :disabled="loading" @click="loadMore">
          <ArtSvgIcon v-if="loading" icon="ri:loader-4-line" class="diag-loading" />
          加载更多
        </button>
      </div>
    </section>
  </div>
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
    opc: '小微公司',
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
.diag-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.overview-panel {
  padding: 24px;
  background: #fff;
  border: 1px solid #e8edf3;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);
}

.overview-panel__head {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 4px;
}

.overview-panel__section-title {
  margin: 0 0 6px;
  font-size: 18px;
  font-weight: 700;
  color: #1a1f36;
}

.overview-panel__section-desc {
  margin: 0 0 20px;
  font-size: 13px;
  color: #94a3b8;
}

.overview-empty {
  padding: 32px 16px;
  text-align: center;

  p {
    margin: 0 0 20px;
    font-size: 14px;
    color: #64748b;
    font-weight: 600;
  }
}

.overview-empty__icon {
  font-size: 40px;
  color: #2563eb;
  margin-bottom: 16px;
}

.overview-empty__actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 12px;
}

.overview-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 20px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  text-decoration: none;
  border: none;
  cursor: pointer;
  transition: all 0.15s ease;

  &--primary {
    color: #fff;
    background: #2563eb;

    &:hover {
      background: #1d4ed8;
    }
  }

  &--ghost {
    color: #334155;
    background: #fff;
    border: 1px solid #d8dee9;

    &:hover:not(:disabled) {
      background: #f8fafc;
    }

    &:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }
  }
}

.diag-list {
  margin: 0;
  padding: 0;
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.diag-item {
  padding: 16px;
  border: 1px solid #e8edf3;
  border-radius: 10px;
  background: #f8fafc;
}

.diag-item__head {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
  margin-bottom: 14px;
}

.diag-tag {
  padding: 2px 8px;
  border-radius: 6px;
  background: #eff6ff;
  color: #2563eb;
  font-size: 11px;
  font-weight: 700;
}

.diag-item__time {
  font-size: 12px;
  color: #94a3b8;
}

.diag-item__grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}

.diag-field__label {
  display: block;
  margin-bottom: 4px;
  font-size: 11px;
  font-weight: 600;
  color: #94a3b8;
}

.diag-field__value {
  font-size: 14px;
  font-weight: 600;
  color: #1a1f36;
}

.diag-item__tax {
  margin-top: 14px;
  padding-top: 14px;
  border-top: 1px solid #e8edf3;
}

.diag-tax-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 8px;
}

.diag-tax-tag {
  padding: 4px 10px;
  border-radius: 6px;
  background: #f0fdf4;
  color: #15803d;
  font-size: 12px;
  font-weight: 600;
}

.diag-item__footer {
  margin-top: 14px;
  text-align: right;
}

.diag-link {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  font-weight: 700;
  color: #2563eb;
  text-decoration: none;

  &:hover {
    text-decoration: underline;
  }
}

.diag-load-more {
  margin-top: 16px;
  text-align: center;
}

.diag-loading {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@media (max-width: 768px) {
  .diag-item__grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
