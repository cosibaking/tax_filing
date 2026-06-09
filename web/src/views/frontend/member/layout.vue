<!-- +----------------------------------------------------------------------
  | XYGo Admin — 会员中心统一壳层（M8 §5.2 侧栏 + §5.3 引导条）
  +---------------------------------------------------------------------- -->
<template>
  <main class="pt-12 pb-8 px-6 max-w-7xl mx-auto min-h-[calc(100vh-128px)] flex flex-col">
    <!-- 全局引导条 -->
    <ComplianceGuideBanner v-if="guideBanner" :banner="guideBanner" class="mb-6" />

    <div class="grid lg:grid-cols-12 gap-8 items-stretch flex-1">
      <!-- 左侧侧栏 -->
      <aside class="lg:col-span-3 flex flex-col gap-6">
        <!-- 用户信息卡片 -->
        <div class="bg-white/70 backdrop-blur-xl rounded-[40px] shadow-clay-card border border-[#d1d9e6]/40 p-8 text-center relative overflow-hidden group">
          <div class="absolute -top-10 -right-10 w-32 h-32 rounded-full bg-blue-500/5 blur-2xl group-hover:scale-150 transition-transform"></div>
          <div class="relative inline-block mb-4">
            <div class="w-24 h-24 rounded-[32px] bg-white shadow-clay-btn p-1 animate-breathe">
              <ElAvatar :size="88" :src="userInfo.avatar" class="!rounded-[28px] !w-full !h-full">
                {{ userInfo.nickname?.charAt(0) || 'U' }}
              </ElAvatar>
            </div>
          </div>
          <h2 class="font-heading font-black text-xl text-clay-foreground mb-1">{{ userInfo.nickname || userInfo.username }}</h2>
          <div class="flex justify-center gap-3 mt-4">
            <div class="px-3 py-1.5 rounded-2xl bg-[#f0f3f8] shadow-clay-pressed text-xs font-bold text-clay-muted">
              积分 <span class="text-clay-accent">{{ userInfo.score ?? 0 }}</span>
            </div>
            <div class="px-3 py-1.5 rounded-2xl bg-[#f0f3f8] shadow-clay-pressed text-xs font-bold text-clay-muted">
              余额 <span class="text-clay-success">{{ formatMoney(userInfo.money) }}</span>
            </div>
          </div>
        </div>

        <!-- 合规服务菜单 -->
        <nav v-if="complianceMenuTree.length > 0" class="bg-white/70 backdrop-blur-xl rounded-[40px] shadow-clay-card border border-[#d1d9e6]/40 p-4 overflow-hidden">
          <template v-for="group in complianceMenuTree" :key="group.id">
            <div class="px-4 py-3 mb-2" :class="{ 'mt-4': group !== complianceMenuTree[0] }">
              <span class="text-xs font-black text-clay-muted uppercase tracking-widest">{{ group.name }}</span>
            </div>
            <ul class="space-y-2">
              <li v-for="item in group.items" :key="item.id">
                <RouterLink
                  :to="item.path"
                  class="flex items-center gap-4 px-6 py-4 rounded-[24px] transition-all duration-300 group"
                  :class="isActive(item.path)
                    ? 'bg-gradient-to-br from-blue-400 to-blue-600 text-white shadow-clay-btn'
                    : 'text-clay-foreground hover:bg-white hover:shadow-clay-card'"
                >
                  <ArtSvgIcon
                    :icon="item.icon"
                    class="text-xl"
                    :class="isActive(item.path) ? 'text-white' : 'text-clay-accent opacity-70 group-hover:opacity-100'"
                  />
                  <span class="font-bold text-sm">{{ item.name }}</span>
                  <ArtSvgIcon v-if="isActive(item.path)" icon="ri:arrow-right-s-line" class="text-base ml-auto" />
                </RouterLink>
              </li>
            </ul>
          </template>
        </nav>

        <!-- 账户菜单（后端动态） -->
        <nav class="bg-white/70 backdrop-blur-xl rounded-[40px] shadow-clay-card border border-[#d1d9e6]/40 p-4 overflow-hidden flex-1">
          <template v-for="group in accountMenuTree" :key="group.id">
            <div class="px-4 py-3 mb-2" :class="{ 'mt-4': group !== accountMenuTree[0] }">
              <span class="text-xs font-black text-clay-muted uppercase tracking-widest">{{ group.name }}</span>
            </div>
            <ul class="space-y-2">
              <li v-for="item in group.children" :key="item.id">
                <RouterLink
                  :to="item.path"
                  class="flex items-center gap-4 px-6 py-4 rounded-[24px] transition-all duration-300 group"
                  :class="isActive(item.path)
                    ? 'bg-gradient-to-br from-blue-400 to-blue-600 text-white shadow-clay-btn'
                    : 'text-clay-foreground hover:bg-white hover:shadow-clay-card'"
                >
                  <ArtSvgIcon
                    :icon="item.icon"
                    class="text-xl"
                    :class="isActive(item.path) ? 'text-white' : 'text-clay-accent opacity-70 group-hover:opacity-100'"
                  />
                  <span class="font-bold text-sm">{{ item.name }}</span>
                  <ArtSvgIcon v-if="isActive(item.path)" icon="ri:arrow-right-s-line" class="text-base ml-auto" />
                </RouterLink>
              </li>
            </ul>
          </template>
        </nav>
      </aside>

      <!-- 右侧内容 -->
      <div class="lg:col-span-9">
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
import { memberMenuHref } from '@/utils/member-nav'
import ComplianceGuideBanner, { type GuideBanner } from '@/components/member/ComplianceGuideBanner.vue'

defineOptions({ name: 'MemberLayout' })

const route = useRoute()
const memberStore = useMemberStore()
const memberMenuStore = useMemberMenuStore()
const userInfo = computed(() => memberStore.getMemberInfo)

const planState = ref<CompliancePlanState>({ hasActiveOrder: false, opcStatus: 'none' })

async function loadPlanState() {
  if (!memberStore.getIsLogin) return
  try { planState.value = await getCompliancePlanState() } catch { /* ignore */ }
}

onMounted(async () => {
  if (!memberStore.getIsLogin) return
  try { await memberMenuStore.fetchMenus() } catch { /* ignore */ }
  await loadPlanState()
})

watch(() => route.path, () => {
  loadPlanState()
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
    name: '概览',
    icon: 'ri:home-4-line',
    path: '/user/overview',
    requiresOpcActive: false,
    visible: () => true
  }
  if (tree.length === 0) {
    return [{ id: 'service', name: '合规服务', items: [overviewItem] }]
  }
  const serviceGroup = tree.find((g) => g.id === 'service')
  if (serviceGroup) {
    return [{ ...serviceGroup, items: [overviewItem, ...serviceGroup.items] }, ...tree.filter((g) => g.id !== 'service')]
  }
  return [{ id: 'service', name: '合规服务', items: [overviewItem] }, ...tree]
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

  if (dirs.length === 0) {
    return [{
      id: 'default',
      name: '我的账户',
      children: items.map(mapItem),
    }]
  }

  return dirs.map((dir) => ({
    id: dir.name || String(dir.id),
    name: dir.title,
    children: items.filter((m) => m.pid === dir.id).map(mapItem),
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
    return {
      type: 'warning',
      message: '完成方案签约，开启 OPC 合规服务',
      actionLabel: '去签约',
      actionPath: '/user/compliance/plan',
    }
  }
  if (planState.value.hasActiveOrder && planState.value.opcStatus !== 'active') {
    return {
      type: 'info',
      message: 'OPC 设立进行中，台账功能将在激活后开放',
      actionLabel: '查看进度',
      actionPath: '/user/compliance/opc',
    }
  }
  return null
})

const isActive = (path: string) => route.path === path || route.path.startsWith(path + '/')

const formatMoney = (v: unknown) => {
  const n = Number(v) || 0
  return n.toFixed(2)
}
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

@keyframes breathe { 0%, 100% { transform: scale(1); } 50% { transform: scale(1.05); } }
.animate-breathe { animation: breathe 6s ease-in-out infinite; }
</style>
