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
 * 点选验证码公共 API（前后台通用）
 * 接口无需鉴权，路径: /captcha/click, /captcha/checkClick
 */
import api from '@/utils/http'

/** 点选验证码响应 */
export interface ClickCaptchaResult {
  id: string
  text: string[]
  base64: string
  width: number
  height: number
}

/**
 * 获取点选验证码
 */
export function getClickCaptcha() {
  return api.get<ClickCaptchaResult>({
    url: '/captcha/click',
    showErrorMessage: false,
  })
}

/**
 * 校验点选验证码
 */
export function checkClickCaptcha(id: string, info: string) {
  return api.post<void>({
    url: '/captcha/checkClick',
    data: { id, info }
  })
}
