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

 * Frontend 前台路由配置（对齐 BuildAdmin 路径隔离模式）

 *

 * 所有前台页面在 FrontendLayout (/) 下静态注册。

 * 菜单 API 仅驱动导航栏 UI，不做动态路由注册。

 * 后台页面全部在 /admin 下，前后台路由互不干扰。

 */



import { AppRouteRecordRaw } from '@/utils/router'

import { complianceRouteGuard } from '@/router/guards/complianceGuard'



const memberCenter = () => import('@/views/frontend/member/center.vue')



export const frontendRoutes: AppRouteRecordRaw[] = [

  {

    path: '/',

    name: 'FrontendLayout',

    component: () => import('@/views/frontend/layouts/FrontendLayout.vue'),

    meta: { title: '首页' },

    children: [

      // 首页（未登录展示免费诊断，已登录展示营销主页）

      {

        path: '',

        name: 'FrontendHome',

        component: () => import('@/views/frontend/index/entry.vue'),

        meta: { title: '首页' }

      },

      // 合规诊断

      {

        path: 'diagnosis',

        name: 'ComplianceDiagnosis',

        component: () => import('@/views/frontend/diagnosis/index.vue'),

        meta: { title: '合规诊断' }

      },

      {

        path: 'diagnosis/result',

        name: 'ComplianceDiagnosisResult',

        component: () => import('@/views/frontend/diagnosis/result.vue'),

        meta: { title: '诊断结果' }

      },

      // 服务套餐

      {

        path: 'pricing',

        name: 'CompliancePricing',

        component: () => import('@/views/frontend/pricing/index.vue'),

        meta: { title: '服务价格' }

      },

      // 应用案例

      {

        path: 'cases',

        name: 'FrontendCases',

        component: () => import('@/views/frontend/cases/index.vue'),

        meta: { title: '应用案例' }

      },

      // 关于我们

      {

        path: 'about',

        name: 'FrontendAbout',

        component: () => import('@/views/frontend/about/index.vue'),

        meta: { title: '关于我们' }

      },

      // 法律页

      {

        path: 'legal/privacy',

        name: 'LegalPrivacy',

        component: () => import('@/views/frontend/legal/privacy.vue'),

        meta: { title: '隐私政策' }

      },

      {

        path: 'legal/terms',

        name: 'LegalTerms',

        component: () => import('@/views/frontend/legal/terms.vue'),

        meta: { title: '用户协议' }

      },

      // 会员登录

      {

        path: 'user/login',

        name: 'MemberLogin',

        component: () => import('@/views/frontend/member/login.vue'),

        meta: { title: '登录' }

      },

      // 会员注册

      {

        path: 'user/register',

        name: 'MemberRegister',

        component: () => import('@/views/frontend/member/register.vue'),

        meta: { title: '注册' }

      },

      // 会员中心壳层（M8 统一侧栏）

      {

        path: 'user',

        component: () => import('@/views/frontend/member/layout.vue'),

        meta: { title: '用户中心', requiresAuth: true },

        redirect: '/user/overview',

        children: [

          {

            path: 'overview',

            name: 'MemberOverview',

            component: memberCenter,

            meta: { title: '账户概览', requiresAuth: true }

          },

          {

            path: 'checkin',

            name: 'MemberCheckin',

            component: memberCenter,

            meta: { title: '每日签到', requiresAuth: true }

          },

          {

            path: 'profile',

            name: 'MemberProfile',

            component: memberCenter,

            meta: { title: '个人资料', requiresAuth: true }

          },

          {

            path: 'password',

            name: 'MemberPassword',

            component: memberCenter,

            meta: { title: '修改密码', requiresAuth: true }

          },

          {

            path: 'points',

            name: 'MemberPoints',

            component: memberCenter,

            meta: { title: '积分记录', requiresAuth: true }

          },

          {

            path: 'balance',

            name: 'MemberBalance',

            component: memberCenter,

            meta: { title: '余额记录', requiresAuth: true }

          },

          {

            path: 'notification',

            name: 'MemberNotification',

            component: memberCenter,

            meta: { title: '系统通知', requiresAuth: true }

          },

          // 合规子页面

          {
            path: 'compliance/profile',
            name: 'MemberComplianceProfile',
            component: () => import('@/views/frontend/compliance/profile.vue'),
            meta: { title: '企业画像', requiresAuth: true },
            beforeEnter: complianceRouteGuard
          },

          {
            path: 'compliance/calendar',
            name: 'MemberComplianceCalendar',
            component: () => import('@/views/frontend/compliance/calendar.vue'),
            meta: { title: '合规日历', requiresAuth: true },
            beforeEnter: complianceRouteGuard
          },

          {
            path: 'compliance/tasks/:id',
            name: 'MemberComplianceTaskDetail',
            component: () => import('@/views/frontend/compliance/task-detail.vue'),
            meta: { title: '合规任务详情', requiresAuth: true },
            beforeEnter: complianceRouteGuard
          },

          {
            path: 'compliance/documents',
            name: 'MemberComplianceDocuments',
            component: () => import('@/views/frontend/compliance/documents.vue'),
            meta: { title: '经营资料库', requiresAuth: true },
            beforeEnter: complianceRouteGuard
          },

          {
            path: 'compliance/risks',
            name: 'MemberComplianceRisks',
            component: () => import('@/views/frontend/compliance/risks.vue'),
            meta: { title: '风险中心', requiresAuth: true },
            beforeEnter: complianceRouteGuard
          },

          {
            path: 'compliance/reports',
            name: 'MemberComplianceReports',
            component: () => import('@/views/frontend/compliance/reports.vue'),
            meta: { title: '月度体检', requiresAuth: true },
            beforeEnter: complianceRouteGuard
          },

          {

            path: 'compliance/diagnosis',

            name: 'MemberComplianceDiagnosis',

            component: () => import('@/views/frontend/compliance/diagnosis.vue'),

            meta: { title: '诊断历史', requiresAuth: true }

          },

          {

            path: 'compliance/plan',

            name: 'MemberCompliancePlan',

            component: () => import('@/views/frontend/compliance/plan.vue'),

            meta: { title: '方案与签约', requiresAuth: true }

          },

          {

            path: 'compliance/opc',

            name: 'MemberComplianceOpc',

            component: () => import('@/views/frontend/compliance/opc.vue'),

            meta: { title: '主体设立进度', requiresAuth: true }

          },

          {

            path: 'compliance/income',

            name: 'MemberComplianceIncome',

            component: () => import('@/views/frontend/compliance/income.vue'),

            meta: { title: '收入台账', requiresAuth: true },

            beforeEnter: complianceRouteGuard

          },

          {

            path: 'compliance/expense',

            name: 'MemberComplianceExpense',

            component: () => import('@/views/frontend/compliance/expense.vue'),

            meta: { title: '费用台账', requiresAuth: true },

            beforeEnter: complianceRouteGuard

          },

          {

            path: 'compliance/ledger',

            name: 'MemberComplianceLedger',

            component: () => import('@/views/frontend/compliance/ledger.vue'),

            meta: { title: '利润报表', requiresAuth: true },

            beforeEnter: complianceRouteGuard

          },

          {

            path: 'compliance/filing',

            name: 'MemberComplianceFiling',

            component: () => import('@/views/frontend/compliance/filing.vue'),

            meta: { title: '报税中心', requiresAuth: true },

            beforeEnter: complianceRouteGuard

          },

          {

            path: 'compliance/tax',

            name: 'MemberComplianceTax',

            component: () => import('@/views/frontend/compliance/tax.vue'),

            meta: { title: '申报管理', requiresAuth: true },

            beforeEnter: complianceRouteGuard

          },

          {

            path: 'compliance/social-guide',

            name: 'MemberComplianceSocialGuide',

            component: () => import('@/views/frontend/compliance/social-guide.vue'),

            meta: { title: '社保指引', requiresAuth: true },

            beforeEnter: complianceRouteGuard

          },

          {

            path: 'compliance/social-consult',

            name: 'MemberComplianceSocialConsult',

            component: () => import('@/views/frontend/compliance/social-consult.vue'),

            meta: { title: '社保咨询', requiresAuth: true },

            beforeEnter: complianceRouteGuard

          },

          {

            path: 'compliance/statement',

            name: 'MemberComplianceStatement',

            component: () => import('@/views/frontend/compliance/statement.vue'),

            meta: { title: '对账单', requiresAuth: true },

            beforeEnter: complianceRouteGuard

          }

        ]

      },

    ]

  }

]



export default frontendRoutes
