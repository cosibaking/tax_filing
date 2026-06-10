import type { AppRouteRecord } from '@/types/router'
import { getComplianceDashboard } from '@/api/backend/compliance'
import { useMenuStore } from '@/store/modules/menu'

function patchMenuBadgeByName(menus: AppRouteRecord[], name: string, showBadge: boolean): void {
  for (const menu of menus) {
    if (menu.name === name) {
      if (!menu.meta) menu.meta = { title: '' }
      menu.meta.showBadge = showBadge
      return
    }
    if (menu.children?.length) {
      patchMenuBadgeByName(menu.children, name, showBadge)
    }
  }
}

/** 根据合规待办动态刷新后台侧栏菜单小红点 */
export function useComplianceMenuBadges() {
  const menuStore = useMenuStore()

  async function refreshComplianceMenuBadges() {
    try {
      const data = await getComplianceDashboard()
      patchMenuBadgeByName(
        menuStore.menuList,
        'ComplianceSocialConsults',
        (data.pendingSocialConsults ?? 0) > 0
      )
    } catch {
      /* ignore */
    }
  }

  return { refreshComplianceMenuBadges }
}
