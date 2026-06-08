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
 * 前台门户 API 模块
 *
 * 目录结构：
 * frontend/
 * ├── member/       # 会员相关 API（登录、注册、个人信息、签到、积分、余额、通知）
 * └── index.ts      # 统一导出
 */
export * from './member/auth'
export * from './member/user'
export * from './member/checkin'
export * from './member/score-log'
export * from './member/money-log'
export * from './member/notice'
export * from './doc'
export * from './compliance/diagnosis'
export * from './compliance/order'
export * from './compliance/opc'
export * from './compliance/member'
