<!-- 会员中心壳层 -->
<template>
  <main class="member-layout">
    <ComplianceGuideBanner v-if="guideBanner" :banner="guideBanner" class="member-layout__banner" />

    <div class="member-layout__grid">
      <aside class="member-layout__aside">
        <!-- 用户卡片 -->
        <div class="member-card member-card--profile">
          <ElAvatar :size="56" :src="userInfo.avatar" class="member-card__avatar">
            {{ userInfo.nickname?.charAt(0) || 'U' }}
          </ElAvatar>
          <h2 class="member-card__name">{{ userInfo.nickname || userInfo.username }}</h2>
          <p class="member-card__meta">金税管家会员</p>
        </div>

        <!-- 合规服务 -->
        <nav v-if="complianceMenuTree.length > 0" class="member-card">
          <template v-for="group in complianceMenuTree" :key="group.id">
            <p class="member-nav__group">{{ group.name }}</p>
            <ul class="member-nav">
              <li v-for="item in group.items" :key="item.id">
                <RouterLink
                  :to="item.path"
                  class="member-nav__link"
                  :class="{ 'is-active': isActive(item.path) }"
                >
                  <ArtSvgIcon :icon="item.icon" class="member-nav__icon" />
                  <span>{{ item.name }}</span>
                  <span v-if="item.showBadge" class="member-nav__badge" aria-label="有新回复" />
                </RouterLink>
              </li>
            </ul>
          </template>
        </nav>

        <!-- 账户设置 -->
        <nav v-if="accountMenuTree.length > 0" class="member-card">
          <template v-for="group in accountMenuTree" :key="group.id">
            <p class="member-nav__group">{{ group.name }}</p>
            <ul class="member-nav">
              <li v-for="item in group.children" :key="item.id">
                <RouterLink
                  :to="item.path"
                  class="member-nav__link"
                  :class="{ 'is-active': isActive(item.path) }"
                >
                  <ArtSvgIcon :icon="item.icon" class="member-nav__icon" />
                  <span>{{ item.name }}</span>
                </RouterLink>
              </li>
            </ul>
          </template>
        </nav>
      </aside>

      <div class="member-layout__main">
        <RouterView />
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { useMemberStore } from '@/store/modules/member'
import { useMemberMenuStore } from '@/store/modules/memberMenu'
import { buildComplianceMenuTree } from '@/config/complianceMenu'
import type { CompliancePlanState } from '@/config/complianceMenu'
import { getCompliancePlanState } from '@/api/frontend/compliance/member'
import { getSocialConsults } from '@/api/frontend/compliance/social'
import { memberMenuHref } from '@/utils/member-nav'
import ComplianceGuideBanner, { type GuideBanner } from '@/components/member/ComplianceGuideBanner.vue'

defineOptions({ name: 'MemberLayout' })

const route = useRoute()
const memberStore = useMemberStore()
const memberMenuStore = useMemberMenuStore()
const userInfo = computed(() => memberStore.getMemberInfo)

/** 与税务合规业务无关的账户菜单 */
const hiddenAccountPaths = new Set(['/user/checkin', '/user/points', '/user/balance'])

const planState = ref<CompliancePlanState>({ hasActiveOrder: false, opcStatus: 'none' })
const socialConsultHasReply = ref(false)

async function loadPlanState() {
  if (!memberStore.getIsLogin) return
  try { planState.value = await getCompliancePlanState() } catch { /* ignore */ }
}

async function loadSocialConsultBadge() {
  if (!memberStore.getIsLogin || planState.value.opcStatus !== 'active') {
    socialConsultHasReply.value = false
    return
  }
  try {
    const res = await getSocialConsults({ page: 1, pageSize: 50 })
    socialConsultHasReply.value = (res.list || []).some((item) => item.status === 'replied')
  } catch {
    socialConsultHasReply.value = false
  }
}

async function refreshMemberShell() {
  await loadPlanState()
  await loadSocialConsultBadge()
}

onMounted(async () => {
  if (!memberStore.getIsLogin) return
  try { await memberMenuStore.fetchMenus() } catch { /* ignore */ }
  await refreshMemberShell()
})

watch(() => route.path, () => {
  refreshMemberShell()
})

interface AccountMenuGroup {
  id: string
  name: string
  children: { id: string; name: string; icon: string; path: string }[]
}

const complianceMenuTree = computed(() => {
  const tree = buildComplianceMenuTree(planState.value)
  const overviewItem = {
    id: 'overview',
    name: '服务概览',
    icon: 'ri:home-4-line',
    path: '/user/overview',
    showBadge: false,
  }
  const withBadges = (items: typeof tree[0]['items']) =>
    items.map((item) => ({
      ...item,
      showBadge: item.id === 'social-consult' && socialConsultHasReply.value,
    }))

  if (tree.length === 0) {
    return [{ id: 'service', name: '合规服务', items: [overviewItem] }]
  }
  const serviceGroup = tree.find((g) => g.id === 'service')
  if (serviceGroup) {
    return [
      { ...serviceGroup, items: [overviewItem, ...withBadges(serviceGroup.items)] },
      ...tree.filter((g) => g.id !== 'service').map((g) => ({ ...g, items: withBadges(g.items) })),
    ]
  }
  return [
    { id: 'service', name: '合规服务', items: [overviewItem] },
    ...tree.map((g) => ({ ...g, items: withBadges(g.items) })),
  ]
})

const accountMenuTree = computed<AccountMenuGroup[]>(() => {
  const raw = memberMenuStore.getCenterMenus
  if (raw.length === 0) return []

  const dirs = raw.filter((m) => m.type === 'menu_dir')
  const items = raw.filter((m) => m.type === 'menu' && m.name !== 'overview')

  const mapItem = (m: typeof items[0]) => {
    const { url } = memberMenuHref(m)
    return {
      id: m.name || String(m.id),
      name: m.title,
      icon: m.icon || 'ri:menu-line',
      path: url,
    }
  }

  const filterItems = (list: ReturnType<typeof mapItem>[]) =>
    list.filter((item) => !hiddenAccountPaths.has(item.path))

  if (dirs.length === 0) {
    const children = filterItems(items.map(mapItem))
    return children.length ? [{ id: 'default', name: '账户设置', children }] : []
  }

  return dirs.map((dir) => ({
    id: dir.name || String(dir.id),
    name: dir.title === '我的账户' ? '账户设置' : dir.title,
    children: filterItems(items.filter((m) => m.pid === dir.id).map(mapItem)),
  })).filter((g) => g.children.length > 0)
})

const guideBanner = computed<GuideBanner | null>(() => {
  if (!planState.value.hasActiveOrder) {
    if (planState.value.hasPendingOrder) {
      return {
        type: 'warning',
        message: '您有未完成的签约流程，请继续完成风险告知与电子签约',
        actionLabel: '继续签约',
        actionPath: '/user/compliance/plan',
      }
    }
    if (planState.value.hasDiagnosis) {
      return {
        type: 'warning',
        message: '您已完成免费诊断，请选择套餐并完成签约',
        actionLabel: '去签约',
        actionPath: '/user/compliance/plan',
      }
    }
    return {
      type: 'warning',
      message: '完成方案签约，开启合规服务',
      actionLabel: '去签约',
      actionPath: '/user/compliance/plan',
    }
  }
  if (planState.value.hasActiveOrder && planState.value.opcStatus !== 'active') {
    return {
      type: 'info',
      message: '经营主体设立进行中，台账功能将在激活后开放',
      actionLabel: '查看进度',
      actionPath: '/user/compliance/opc',
    }
  }
  return null
})

const isActive = (path: string) => route.path === path || route.path.startsWith(path + '/')
</script>

<style lang="scss" scoped>
.member-layout {
  max-width: 1120px;
  margin: 0 auto;
  padding: 24px 24px 64px;
  min-height: calc(100vh - 128px);
}

.member-layout__banner {
  margin-bottom: 20px;
}

.member-layout__grid {
  display: grid;
  grid-template-columns: 240px 1fr;
  gap: 24px;
  align-items: start;
}

.member-layout__aside {
  display: flex;
  flex-direction: column;
  gap: 16px;
  position: sticky;
  top: 96px;
}

.member-card {
  padding: 16px;
  background: #fff;
  border: 1px solid #e8edf3;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.04);
}

.member-card--profile {
  text-align: center;
  padding: 20px 16px;
}

.member-card__avatar {
  margin: 0 auto 12px;
}

.member-card__name {
  margin: 0 0 4px;
  font-size: 16px;
  font-weight: 700;
  color: #1a1f36;
}

.member-card__meta {
  margin: 0;
  font-size: 12px;
  color: #94a3b8;
}

.member-nav__group {
  margin: 0 0 8px;
  padding: 0 8px;
  font-size: 11px;
  font-weight: 700;
  color: #94a3b8;
  letter-spacing: 0.06em;
  text-transform: uppercase;
}

.member-nav {
  margin: 0 0 12px;
  padding: 0;
  list-style: none;

  &:last-child {
    margin-bottom: 0;
  }
}

.member-nav__link {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #475569;
  text-decoration: none;
  transition: background 0.15s ease, color 0.15s ease;

  &:hover {
    background: #f8fafc;
    color: #2563eb;
  }

  &.is-active {
    background: #eff6ff;
    color: #2563eb;
  }
}

.member-nav__icon {
  font-size: 18px;
  flex-shrink: 0;
}

.member-nav__badge {
  width: 6px;
  height: 6px;
  margin-left: auto;
  background: #ef4444;
  border-radius: 50%;
  flex-shrink: 0;
  animation: member-nav-badge-breathe 1.5s ease-in-out infinite;
}

@keyframes member-nav-badge-breathe {
  0%,
  100% {
    opacity: 1;
    transform: scale(1);
  }

  50% {
    opacity: 0.65;
    transform: scale(0.92);
  }
}

@media (max-width: 960px) {
  .member-layout__grid {
    grid-template-columns: 1fr;
  }

  .member-layout__aside {
    position: static;
  }
}
</style>
