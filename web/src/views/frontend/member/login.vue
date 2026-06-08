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
  <main class="pt-20 pb-8 px-6 flex items-center justify-center min-h-[80vh]">
    <!-- 会员中心已禁用提示 -->
    <div v-if="!memberCenterOpen" class="w-full max-w-md">
      <div class="bg-white/70 backdrop-blur-2xl rounded-[48px] shadow-clay-deep border border-[#d1d9e6]/40 p-10 md:p-12 text-center">
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

    <!-- 登录卡片 -->
    <div v-else class="w-full max-w-md relative">
      <div class="bg-white/70 backdrop-blur-2xl rounded-[48px] shadow-clay-deep border border-[#d1d9e6]/40 p-10 md:p-12 relative z-10">
        <!-- Header -->
        <div class="text-center mb-10">
          <img v-if="siteStore.getLogo()" :src="siteStore.getLogo()" alt="logo" class="w-16 h-16 rounded-[20px] shadow-clay-btn mb-6 mx-auto object-cover animate-breathe" />
          <div v-else class="inline-flex w-16 h-16 rounded-[20px] bg-gradient-to-br from-blue-400 to-blue-600 shadow-clay-btn items-center justify-center text-white text-3xl font-black mb-6 animate-breathe">
            {{ siteName.charAt(0) }}
          </div>
          <h1 class="font-heading font-black text-3xl text-clay-foreground mb-2">欢迎回来</h1>
          <p class="text-clay-muted font-medium">请登录您的 {{ siteName }} 账户</p>
        </div>

        <!-- 表单 -->
        <ElForm ref="formRef" :model="formData" :rules="rules" @keyup.enter="handleSubmit" class="space-y-6">
          <ElFormItem prop="username" class="!mb-0">
            <label class="block text-sm font-bold text-clay-foreground mb-3 ml-1">账号</label>
            <ElInput
              v-model.trim="formData.username"
              placeholder="用户名 / 邮箱"
              size="large"
              class="clay-input"
            >
              <template #prefix>
                <ArtSvgIcon icon="ri:user-line" class="text-lg text-clay-muted" />
              </template>
            </ElInput>
          </ElFormItem>

          <ElFormItem prop="password" class="!mb-0">
            <label class="block text-sm font-bold text-clay-foreground mb-3 ml-1">密码</label>
            <ElInput
              v-model.trim="formData.password"
              placeholder="请输入密码"
              type="password"
              show-password
              size="large"
              class="clay-input"
            >
              <template #prefix>
                <ArtSvgIcon icon="ri:lock-line" class="text-lg text-clay-muted" />
              </template>
            </ElInput>
          </ElFormItem>

          <!-- 记住我 -->
          <div class="flex items-center justify-between px-1 pt-2">
            <ElCheckbox v-model="formData.rememberMe" class="!text-clay-muted">
              <span class="text-xs font-bold text-clay-muted">记住我</span>
            </ElCheckbox>
            <a href="javascript:;" class="text-xs font-bold text-clay-accent hover:underline">忘记密码？</a>
          </div>

          <!-- 提交按钮：点击后先弹出点选验证码 -->
          <button
            type="button"
            class="w-full h-14 rounded-2xl bg-gradient-to-br from-blue-400 to-blue-600 text-white font-black text-lg shadow-clay-btn hover:shadow-clay-btn-hover hover:-translate-y-1 active:scale-95 active:shadow-clay-pressed transition-all duration-300 mt-4 flex items-center justify-center gap-2"
            :disabled="loading"
            @click="handleSubmit"
          >
            <ArtSvgIcon v-if="loading" icon="ri:loader-4-line" class="text-xl animate-spin" />
            {{ loading ? '登录中...' : '立即登录' }}
          </button>
        </ElForm>

        <!-- Footer -->
        <div class="mt-10 text-center">
          <p class="text-sm text-clay-muted font-medium">
            还没有账号？
            <RouterLink to="/user/register" class="text-clay-accent font-black hover:underline ml-1">立即注册</RouterLink>
          </p>
        </div>
      </div>

      <!-- 装饰球 -->
      <div class="absolute -top-6 -right-6 w-20 h-20 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 shadow-clay-btn animate-breathe flex items-center justify-center text-white text-2xl font-bold z-20">
        Hi!
      </div>
      <div class="absolute -bottom-10 -left-10 w-32 h-32 rounded-full bg-gradient-to-br from-cyan-300 to-cyan-500 opacity-20 blur-2xl animate-float z-0"></div>
    </div>

    <!-- 点选验证码组件 -->
    <ArtClickCaptcha ref="captchaRef" :on-success="onCaptchaSuccess" />
  </main>
</template>

<script setup lang="ts">
import { useSiteStore } from '@/store/modules/site'
import { useMemberStore } from '@/store/modules/member'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { memberLogin, getMemberInfo } from '@/api/frontend'
import ArtClickCaptcha from '@/components/core/forms/art-click-captcha/index.vue'

defineOptions({ name: 'UserLogin' })

const router = useRouter()
const route = useRoute()

const siteStore = useSiteStore()
const memberStore = useMemberStore()

const siteName = computed(() => siteStore.getSiteName())
const memberCenterOpen = computed(() => siteStore.isUserCenterEnabled())

const formRef = ref<FormInstance>()
const captchaRef = ref<InstanceType<typeof ArtClickCaptcha>>()

const formData = reactive({
  username: siteStore.isDemoMode() ? 'user' : '',
  password: siteStore.isDemoMode() ? '123456' : '',
  rememberMe: false
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const loading = ref(false)

// 点击登录按钮 → 先验证表单 → 再弹出点选验证码
const handleSubmit = async () => {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  // 弹出点选验证码
  captchaRef.value?.open()
}

// 验证码收集完成 → 连同 captchaId/Info 一起发给登录接口校验
const onCaptchaSuccess = async (captchaId: string, captchaInfo: string) => {
  loading.value = true
  try {
    const { token } = await memberLogin({
      username: formData.username,
      password: formData.password,
      captchaId,
      captcha: captchaInfo, // 后端字段名是 captcha，值是点选坐标
    })

    if (!token) throw new Error('Login failed')

    memberStore.setToken(token)
    const info = await getMemberInfo()
    memberStore.setMemberInfo(info)

    ElMessage.success('登录成功')
    const redirect = route.query.redirect as string
    router.push(redirect || '/user/overview')
  } catch {
    // 错误已由拦截器处理
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
.text-clay-foreground { color: #32325d; }
.text-clay-muted { color: #8898aa; }
.text-clay-accent { color: #5a8dee; }
.font-heading { font-family: 'Nunito', 'PingFang SC', sans-serif; }

.shadow-clay-deep {
  box-shadow: 30px 30px 60px #d1d9e6, -30px -30px 60px #ffffff,
    inset 10px 10px 20px rgba(90, 141, 238, 0.05), inset -10px -10px 20px rgba(255, 255, 255, 0.8);
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
@keyframes float { 0%, 100% { transform: translateY(0); } 50% { transform: translateY(-20px); } }
.animate-float { animation: float 8s ease-in-out infinite; }

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
