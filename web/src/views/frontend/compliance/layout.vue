<!-- +----------------------------------------------------------------------
  | XYGo Admin — 会员合规壳层布局
  +---------------------------------------------------------------------- -->
<template>
  <main class="pt-12 pb-8 px-6 max-w-7xl mx-auto min-h-[calc(100vh-128px)] flex flex-col">
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
          <p class="text-xs text-clay-muted font-medium">税务合规服务</p>
          <RouterLink
            to="/user"
            class="inline-flex items-center gap-2 mt-4 px-4 py-2 rounded-xl bg-[#f0f3f8] shadow-clay-pressed text-xs font-bold text-clay-muted hover:text-clay-accent transition-colors"
          >
            <ArtSvgIcon icon="ri:arrow-left-line" class="text-sm" />
            返回会员中心
          </RouterLink>
        </div>

        <!-- 合规菜单 -->
        <nav class="bg-white/70 backdrop-blur-xl rounded-[40px] shadow-clay-card border border-[#d1d9e6]/40 p-4 overflow-hidden flex-1">
          <template v-for="group in menuTree" :key="group.id">
            <div class="px-4 py-3 mb-2" :class="{ 'mt-4': group !== menuTree[0] }">
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
.text-clay-foreground { color: #32325d; }
.text-clay-muted { color: #8898aa; }
.text-clay-accent { color: #5a8dee; }
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
