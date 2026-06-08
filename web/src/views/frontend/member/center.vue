<!-- +----------------------------------------------------------------------
  | XYGo Admin [ Vue3 + GoFrame 企业级中后台管理系统 ]
  +----------------------------------------------------------------------
  | Copyright (c) 2026 大连星韵网络科技有限公司 All rights reserved.
  +----------------------------------------------------------------------
  | Licensed ( https://opensource.org/licenses/MIT )
  +----------------------------------------------------------------------
  | Author: 喜羊羊 <751300685@qq.com>
  +---------------------------------------------------------------------- -->
<template>
  <!-- 会员中心已禁用提示 -->
  <main v-if="!memberCenterOpen" class="pt-32 pb-8 px-6 flex items-center justify-center min-h-[60vh]">
    <div class="w-full max-w-md">
      <div class="bg-white/70 backdrop-blur-2xl rounded-[48px] shadow-clay-card border border-[#d1d9e6]/40 p-10 md:p-12 text-center">
        <div class="w-20 h-20 rounded-[24px] bg-[#f0f3f8] shadow-clay-pressed flex items-center justify-center mx-auto mb-6">
          <ArtSvgIcon icon="ri:lock-2-line" class="text-[36px] text-clay-muted" />
        </div>
        <h2 class="font-heading font-black text-2xl text-clay-foreground mb-3">会员中心已关闭</h2>
        <p class="text-clay-muted font-medium leading-relaxed">会员中心已禁用，请联系网站管理员开启。</p>
        <RouterLink to="/" class="inline-block mt-8 px-8 py-3 rounded-2xl bg-white shadow-clay-btn hover:shadow-clay-btn-hover font-bold text-clay-foreground active:scale-95 transition-all">
          返回首页
        </RouterLink>
      </div>
    </div>
  </main>

  <div v-else>

        <!-- 1. 账户概览（对齐 homesite/user-dashboard.html） -->
        <div v-if="activeMenu === 'overview'" class="h-full flex flex-col gap-8">
          <!-- 账户信息卡片 -->
          <section class="bg-white/70 backdrop-blur-xl rounded-[48px] shadow-clay-card border border-[#d1d9e6]/40 p-8 md:p-10 relative overflow-hidden group">
            <div class="absolute -bottom-20 -right-20 w-64 h-64 rounded-full bg-blue-500/5 blur-3xl group-hover:scale-110 transition-transform"></div>
            <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-6 mb-10 relative z-10">
              <div>
                <h3 class="text-sm font-black text-clay-muted uppercase tracking-wider mb-2">账户信息</h3>
                <div class="flex items-center gap-4">
                  <div class="w-16 h-16 rounded-2xl bg-white shadow-clay-pressed flex items-center justify-center text-3xl">👋</div>
                  <div>
                    <h1 class="font-heading font-black text-2xl text-clay-foreground">{{ userInfo.nickname || userInfo.username }}，{{ greeting }}！</h1>
                    <p class="text-clay-muted font-medium">欢迎回到 {{ siteName }} 门户中心</p>
                  </div>
                </div>
              </div>
              <RouterLink to="/user/profile" class="px-8 py-3 rounded-2xl bg-white shadow-clay-btn hover:shadow-clay-btn-hover active:scale-95 transition-all duration-300 font-bold text-clay-foreground flex items-center gap-2">
                <ArtSvgIcon icon="ri:user-line" class="text-lg" />
                个人资料
              </RouterLink>
            </div>
            <div class="grid sm:grid-cols-2 md:grid-cols-4 gap-6 relative z-10">
              <div class="p-6 rounded-[32px] bg-[#f0f3f8] shadow-clay-pressed">
                <span class="block text-xs font-black text-clay-muted uppercase tracking-widest mb-3">我的积分</span>
                <div class="flex items-end gap-2">
                  <span class="text-3xl font-black text-clay-accent">{{ userInfo.score ?? 0 }}</span>
                  <span class="text-xs font-bold text-clay-muted mb-1">分</span>
                </div>
              </div>
              <div class="p-6 rounded-[32px] bg-[#f0f3f8] shadow-clay-pressed">
                <span class="block text-xs font-black text-clay-muted uppercase tracking-widest mb-3">账户余额</span>
                <div class="flex items-end gap-2">
                  <span class="text-3xl font-black text-clay-success">{{ formatMoney(userInfo.money) }}</span>
                  <span class="text-xs font-bold text-clay-muted mb-1">元</span>
                </div>
              </div>
              <div class="p-6 rounded-[32px] bg-[#f0f3f8] shadow-clay-pressed sm:col-span-2">
                <div class="flex justify-between items-start">
                  <div>
                    <span class="block text-xs font-black text-clay-muted uppercase tracking-widest mb-3">最后登录</span>
                    <div class="font-bold text-clay-foreground">{{ formatTimestamp(userInfo.lastLoginAt) }}</div>
                    <div class="text-xs text-clay-muted mt-1">IP: {{ userInfo.lastLoginIp || '-' }}</div>
                  </div>
                  <div class="w-12 h-12 rounded-2xl bg-white shadow-clay-btn flex items-center justify-center text-xl">📍</div>
                </div>
              </div>
            </div>
          </section>

          <!-- 合规仪表盘（P-07） -->
          <section class="bg-white/70 backdrop-blur-xl rounded-[48px] shadow-clay-card border border-[#d1d9e6]/40 p-8 md:p-10 relative overflow-hidden">
            <div class="flex justify-between items-center mb-8">
              <div>
                <h3 class="font-heading font-black text-xl text-clay-foreground">合规服务概览</h3>
                <p class="text-xs text-clay-muted mt-1">OPC 状态与本月经营摘要</p>
              </div>
            </div>

            <!-- 未签约空态 -->
            <div v-if="!complianceState.hasActiveOrder" class="py-12 text-center">
              <div class="w-20 h-20 rounded-[24px] bg-[#f0f3f8] shadow-clay-pressed flex items-center justify-center mx-auto mb-6">
                <ArtSvgIcon icon="ri:shield-check-line" class="text-[36px] text-clay-accent" />
              </div>
              <p class="text-clay-muted font-bold mb-6">
                {{ complianceState.hasPendingOrder
                  ? '您已选择套餐，请继续完成风险告知与电子签约'
                  : '尚未签约合规服务，完成免费诊断后可选择方案签约' }}
              </p>
              <div class="flex flex-col sm:flex-row items-center justify-center gap-4">
                <RouterLink
                  to="/user/compliance/plan"
                  class="inline-flex items-center gap-2 px-8 py-3 rounded-2xl bg-gradient-to-br from-blue-400 to-blue-600 text-white font-bold shadow-clay-btn"
                >
                  {{ complianceState.hasPendingOrder ? '继续签约' : '方案与签约' }}
                  <ArtSvgIcon icon="ri:arrow-right-line" />
                </RouterLink>
                <RouterLink
                  v-if="!complianceState.hasPendingOrder"
                  to="/diagnosis"
                  class="inline-flex items-center gap-2 px-8 py-3 rounded-2xl bg-white font-bold text-clay-foreground shadow-clay-btn"
                >
                  开始免费诊断
                </RouterLink>
              </div>
            </div>

            <template v-else>
              <div class="grid sm:grid-cols-3 gap-6 mb-8">
                <div class="p-6 rounded-[32px] bg-[#f0f3f8] shadow-clay-pressed">
                  <span class="block text-xs font-black text-clay-muted uppercase tracking-widest mb-3">OPC 状态</span>
                  <span class="text-xl font-black" :class="complianceState.opcStatus === 'active' ? 'text-clay-success' : 'text-clay-accent'">
                    {{ opcStatusLabel }}
                  </span>
                </div>
                <div class="p-6 rounded-[32px] bg-[#f0f3f8] shadow-clay-pressed">
                  <span class="block text-xs font-black text-clay-muted uppercase tracking-widest mb-3">本月收入</span>
                  <span class="text-xl font-black text-clay-foreground">¥ {{ formatComplianceMoney(monthSummary.revenue) }}</span>
                </div>
                <div class="p-6 rounded-[32px] bg-[#f0f3f8] shadow-clay-pressed">
                  <span class="block text-xs font-black text-clay-muted uppercase tracking-widest mb-3">本月利润</span>
                  <span class="text-xl font-black text-clay-success">¥ {{ formatComplianceMoney(monthSummary.profit) }}</span>
                </div>
              </div>

              <!-- 待办列表 -->
              <div v-if="todoList.length > 0" class="mb-8">
                <h4 class="text-sm font-black text-clay-muted uppercase tracking-widest mb-4">待办事项</h4>
                <ul class="space-y-3">
                  <li
                    v-for="(todo, idx) in todoList"
                    :key="idx"
                    class="flex items-center justify-between gap-4 p-4 rounded-2xl"
                    :class="todo.urgent ? 'bg-red-50 border border-red-100' : 'bg-[#f0f3f8] shadow-clay-pressed'"
                  >
                    <div class="flex items-center gap-3">
                      <ArtSvgIcon :icon="todo.urgent ? 'ri:error-warning-line' : 'ri:checkbox-circle-line'" :class="todo.urgent ? 'text-red-500' : 'text-clay-accent'" />
                      <span class="font-bold text-sm text-clay-foreground">{{ todo.text }}</span>
                    </div>
                    <RouterLink v-if="todo.path" :to="todo.path" class="text-xs font-bold text-clay-accent hover:underline shrink-0">
                      去处理 →
                    </RouterLink>
                  </li>
                </ul>
              </div>

              <!-- 快捷入口 -->
              <div class="flex flex-wrap gap-3">
                <RouterLink
                  v-for="entry in quickEntries"
                  :key="entry.path"
                  :to="entry.path"
                  class="px-5 py-2.5 rounded-xl bg-white shadow-clay-btn text-sm font-bold text-clay-foreground hover:shadow-clay-btn-hover transition-all"
                >
                  {{ entry.label }}
                </RouterLink>
              </div>
            </template>
          </section>

          <!-- 增长趋势统计图表（对齐 homesite） -->
          <section class="bg-white/70 backdrop-blur-xl rounded-[48px] shadow-clay-card border border-[#d1d9e6]/40 p-8 md:p-10 relative overflow-hidden group flex-1 flex flex-col">
            <div class="flex justify-between items-center mb-10">
              <div>
                <h3 class="font-heading font-black text-xl text-clay-foreground">增长趋势统计</h3>
                <p class="text-xs text-clay-muted mt-1">最近 7 天账户积分与余额增长情况</p>
              </div>
              <div class="flex gap-3">
                <div class="flex items-center gap-2 px-4 py-2 rounded-2xl bg-white shadow-clay-btn text-[10px] font-bold text-blue-500">
                  <span class="w-2 h-2 rounded-full bg-blue-500"></span> 积分增长
                </div>
                <div class="flex items-center gap-2 px-4 py-2 rounded-2xl bg-white shadow-clay-btn text-[10px] font-bold text-green-500">
                  <span class="w-2 h-2 rounded-full bg-green-500"></span> 余额增长
                </div>
              </div>
            </div>
            <!-- 拟态柱状图 -->
            <div class="flex-1 w-full relative min-h-[320px] flex flex-col justify-end">
              <div class="absolute inset-0 flex flex-col justify-between py-6 pointer-events-none">
                <div v-for="i in 5" :key="i" class="w-full h-[1px] bg-gray-200/30"></div>
              </div>
              <div class="relative z-10 flex justify-between items-end h-[240px] px-4">
                <div v-for="(day, index) in chartDays" :key="index" class="flex flex-col items-center gap-4 group/bar">
                  <div class="flex items-end gap-3 h-[200px]">
                    <!-- 积分柱 -->
                    <div class="w-4 rounded-full bg-gradient-to-t from-blue-500 to-blue-300 shadow-clay-btn relative overflow-hidden group/item transition-all duration-500"
                         :style="{ height: chartScoreHeights[index] + '%' }">
                      <div class="absolute inset-0 opacity-0 group-hover/item:opacity-100 transition-opacity bg-white/20"></div>
                      <div class="absolute top-1 left-1 right-1 h-2 rounded-full bg-white/30 blur-[1px]"></div>
                      <div class="absolute -top-10 left-1/2 -translate-x-1/2 bg-blue-600 text-white text-[10px] px-2 py-1 rounded-lg opacity-0 group-hover/item:opacity-100 transition-all duration-300 pointer-events-none whitespace-nowrap shadow-lg">
                        +{{ chartScoreValues[index] }} 积分
                      </div>
                    </div>
                    <!-- 余额柱 -->
                    <div class="w-4 rounded-full bg-gradient-to-t from-green-500 to-green-300 shadow-clay-btn relative overflow-hidden group/item transition-all duration-500"
                         :style="{ height: chartMoneyHeights[index] + '%' }">
                      <div class="absolute inset-0 opacity-0 group-hover/item:opacity-100 transition-opacity bg-white/20"></div>
                      <div class="absolute top-1 left-1 right-1 h-2 rounded-full bg-white/30 blur-[1px]"></div>
                      <div class="absolute -top-10 left-1/2 -translate-x-1/2 bg-green-600 text-white text-[10px] px-2 py-1 rounded-lg opacity-0 group-hover/item:opacity-100 transition-all duration-300 pointer-events-none whitespace-nowrap shadow-lg">
                        +{{ chartMoneyValues[index] }}.00 元
                      </div>
                    </div>
                  </div>
                  <span class="text-[10px] font-bold text-clay-muted group-hover/bar:text-clay-accent transition-colors">{{ day }}</span>
                </div>
              </div>
            </div> 
          </section>
        </div>

        <!-- 2. 每日签到（对齐 homesite：7天日历+签到按钮） -->
        <div v-if="activeMenu === 'checkin'" class="animate-in">
          <section class="bg-white/70 backdrop-blur-xl rounded-[48px] shadow-clay-card border border-[#d1d9e6]/40 p-8 md:p-12 text-center">
            <div class="w-24 h-24 rounded-[32px] bg-gradient-to-br from-amber-300 to-amber-500 shadow-clay-btn flex items-center justify-center text-4xl mx-auto mb-6">✨</div>
            <h2 class="font-heading font-black text-3xl text-clay-foreground mb-4">今日签到</h2>
            <p class="text-clay-muted mb-10 max-w-sm mx-auto">
              每日签到可获得随机积分奖励，连续签到更有额外大奖！
              <span v-if="checkinData.continuousDays > 0" class="block mt-2 text-clay-accent font-bold">
                已连续签到 {{ checkinData.continuousDays }} 天
              </span>
            </p>
            <!-- 7天签到日历 -->
            <div class="grid grid-cols-7 gap-3 mb-12 max-w-lg mx-auto">
              <div v-for="(day, idx) in checkinData.weekDays" :key="idx"
                   class="aspect-square rounded-2xl flex flex-col items-center justify-center gap-1 transition-all"
                   :class="day.checked ? 'bg-gradient-to-br from-blue-400 to-blue-600 text-white shadow-clay-btn' : 'bg-[#f0f3f8] shadow-clay-pressed text-clay-muted opacity-50'">
                <span class="text-[10px] font-black uppercase">{{ day.date.slice(5) }}</span>
                <span class="text-lg font-black">{{ idx + 1 }}</span>
                <span v-if="day.checked" class="text-[10px]">+{{ day.score }}</span>
              </div>
            </div>
            <button
              class="w-full max-w-xs py-5 rounded-[24px] text-white text-xl font-black shadow-clay-btn hover:shadow-clay-btn-hover active:scale-95 transition-all"
              :class="checkinData.todayChecked ? 'bg-gray-400 cursor-not-allowed' : 'bg-gradient-to-br from-blue-400 to-blue-600'"
              :disabled="checkinData.todayChecked || checkinLoading"
              @click="handleCheckin"
            >
              {{ checkinLoading ? '签到中...' : checkinData.todayChecked ? `已签到 (+${checkinData.todayScore}积分)` : '立即签到' }}
            </button>
          </section>
        </div>

        <!-- 3. 个人资料 -->
        <div v-if="activeMenu === 'profile'" class="animate-in">
          <section class="bg-white/70 backdrop-blur-xl rounded-[48px] shadow-clay-card border border-[#d1d9e6]/40 p-8 md:p-12">
            <h2 class="font-heading font-black text-2xl text-clay-foreground mb-8">个人资料设置</h2>
            <ElForm ref="profileFormRef" :model="profileForm" :rules="profileRules" label-position="top" class="grid md:grid-cols-2 gap-8">
              <div class="space-y-6">
                <ElFormItem label="昵称" prop="nickname">
                  <ElInput v-model="profileForm.nickname" size="large" class="clay-input" />
                </ElFormItem>
                <ElFormItem label="邮箱" prop="email">
                  <ElInput v-model="profileForm.email" size="large" class="clay-input" />
                </ElFormItem>
                <ElFormItem label="手机号">
                  <ElInput v-model="profileForm.mobile" size="large" class="clay-input" />
                </ElFormItem>
                <ElFormItem label="性别">
                  <ElRadioGroup v-model="profileForm.gender">
                    <ElRadio :value="1">男</ElRadio>
                    <ElRadio :value="2">女</ElRadio>
                    <ElRadio :value="0">保密</ElRadio>
                  </ElRadioGroup>
                </ElFormItem>
              </div>
              <div class="flex flex-col items-center justify-center p-8 rounded-[40px] bg-[#f0f3f8] shadow-clay-pressed border-2 border-dashed border-gray-200">
                <div class="w-32 h-32 rounded-[40px] bg-white shadow-clay-btn p-1 mb-6">
                  <ElAvatar :size="120" :src="profileForm.avatar" class="!rounded-[36px] !w-full !h-full">
                    {{ profileForm.nickname?.charAt(0) || 'U' }}
                  </ElAvatar>
                </div>
                <input ref="avatarInputRef" type="file" accept="image/jpeg,image/png,image/gif" class="hidden" @change="handleAvatarUpload" />
                <button type="button" class="px-6 py-2 rounded-xl bg-white shadow-clay-btn text-xs font-bold text-clay-accent transition-all" @click="avatarInputRef?.click()">
                  {{ avatarUploading ? '上传中...' : '更换头像' }}
                </button>
                <p class="mt-4 text-[10px] text-clay-muted text-center">支持 JPG, PNG, GIF 格式<br>最大 2MB</p>
              </div>
            </ElForm>
            <div class="mt-12 flex justify-end gap-4">
              <button class="px-8 py-3 rounded-2xl bg-white shadow-clay-btn font-bold text-clay-muted active:scale-95 transition-all">取消</button>
              <button class="px-10 py-3 rounded-2xl bg-gradient-to-br from-blue-400 to-blue-600 text-white shadow-clay-btn font-bold active:scale-95 transition-all" :disabled="saving" @click="handleSaveProfile">
                {{ saving ? '保存中...' : '保存修改' }}
              </button>
            </div>
          </section>
        </div>

        <!-- 4. 修改密码 -->
        <div v-if="activeMenu === 'password'" class="animate-in">
          <section class="bg-white/70 backdrop-blur-xl rounded-[48px] shadow-clay-card border border-[#d1d9e6]/40 p-8 md:p-12">
            <h2 class="font-heading font-black text-2xl text-clay-foreground mb-8">安全设置 - 修改密码</h2>
            <ElForm ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-position="top" class="max-w-md space-y-6">
              <ElFormItem label="当前密码" prop="oldPassword">
                <ElInput v-model="passwordForm.oldPassword" type="password" show-password placeholder="请输入当前密码" size="large" class="clay-input" />
              </ElFormItem>
              <ElFormItem label="新密码" prop="newPassword">
                <ElInput v-model="passwordForm.newPassword" type="password" show-password placeholder="不少于6位" size="large" class="clay-input" />
              </ElFormItem>
              <ElFormItem label="确认新密码" prop="confirmPassword">
                <ElInput v-model="passwordForm.confirmPassword" type="password" show-password placeholder="再次输入新密码" size="large" class="clay-input" />
              </ElFormItem>
              <div class="pt-4">
                <button class="w-full py-4 rounded-2xl bg-gradient-to-br from-blue-400 to-blue-600 text-white shadow-clay-btn font-bold active:scale-95 transition-all" :disabled="changingPassword" @click="handleChangePassword">
                  {{ changingPassword ? '更新中...' : '更新密码' }}
                </button>
              </div>
            </ElForm>
          </section>
        </div>

        <!-- 5. 积分记录（对齐 homesite 拟态表格） -->
        <div v-if="activeMenu === 'points'" class="animate-in">
          <section class="bg-white/70 backdrop-blur-xl rounded-[48px] shadow-clay-card border border-[#d1d9e6]/40 p-8 md:p-10">
            <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
              <h2 class="font-heading font-black text-2xl text-clay-foreground">积分记录</h2>
              <div class="px-6 py-3 rounded-2xl bg-[#f0f3f8] shadow-clay-pressed font-bold text-clay-muted">
                当前积分：<span class="text-clay-accent">{{ userInfo.score ?? 0 }}</span>
              </div>
            </div>
            <div v-if="scoreLogList.length === 0 && !scoreLogLoading" class="py-20 text-center">
              <div class="w-24 h-24 rounded-full bg-white shadow-clay-btn flex items-center justify-center mb-6 mx-auto">
                <ArtSvgIcon icon="ri:coin-line" class="text-[40px] text-clay-muted opacity-50" />
              </div>
              <p class="text-clay-muted font-bold text-lg">暂无积分变动记录</p>
            </div>
            <div v-else class="overflow-hidden rounded-[32px] bg-[#f0f3f8] shadow-clay-pressed p-2">
              <table class="w-full text-left">
                <thead>
                  <tr class="text-xs font-black text-clay-muted uppercase tracking-widest border-b border-gray-100">
                    <th class="px-6 py-4">变动原因</th>
                    <th class="px-6 py-4 text-center">变动数值</th>
                    <th class="px-6 py-4 text-right">变动时间</th>
                  </tr>
                </thead>
                <tbody class="text-sm font-bold">
                  <tr v-for="item in scoreLogList" :key="item.id" class="hover:bg-white/50 transition-colors">
                    <td class="px-6 py-5 text-clay-foreground">{{ item.memo || '积分变动' }}</td>
                    <td class="px-6 py-5 text-center" :class="item.score > 0 ? 'text-clay-success' : 'text-red-500'">
                      {{ item.score > 0 ? '+' : '' }}{{ item.score }}
                    </td>
                    <td class="px-6 py-5 text-right text-clay-muted font-medium">{{ formatTimestamp(item.createdAt) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <!-- 分页 -->
            <div v-if="scoreLogTotal > scoreLogPageSize" class="pt-6 flex justify-center">
              <button v-if="scoreLogPage * scoreLogPageSize < scoreLogTotal"
                class="flex items-center gap-2 text-sm font-bold text-clay-muted hover:text-clay-accent transition-colors"
                @click="loadScoreLog(scoreLogPage + 1)">
                加载更多 <ArtSvgIcon icon="ri:arrow-down-s-line" class="text-base" />
              </button>
            </div>
          </section>
        </div>

        <!-- 6. 余额记录 -->
        <div v-if="activeMenu === 'balance'" class="animate-in">
          <section class="bg-white/70 backdrop-blur-xl rounded-[48px] shadow-clay-card border border-[#d1d9e6]/40 p-8 md:p-10">
            <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
              <h2 class="font-heading font-black text-2xl text-clay-foreground">余额变动记录</h2>
              <div class="px-6 py-3 rounded-2xl bg-[#f0f3f8] shadow-clay-pressed font-bold text-clay-muted">
                可用余额：<span class="text-clay-success">{{ formatMoney(userInfo.money) }}</span> 元
              </div>
            </div>
            <div v-if="moneyLogList.length === 0 && !moneyLogLoading" class="py-20 flex flex-col items-center justify-center text-center">
              <div class="w-24 h-24 rounded-full bg-white shadow-clay-btn flex items-center justify-center mb-6">
                <ArtSvgIcon icon="ri:wallet-3-line" class="text-[40px] text-clay-muted opacity-50" />
              </div>
              <p class="text-clay-muted font-bold text-lg">暂无余额变动记录</p>
            </div>
            <div v-else class="overflow-hidden rounded-[32px] bg-[#f0f3f8] shadow-clay-pressed p-2">
              <table class="w-full text-left">
                <thead>
                  <tr class="text-xs font-black text-clay-muted uppercase tracking-widest border-b border-gray-100">
                    <th class="px-6 py-4">变动原因</th>
                    <th class="px-6 py-4 text-center">变动金额</th>
                    <th class="px-6 py-4 text-right">变动时间</th>
                  </tr>
                </thead>
                <tbody class="text-sm font-bold">
                  <tr v-for="item in moneyLogList" :key="item.id" class="hover:bg-white/50 transition-colors">
                    <td class="px-6 py-5 text-clay-foreground">{{ item.memo || '余额变动' }}</td>
                    <td class="px-6 py-5 text-center" :class="item.money > 0 ? 'text-clay-success' : 'text-red-500'">
                      {{ item.money > 0 ? '+' : '' }}{{ formatCent(item.money) }} 元
                    </td>
                    <td class="px-6 py-5 text-right text-clay-muted font-medium">{{ formatTimestamp(item.createdAt) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div v-if="moneyLogTotal > moneyLogPageSize" class="pt-6 flex justify-center">
              <button v-if="moneyLogPage * moneyLogPageSize < moneyLogTotal"
                class="flex items-center gap-2 text-sm font-bold text-clay-muted hover:text-clay-accent transition-colors"
                @click="loadMoneyLog(moneyLogPage + 1)">
                加载更多 <ArtSvgIcon icon="ri:arrow-down-s-line" class="text-base" />
              </button>
            </div>
          </section>
        </div>

        <!-- 7. 系统通知（对齐 homesite：富通知卡片+未读标记） -->
        <div v-if="activeMenu === 'notification'" class="animate-in">
          <section class="bg-white/70 backdrop-blur-xl rounded-[48px] shadow-clay-card border border-[#d1d9e6]/40 p-8 md:p-12">
            <div class="flex justify-between items-center mb-10">
              <div class="flex items-center gap-4">
                <h2 class="font-heading font-black text-2xl text-clay-foreground">系统通知</h2>
                <span v-if="noticeUnread > 0" class="px-3 py-1 rounded-full bg-blue-500 text-white text-[10px] font-black">{{ noticeUnread }} 未读</span>
              </div>
              <button
                class="px-6 py-2 rounded-xl bg-[#f0f3f8] shadow-clay-pressed text-xs font-bold text-clay-muted hover:bg-white hover:shadow-clay-btn transition-all"
                @click="handleReadAllNotice"
              >全部已读</button>
            </div>
            <div v-if="noticeList.length === 0 && !noticeLoading" class="py-20 text-center">
              <div class="w-24 h-24 rounded-full bg-white shadow-clay-btn flex items-center justify-center mb-6 mx-auto">
                <ArtSvgIcon icon="ri:notification-line" class="text-[40px] text-clay-muted opacity-50" />
              </div>
              <p class="text-clay-muted font-bold text-lg">暂无通知</p>
            </div>
            <div v-else class="space-y-8">
              <div v-for="item in noticeList" :key="item.id"
                   class="group relative bg-[#f0f3f8] shadow-clay-pressed rounded-[32px] p-6 md:p-8 hover:bg-white hover:shadow-clay-card transition-all duration-500 cursor-pointer"
                   @click="handleReadNotice(item)">
                <!-- 未读标记 -->
                <div v-if="!item.isRead" class="absolute top-8 left-8 w-2.5 h-2.5 rounded-full bg-blue-500 shadow-[0_0_10px_rgba(59,130,246,0.5)]"></div>
                <div :class="{ 'ml-6': !item.isRead }">
                  <div class="flex flex-wrap items-center gap-3 mb-3">
                    <span class="px-3 py-1 rounded-lg bg-white shadow-clay-btn text-[10px] font-black text-clay-accent uppercase">{{ noticeTypeLabel(item.type) }}</span>
                    <h3 class="font-bold text-lg text-clay-foreground group-hover:text-blue-600 transition-colors">{{ item.title }}</h3>
                  </div>
                  <div class="text-sm text-clay-muted leading-relaxed mb-6" v-html="item.content"></div>
                  <div class="flex items-center justify-between">
                    <div class="flex items-center gap-2 text-xs font-bold text-clay-muted">
                      <ArtSvgIcon icon="ri:time-line" class="text-sm" />
                      {{ formatTimestamp(item.createdAt) }}
                    </div>
                  </div>
                </div>
              </div>
              <!-- 加载更多 -->
              <div v-if="noticeTotal > noticePageSize" class="pt-4 flex justify-center">
                <button v-if="noticePage * noticePageSize < noticeTotal"
                  class="flex items-center gap-2 text-sm font-bold text-clay-muted hover:text-clay-accent transition-colors"
                  @click="loadNoticeList(noticePage + 1)">
                  查看更多通知 <ArtSvgIcon icon="ri:arrow-down-s-line" class="text-base" />
                </button>
              </div>
            </div>
          </section>
        </div>

  </div>
</template>

<script setup lang="ts">
import { useMemberStore } from '@/store/modules/member'
import { useSiteStore } from '@/store/modules/site'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import {
  updateMemberProfile, changeMemberPassword, getMemberInfo,
  getCheckinInfo, doCheckin,
  getScoreLogList,
  getMoneyLogList,
  getNoticeList, markNoticeRead, markAllNoticeRead,
} from '@/api/frontend'
import type { CheckinDayItem, ScoreLogItem, MoneyLogItem, NoticeItem } from '@/api/frontend'
import {
  getCompliancePlanState,
  getProfitSummary,
  getTaxCalendar,
  getBankUnmatched,
  formatMoney as formatComplianceMoney,
} from '@/api/frontend/compliance/member'
import type { CompliancePlanState } from '@/config/complianceMenu'
import { formatTimestamp } from '@/utils/time'

defineOptions({ name: 'UserCenter' })

const route = useRoute()
const memberStore = useMemberStore()
const siteStore = useSiteStore()
const siteName = computed(() => siteStore.getSiteName())
const userInfo = computed(() => memberStore.getMemberInfo)
const memberCenterOpen = computed(() => siteStore.isUserCenterEnabled())

// 页面加载时主动拉取数据（解决重启/刷新后数据丢失）
onMounted(async () => {
  if (!memberStore.getIsLogin) return
  // 拉取会员信息
  try {
    const res = await getMemberInfo()
    if (res) memberStore.setMemberInfo(res)
  } catch {
    // token 过期或无效，自动登出
    memberStore.logOut()
    return
  }
  loadComplianceOverview()
  loadGrowthChart()
})

// 路由驱动当前面板（/user/overview → overview）
const activeMenu = computed(() => {
  const seg = route.path.replace(/^\/user\/?/, '').split('/')[0]
  return seg || 'overview'
})

// ===== 合规概览数据（P-07） =====
const complianceState = ref<CompliancePlanState>({ hasActiveOrder: false, opcStatus: 'none' })
const monthSummary = reactive({ revenue: 0, profit: 0 })
const todoList = ref<{ text: string; path?: string; urgent?: boolean }[]>([])

const opcStatusLabel = computed(() => {
  if (complianceState.value.opcStatus === 'active') return '已激活'
  if (complianceState.value.hasActiveOrder) return '设立中'
  return '未签约'
})

const quickEntries = computed(() => {
  if (complianceState.value.opcStatus === 'active') {
    return [
      { label: '记收入', path: '/user/compliance/income' },
      { label: '记费用', path: '/user/compliance/expense' },
      { label: '看对账单', path: '/user/compliance/statement' },
    ]
  }
  if (complianceState.value.hasActiveOrder) {
    return [{ label: '查看 OPC 进度', path: '/user/compliance/opc' }]
  }
  if (complianceState.value.hasPendingOrder) {
    return [{ label: '继续签约', path: '/user/compliance/plan' }]
  }
  return [{ label: '方案与签约', path: '/user/compliance/plan' }]
})

async function loadComplianceOverview() {
  try {
    complianceState.value = await getCompliancePlanState()
    if (!complianceState.value.hasActiveOrder) return

    const now = new Date()
    const period = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
    const todos: typeof todoList.value = []

    if (complianceState.value.opcStatus === 'active') {
      try {
        const profit = await getProfitSummary({ period, periodType: 'month' })
        monthSummary.revenue = profit?.revenue ?? 0
        monthSummary.profit = profit?.profit ?? 0
      } catch { /* ignore */ }

      try {
        const tax = await getTaxCalendar({ year: now.getFullYear(), month: now.getMonth() + 1 })
        for (const task of tax?.tasks || []) {
          if (task.status === 'pending' || task.status === 'overdue') {
            todos.push({
              text: `${task.taxTypeLabel}申报（截止 ${task.dueDate}）`,
              path: '/user/compliance/tax',
              urgent: task.status === 'overdue',
            })
          }
        }
      } catch { /* ignore */ }

      try {
        const bank = await getBankUnmatched({ month: period })
        if (bank?.count > 0) {
          todos.push({
            text: `有 ${bank.count} 笔银行流水待匹配入账`,
            path: '/user/compliance/income',
            urgent: true,
          })
        }
      } catch { /* ignore */ }
    }

    todoList.value = todos
  } catch { /* ignore */ }
}

// ===== 图表数据（最近 7 天积分/余额增长） =====
const chartDays = ref<string[]>([])
const chartScoreHeights = ref<number[]>([])
const chartScoreValues = ref<number[]>([])
const chartMoneyHeights = ref<number[]>([])
const chartMoneyValues = ref<number[]>([])

async function loadGrowthChart() {
  const days: string[] = []
  const scoreByDay: number[] = []
  const moneyByDay: number[] = []
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  for (let i = 6; i >= 0; i--) {
    const d = new Date(today)
    d.setDate(d.getDate() - i)
    days.push(`${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`)
    scoreByDay.push(0)
    moneyByDay.push(0)
  }

  try {
    const [scoreRes, moneyRes] = await Promise.all([
      getScoreLogList({ page: 1, pageSize: 200 }),
      getMoneyLogList({ page: 1, pageSize: 200 }),
    ])
    const dayKey = (ts: string) => {
      const d = new Date(ts)
      return `${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
    }
    const scoreMap = new Map(days.map((d) => [d, 0]))
    const moneyMap = new Map(days.map((d) => [d, 0]))

    for (const item of scoreRes?.list || []) {
      if (item.score <= 0) continue
      const key = dayKey(item.createdAt)
      if (scoreMap.has(key)) scoreMap.set(key, (scoreMap.get(key) || 0) + item.score)
    }
    for (const item of moneyRes?.list || []) {
      if (item.money <= 0) continue
      const key = dayKey(item.createdAt)
      if (moneyMap.has(key)) moneyMap.set(key, (moneyMap.get(key) || 0) + item.money / 100)
    }

    days.forEach((d, i) => {
      scoreByDay[i] = scoreMap.get(d) || 0
      moneyByDay[i] = moneyMap.get(d) || 0
    })
  } catch { /* ignore */ }

  const maxScore = Math.max(...scoreByDay, 1)
  const maxMoney = Math.max(...moneyByDay, 1)

  chartDays.value = days
  chartScoreValues.value = scoreByDay
  chartMoneyValues.value = moneyByDay.map((v) => Math.round(v))
  chartScoreHeights.value = scoreByDay.map((v) => Math.max(8, Math.round((v / maxScore) * 100)))
  chartMoneyHeights.value = moneyByDay.map((v) => Math.max(8, Math.round((v / maxMoney) * 100)))
}

// ===== 签到数据 =====
const checkinData = reactive({
  continuousDays: 0,
  todayChecked: false,
  todayScore: 0,
  weekDays: [] as CheckinDayItem[]
})
const checkinLoading = ref(false)

async function loadCheckinInfo() {
  try {
    const res = await getCheckinInfo()
    if (res) {
      checkinData.continuousDays = res.continuousDays
      checkinData.todayChecked = res.todayChecked
      checkinData.todayScore = res.todayScore
      checkinData.weekDays = res.weekDays || []
    }
  } catch { /* ignore */ }
}

async function handleCheckin() {
  checkinLoading.value = true
  try {
    const res = await doCheckin()
    if (res) {
      checkinData.todayChecked = true
      checkinData.todayScore = res.score
      checkinData.continuousDays = res.continuousDays
      ElMessage.success(`签到成功！获得 ${res.score} 积分`)
      // 刷新签到日历和会员信息（积分已变）
      await loadCheckinInfo()
      const info = await getMemberInfo()
      if (info) memberStore.setMemberInfo(info)
    }
  } catch { /* 拦截器已处理 */ } finally {
    checkinLoading.value = false
  }
}

// ===== 积分记录 =====
const scoreLogList = ref<ScoreLogItem[]>([])
const scoreLogLoading = ref(false)
const scoreLogPage = ref(1)
const scoreLogPageSize = 20
const scoreLogTotal = ref(0)

async function loadScoreLog(page = 1) {
  scoreLogLoading.value = true
  try {
    const res = await getScoreLogList({ page, pageSize: scoreLogPageSize })
    if (res) {
      if (page === 1) {
        scoreLogList.value = res.list || []
      } else {
        scoreLogList.value.push(...(res.list || []))
      }
      scoreLogPage.value = res.page
      scoreLogTotal.value = res.total
    }
  } catch { /* ignore */ } finally {
    scoreLogLoading.value = false
  }
}

// ===== 余额记录 =====
const moneyLogList = ref<MoneyLogItem[]>([])
const moneyLogLoading = ref(false)
const moneyLogPage = ref(1)
const moneyLogPageSize = 20
const moneyLogTotal = ref(0)

async function loadMoneyLog(page = 1) {
  moneyLogLoading.value = true
  try {
    const res = await getMoneyLogList({ page, pageSize: moneyLogPageSize })
    if (res) {
      if (page === 1) {
        moneyLogList.value = res.list || []
      } else {
        moneyLogList.value.push(...(res.list || []))
      }
      moneyLogPage.value = res.page
      moneyLogTotal.value = res.total
    }
  } catch { /* ignore */ } finally {
    moneyLogLoading.value = false
  }
}

// 分转元显示
const formatCent = (v: number) => (v / 100).toFixed(2)

// ===== 系统通知 =====
const noticeList = ref<NoticeItem[]>([])
const noticeLoading = ref(false)
const noticePage = ref(1)
const noticePageSize = 20
const noticeTotal = ref(0)
const noticeUnread = ref(0)

async function loadNoticeList(page = 1) {
  noticeLoading.value = true
  try {
    const res = await getNoticeList({ page, pageSize: noticePageSize })
    if (res) {
      if (page === 1) {
        noticeList.value = res.list || []
      } else {
        noticeList.value.push(...(res.list || []))
      }
      noticePage.value = res.page
      noticeTotal.value = res.total
      noticeUnread.value = res.unread
    }
  } catch { /* ignore */ } finally {
    noticeLoading.value = false
  }
}

async function handleReadNotice(item: NoticeItem) {
  if (!item.isRead) {
    try {
      await markNoticeRead(item.id)
      item.isRead = true
      noticeUnread.value = Math.max(0, noticeUnread.value - 1)
    } catch { /* ignore */ }
  }
}

async function handleReadAllNotice() {
  try {
    await markAllNoticeRead()
    noticeList.value.forEach(n => n.isRead = true)
    noticeUnread.value = 0
    ElMessage.success('已全部标为已读')
  } catch { /* ignore */ }
}

const noticeTypeLabel = (type: string) => {
  const map: Record<string, string> = { system: '系统', announce: '公告', feature: '功能', maintain: '维护' }
  return map[type] || type
}

// ===== 路由切换时按需加载数据 =====
watch(activeMenu, (menu) => {
  if (menu === 'overview') {
    loadComplianceOverview()
    loadGrowthChart()
  }
  if (menu === 'checkin') loadCheckinInfo()
  if (menu === 'points') loadScoreLog(1)
  if (menu === 'balance') loadMoneyLog(1)
  if (menu === 'notification') loadNoticeList(1)
}, { immediate: true })

// 问候语
const greeting = computed(() => {
  const h = new Date().getHours()
  if (h < 6) return '夜深了'
  if (h < 12) return '早上好'
  if (h < 14) return '中午好'
  if (h < 18) return '下午好'
  return '晚上好'
})

// 格式化金额
const formatMoney = (v: any) => {
  const n = Number(v) || 0
  return n.toFixed(2)
}

// ===== 个人资料表单 =====
const profileFormRef = ref<FormInstance>()
const profileForm = reactive({
  avatar: userInfo.value.avatar || '',
  nickname: userInfo.value.nickname || '',
  email: userInfo.value.email || '',
  mobile: userInfo.value.mobile || '',
  gender: userInfo.value.gender || 0,
})

const profileRules: FormRules = {
  nickname: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
  email: [{ type: 'email', message: '邮箱格式不正确', trigger: 'blur' }],
}

const saving = ref(false)
const avatarInputRef = ref<HTMLInputElement>()
const avatarUploading = ref(false)

const handleAvatarUpload = async (e: Event) => {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  if (file.size > 2 * 1024 * 1024) {
    ElMessage.warning('图片不能超过 2MB')
    return
  }
  avatarUploading.value = true
  try {
    const fd = new FormData()
    fd.append('file', file)
    const { memberRequest } = await import('@/utils/http')
    const res = await memberRequest.post<any>({ url: '/user/upload', data: fd, headers: { 'Content-Type': 'multipart/form-data' } })
    if (res?.url) {
      profileForm.avatar = res.url
      ElMessage.success('头像上传成功，请点击保存修改')
    }
  } catch { ElMessage.error('上传失败') } finally {
    avatarUploading.value = false
    if (avatarInputRef.value) avatarInputRef.value.value = ''
  }
}

const handleSaveProfile = async () => {
  if (!profileFormRef.value) return
  const valid = await profileFormRef.value.validate().catch(() => false)
  if (!valid) return

  saving.value = true
  try {
    await updateMemberProfile({
      nickname: profileForm.nickname,
      avatar: profileForm.avatar,
      gender: profileForm.gender,
      email: profileForm.email,
      mobile: profileForm.mobile,
    })
    ElMessage.success('保存成功')
    const info = await getMemberInfo()
    if (info) memberStore.setMemberInfo(info)
  } catch { /* 拦截器已处理 */ } finally {
    saving.value = false
  }
}

// ===== 修改密码表单 =====
const passwordFormRef = ref<FormInstance>()
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

const passwordRules: FormRules = {
  oldPassword: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码不少于 6 位', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (_r: any, v: string, cb: any) => {
        v !== passwordForm.newPassword ? cb(new Error('两次密码不一致')) : cb()
      },
      trigger: 'blur',
    },
  ],
}

const changingPassword = ref(false)

const handleChangePassword = async () => {
  if (!passwordFormRef.value) return
  const valid = await passwordFormRef.value.validate().catch(() => false)
  if (!valid) return

  changingPassword.value = true
  try {
    await changeMemberPassword({
      oldPassword: passwordForm.oldPassword,
      newPassword: passwordForm.newPassword,
    })
    ElMessage.success('密码修改成功')
    passwordFormRef.value.resetFields()
  } catch { /* 拦截器已处理 */ } finally {
    changingPassword.value = false
  }
}
</script>

<style lang="scss" scoped>
.text-clay-foreground { color: #32325d; }
.text-clay-muted { color: #8898aa; }
.text-clay-accent { color: #5a8dee; }
.text-clay-success { color: #71dd37; }
.bg-clay-accent { background-color: #5a8dee; }
.font-heading { font-family: 'Nunito', 'PingFang SC', sans-serif; }

.shadow-clay-card {
  box-shadow: 16px 16px 32px rgba(165, 175, 190, 0.3), -10px -10px 24px rgba(255, 255, 255, 0.9),
    inset 6px 6px 12px rgba(90, 141, 238, 0.03), inset -6px -6px 12px rgba(255, 255, 255, 1);
}
.shadow-clay-btn {
  box-shadow: 12px 12px 24px rgba(90, 141, 238, 0.3), -8px -8px 16px rgba(255, 255, 255, 0.4),
    inset 4px 4px 8px rgba(255, 255, 255, 0.4), inset -4px -4px 8px rgba(0, 0, 0, 0.05);
}
.shadow-clay-btn-hover {
  box-shadow: 16px 16px 32px rgba(90, 141, 238, 0.4), -10px -10px 20px rgba(255, 255, 255, 0.5),
    inset 4px 4px 8px rgba(255, 255, 255, 0.4), inset -4px -4px 8px rgba(0, 0, 0, 0.05);
}
.shadow-clay-pressed {
  box-shadow: inset 10px 10px 20px #e0e5ec, inset -10px -10px 20px #ffffff;
}

@keyframes breathe { 0%, 100% { transform: scale(1); } 50% { transform: scale(1.05); } }
.animate-breathe { animation: breathe 6s ease-in-out infinite; }

.animate-in {
  animation: slideIn 0.4s ease-out;
}
@keyframes slideIn {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}

:deep(.clay-input) {
  .el-input__wrapper {
    height: 48px;
    padding: 0 16px;
    border-radius: 16px;
    background: #f0f3f8;
    box-shadow: inset 10px 10px 20px #e0e5ec, inset -10px -10px 20px #ffffff;
    border: none;
    transition: all 0.3s;
    &.is-focus {
      background: #fff;
      box-shadow: 16px 16px 32px rgba(165, 175, 190, 0.3), -10px -10px 24px rgba(255, 255, 255, 0.9),
        inset 6px 6px 12px rgba(90, 141, 238, 0.03), inset -6px -6px 12px rgba(255, 255, 255, 1);
    }
  }
  .el-input__inner { font-weight: 500; color: #32325d; }
}
</style>
