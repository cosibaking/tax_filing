// +----------------------------------------------------------------------
// | XYGo Admin [ Vue3 + GoFrame 企业级中后台管理系统 ]
// +----------------------------------------------------------------------
// | Copyright (c) 2026 大连星韵网络科技有限公司 All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( https://opensource.org/licenses/MIT )
// +----------------------------------------------------------------------
// | Author: 喜羊羊 <751300685@qq.com>
// +----------------------------------------------------------------------

/**
 * useCommon - 通用功能集合
 *
 * 提供常用的页面操作功能，包括页面刷新、滚动控制、路径获取等。
 * 这些功能在多个页面和组件中都会用到，统一封装便于复用。
 *
 * ## 主要功能
 *
 * 1. 首页路径 - 获取系统配置的首页路径
 * 2. 页面刷新 - 刷新当前页面内容
 * 3. 滚动控制 - 提供多种滚动到顶部和指定位置的方法
 * 4. 平滑滚动 - 支持平滑滚动动画效果
 *
 * @module useCommon
 * @author Art Design Pro Team
 */

import { computed } from 'vue'
import { useMenuStore } from '@/store/modules/menu'
import { useSettingStore } from '@/store/modules/setting'

export function useCommon() {
  const menuStore = useMenuStore()
  const settingStore = useSettingStore()

  /**
   * 首页路径
   * 从菜单 store 中获取配置的首页路径
   */
  const homePath = computed(() => menuStore.getHomePath())

  /**
   * 刷新当前页面
   * 通过切换 setting store 中的 refresh 状态触发页面重新渲染
   */
  const refresh = () => {
    settingStore.reload()
  }

  const getScrollContainer = () => document.getElementById('app-main')

  /**
   * 滚动到页面顶部
   * 后台滚动 #app-main，前台回退到 window
   */
  const scrollToTop = () => {
    const scrollContainer = getScrollContainer()
    if (scrollContainer) {
      scrollContainer.scrollTop = 0
      return
    }
    window.scrollTo(0, 0)
  }

  /**
   * 平滑滚动到页面顶部
   * 使用 smooth 行为实现平滑滚动效果
   */
  const smoothScrollToTop = () => {
    const scrollContainer = getScrollContainer()
    if (scrollContainer) {
      scrollContainer.scrollTo({
        top: 0,
        behavior: 'smooth'
      })
      return
    }
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  /**
   * 滚动到指定位置
   * @param top 目标滚动位置（像素）
   * @param smooth 是否使用平滑滚动
   */
  const scrollTo = (top: number, smooth: boolean = false) => {
    const scrollContainer = getScrollContainer()
    if (scrollContainer) {
      scrollContainer.scrollTo({
        top,
        behavior: smooth ? 'smooth' : 'auto'
      })
      return
    }
    window.scrollTo({
      top,
      behavior: smooth ? 'smooth' : 'auto'
    })
  }

  return {
    homePath,
    refresh,
    scrollTo,
    scrollToTop,
    smoothScrollToTop
  }
}
