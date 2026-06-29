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
import defaultLogo from '@imgs/common/jinshui-logo.png'

const BRAND_SITE_NAME = '金税管家'
const LEGACY_SITE_NAMES = new Set(['XYGo Admin', 'XYgoAdmin', 'XYGo'])
const LEGACY_LOGO_MARKERS = ['995f9919-23d3-4ce2-8564-2460e4b1261d', 'logo.webp', '/attachment/upload/20260610/jinshui-logo.png']

function normalizeSiteInfo(info: Partial<SiteInfo>): Partial<SiteInfo> {
  const normalized = { ...info }
  if (!normalized.siteName || LEGACY_SITE_NAMES.has(normalized.siteName)) {
    normalized.siteName = BRAND_SITE_NAME
  }
  if (!normalized.logo || LEGACY_LOGO_MARKERS.some((marker) => normalized.logo?.includes(marker))) {
    normalized.logo = ''
  }
  return normalized
}

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

    // 是否已加载
    const loaded = ref(false)

    // 时区模式：从 .env 配置读取，client=浏览器本地时区 | server=后端配置时区
    const envMode = (import.meta.env.VITE_TIMEZONE_MODE || 'client') as 'client' | 'server'
    const timeZoneMode = ref<'client' | 'server'>(envMode)

    // 前台配置
    const frontendConfig = ref({
      userCenterEnabled: true,
      portalEnabled: true,
      registerEnabled: true
    })

    /**
     * 加载站点信息
     */
    const loadSiteInfo = async (): Promise<void> => {
      try {
        const res = await fetchSiteIndex()
        siteInfo.value = normalizeSiteInfo(res)
        loaded.value = true

        // 同步后端开关到前台配置
        frontendConfig.value.userCenterEnabled = res.openMemberCenter !== false
      } catch (error) {
        console.error('[SiteStore] 加载站点信息失败:', error)
        // 加载失败时使用默认值
      }
    }

    /**
     * 获取站点名称
     */
    const getSiteName = (): string => {
      return siteInfo.value.siteName || BRAND_SITE_NAME
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
      return siteInfo.value.logo || defaultLogo
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
     * 是否为演示模式
     */
    const isDemoMode = (): boolean => {
      return !!(siteInfo.value as any).demoMode
    }

    /**
     * 是否启用门户首页
     */
    const isPortalEnabled = (): boolean => {
      return frontendConfig.value.portalEnabled
    }

    /**
     * 获取时区
     */
    const getTimezone = (): string => {
      return siteInfo.value.timezone || 'Asia/Shanghai'
    }

    /**
     * 获取时区模式
     */
    const getTimeZoneMode = (): 'client' | 'server' => {
      return timeZoneMode.value
    }

    /**
     * 设置时区模式
     */
    const setTimeZoneMode = (mode: 'client' | 'server') => {
      timeZoneMode.value = mode
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
      timeZoneMode,
      loadSiteInfo,
      getSiteName,
      getSiteSubtitle,
      getLogo,
      getThemeColor,
      getTimezone,
      getTimeZoneMode,
      setTimeZoneMode,
      isDemoMode,
      isUserCenterEnabled,
      isPortalEnabled,
      setFrontendConfig
    }
  }
)
