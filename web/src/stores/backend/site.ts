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
 * 站点信息状态管理模块
 *
 * 提供站点配置信息的统一管理
 *
 * @module store/modules/site
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { fetchSiteIndex, type SiteInfo } from '@/api/site'

/**
 * 站点信息状态管理
 */
export const useSiteStore = defineStore(
  'siteStore',
  () => {
    // 站点信息
    const siteInfo = ref<Partial<SiteInfo>>({
      siteName: '金税管家',
      siteSubtitle: '个人税务合规与申报一站式服务',
      logo: '',
      themeColor: '#1890ff',
      description: '',
      icp: '',
      timezone: 'Asia/Shanghai',
      closed: '0'
    })

    // 前台配置
    const frontendConfig = ref({
      // 是否启用用户中心
      userCenterEnabled: true,
      // 是否启用门户首页
      portalEnabled: true,
      // 是否允许注册
      registerEnabled: true
    })

    // 是否已加载
    const loaded = ref(false)

    /**
     * 加载站点信息
     */
    const loadSiteInfo = async (): Promise<void> => {
      try {
        const res = await fetchSiteIndex()
        siteInfo.value = res
        loaded.value = true
      } catch (error) {
        console.error('[SiteStore] 加载站点信息失败:', error)
        // 加载失败时使用默认值
      }
    }

    /**
     * 获取站点名称
     */
    const getSiteName = (): string => {
      return siteInfo.value.siteName || '金税管家'
    }

    /**
     * 获取站点副标题
     */
    const getSiteSubtitle = (): string => {
      return siteInfo.value.siteSubtitle || ''
    }

    /**
     * 获取Logo
     */
    const getLogo = (): string => {
      return siteInfo.value.logo || ''
    }

    /**
     * 获取主题色
     */
    const getThemeColor = (): string => {
      return siteInfo.value.themeColor || '#1890ff'
    }

    /**
     * 是否启用用户中心
     */
    const isUserCenterEnabled = (): boolean => {
      return frontendConfig.value.userCenterEnabled
    }

    /**
     * 是否启用门户首页
     */
    const isPortalEnabled = (): boolean => {
      return frontendConfig.value.portalEnabled
    }

    /**
     * 是否允许注册
     */
    const isRegisterEnabled = (): boolean => {
      return frontendConfig.value.registerEnabled
    }

    /**
     * 设置前台配置
     */
    const setFrontendConfig = (config: Partial<typeof frontendConfig.value>) => {
      frontendConfig.value = { ...frontendConfig.value, ...config }
    }

    return {
      siteInfo,
      loaded,
      frontendConfig,
      loadSiteInfo,
      getSiteName,
      getSiteSubtitle,
      getLogo,
      getThemeColor,
      isUserCenterEnabled,
      isPortalEnabled,
      isRegisterEnabled,
      setFrontendConfig
    }
  }
)
