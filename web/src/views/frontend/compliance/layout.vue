<!-- +----------------------------------------------------------------------
  | XYGo Admin — 会员合规壳层布局
  +---------------------------------------------------------------------- -->
<template>
  <main class="compliance-layout">
    <div class="compliance-layout__grid">

      <!-- 左侧侧栏 -->
      <aside class="compliance-layout__sidebar">
        <!-- 用户信息卡片 -->
        <div class="compliance-card compliance-card--profile">
          <ElAvatar :size="72" :src="userInfo.avatar" class="compliance-card__avatar">
            {{ userInfo.nickname?.charAt(0) || 'U' }}
          </ElAvatar>
          <h2 class="compliance-card__name">{{ userInfo.nickname || userInfo.username }}</h2>
          <p class="compliance-card__role">税务合规服务</p>
          <RouterLink to="/user" class="compliance-card__back">
            <ArtSvgIcon icon="ri:arrow-left-line" />
            返回会员中心
          </RouterLink>
        </div>

        <!-- 合规菜单 -->
        <nav class="compliance-card compliance-card--nav">
          <template v-for="group in menuTree" :key="group.id">
            <div class="compliance-nav__group" :class="{ 'is-first': group === menuTree[0] }">
              <span class="compliance-nav__label">{{ group.name }}</span>
            </div>
            <ul class="compliance-nav__list">
              <li v-for="item in group.items" :key="item.id">
                <RouterLink
                  :to="item.path"
                  class="compliance-nav__item"
                  :class="{ 'is-active': isActive(item.path) }"
                >
                  <ArtSvgIcon :icon="item.icon" class="compliance-nav__icon" />
                  <span>{{ item.name }}</span>
                  <ArtSvgIcon v-if="isActive(item.path)" icon="ri:arrow-right-s-line" class="compliance-nav__arrow" />
                </RouterLink>
              </li>
            </ul>
          </template>
        </nav>
      </aside>

      <!-- 右侧内容 -->
      <div class="compliance-layout__content">
        <RouterView />
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { useMemberStore } from '@/store/modules/member'
import { buildComplianceMenuTree } from '@/config/complianceMenu'
import { getCompliancePlanState } from '@/api/frontend/compliance/member'
import type { CompliancePlanState } from '@/config/complianceMenu'

defineOptions({ name: 'ComplianceLayout' })

const route = useRoute()
const memberStore = useMemberStore()
const userInfo = computed(() => memberStore.getMemberInfo)

const planState = ref<CompliancePlanState>({ hasActiveOrder: false, opcStatus: 'none' })

onMounted(async () => {
  try {
    planState.value = await getCompliancePlanState()
  } catch { /* ignore */ }
})

const menuTree = computed(() => buildComplianceMenuTree(planState.value))

const isActive = (path: string) => route.path === path || route.path.startsWith(path + '/')
</script>

<style lang="scss" scoped>
.compliance-layout {
  padding: 48px 24px 32px;
  max-width: 1280px;
  margin: 0 auto;
  min-height: calc(100vh - 128px);
}

.compliance-layout__grid {
  display: grid;
  grid-template-columns: 260px 1fr;
  gap: 24px;
  align-items: start;
}

.compliance-card {
  padding: 24px;
  background: #fff;
  border: 1px solid #e8edf3;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);
}

.compliance-card--profile {
  text-align: center;
  margin-bottom: 16px;
}

.compliance-card__avatar {
  margin-bottom: 12px;
}

.compliance-card__name {
  margin: 0 0 4px;
  font-size: 18px;
  font-weight: 800;
  color: #1a1f36;
}

.compliance-card__role {
  margin: 0 0 16px;
  font-size: 12px;
  color: #94a3b8;
}

.compliance-card__back {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border-radius: 8px;
  background: #f8fafc;
  border: 1px solid #e8edf3;
  font-size: 12px;
  font-weight: 600;
  color: #64748b;
  text-decoration: none;
  transition: all 0.15s ease;

  &:hover {
    color: #2563eb;
    border-color: #bfdbfe;
    background: #eff6ff;
  }
}

.compliance-nav__group {
  padding: 12px 8px 8px;

  &.is-first {
    padding-top: 0;
  }
}

.compliance-nav__label {
  font-size: 11px;
  font-weight: 700;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.compliance-nav__list {
  margin: 0 0 8px;
  padding: 0;
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.compliance-nav__item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #334155;
  text-decoration: none;
  transition: all 0.15s ease;

  &:hover:not(.is-active) {
    background: #f8fafc;
    color: #2563eb;
  }

  &.is-active {
    color: #fff;
    background: #2563eb;
  }
}

.compliance-nav__icon {
  font-size: 18px;
  opacity: 0.85;
}

.compliance-nav__arrow {
  margin-left: auto;
  font-size: 16px;
}

@media (max-width: 1024px) {
  .compliance-layout__grid {
    grid-template-columns: 1fr;
  }
}
</style>
