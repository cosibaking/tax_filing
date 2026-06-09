<!-- 登录页面 v1.2.7 -->
<template>
  <div class="flex w-full h-screen">
    <LoginLeftView />

    <div class="relative flex-1">
      <AuthTopBar />

      <div class="auth-right-wrap">
        <div class="form">
          <h3 class="title">{{ $t('login.title') }}</h3>
          <p class="sub-title">{{ $t('login.subTitle') }}</p>
          <ElForm
            ref="formRef"
            :model="formData"
            :rules="rules"
            :key="formKey"
            @keyup.enter="handleSubmit"
            style="margin-top: 25px"
          >
            <ElFormItem prop="username">
              <ElInput
                class="custom-height"
                :placeholder="$t('login.placeholder.username')"
                v-model.trim="formData.username"
                @focus="onInputFocus"
                @blur="onInputBlur"
              />
            </ElFormItem>
            <ElFormItem prop="password">
              <ElInput
                class="custom-height"
                :placeholder="$t('login.placeholder.password')"
                v-model.trim="formData.password"
                :type="passwordVisible ? 'text' : 'password'"
                autocomplete="off"
                @focus="onInputFocus"
                @blur="onInputBlur"
              >
                <template #suffix>
                  <ElIcon
                    class="cursor-pointer"
                    @click="passwordVisible = !passwordVisible"
                  >
                    <svg v-if="passwordVisible" viewBox="0 0 1024 1024" width="1em" height="1em" fill="currentColor">
                      <path d="M512 256c-211.2 0-390.4 121.6-480 307.2 89.6 185.6 268.8 307.2 480 307.2s390.4-121.6 480-307.2C902.4 377.6 723.2 256 512 256zm0 512c-112 0-204.8-92.8-204.8-204.8S400 358.4 512 358.4s204.8 92.8 204.8 204.8S624 768 512 768zm0-332.8c-70.4 0-128 57.6-128 128s57.6 128 128 128 128-57.6 128-128-57.6-128-128-128z" />
                    </svg>
                    <svg v-else viewBox="0 0 1024 1024" width="1em" height="1em" fill="currentColor">
                      <path d="M942.4 521.6c-25.6-44.8-57.6-86.4-92.8-121.6l-64 64c28.8 28.8 54.4 60.8 76.8 99.2-89.6 160-249.6 256-400 256-51.2 0-99.2-9.6-144-28.8l-70.4 70.4c64 32 134.4 51.2 214.4 51.2 211.2 0 390.4-121.6 480-307.2-9.6-28.8-25.6-57.6-44.8-83.2zM876.8 128l-48-48-188.8 188.8c-38.4-16-80-25.6-128-25.6-211.2 0-390.4 121.6-480 307.2 51.2 102.4 128 188.8 220.8 243.2L76.8 972.8l48 48 752-752zM320 563.2c0-105.6 86.4-192 192-192 38.4 0 73.6 12.8 105.6 32L352 668.8c-19.2-28.8-32-67.2-32-105.6z" />
                    </svg>
                  </ElIcon>
                </template>
              </ElInput>
            </ElFormItem>

            <!-- 点选验证码 -->
            <ArtClickCaptcha ref="clickCaptchaRef" :on-success="onCaptchaSuccess" />

            <div class="flex-cb mt-2 text-sm">
              <ElCheckbox v-model="formData.rememberPassword">{{
                $t('login.rememberPwd')
              }}</ElCheckbox>
              <RouterLink class="text-theme" :to="{ name: 'ForgetPassword' }">{{
                $t('login.forgetPwd')
              }}</RouterLink>
            </div>

            <div style="margin-top: 30px">
              <ElButton
                class="w-full custom-height"
                type="primary"
                @click="handleSubmit"
                :loading="loading"
                v-ripple
              >
                {{ $t('login.btnText') }}
              </ElButton>
            </div>

            <div class="mt-5 text-sm text-gray-600">
              <span>{{ $t('login.noAccount') }}</span>
              <RouterLink class="text-theme" :to="{ name: 'Register' }">{{
                $t('login.register')
              }}</RouterLink>
            </div>
          </ElForm>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { useUserStore } from '@/store/modules/user'
  import { useSiteStore } from '@/store/modules/site'
  import { useI18n } from 'vue-i18n'
  import { HttpError } from '@/utils/http/error'
  import { fetchLogin } from '@/api/backend/auth'
  import { ADMIN_BASE_PATH } from '@/router/routesAlias'
  import { ElNotification, type FormInstance, type FormRules } from 'element-plus'
  import ArtClickCaptcha from '@/components/core/forms/art-click-captcha/index.vue'

  defineOptions({ name: 'Login' })

  const { t, locale } = useI18n()
  const formKey = ref(0)
  watch(locale, () => { formKey.value++ })

  const clickCaptchaRef = ref<InstanceType<typeof ArtClickCaptcha>>()
  const userStore = useUserStore()
  const router = useRouter()
  const route = useRoute()
  const formRef = ref<FormInstance>()

  const siteStore = useSiteStore()
  const systemName = computed(() => siteStore.getSiteName())

  const isInputFocused = ref(false)
  const passwordVisible = ref(false)

  let blurTimer: ReturnType<typeof setTimeout> | null = null
  const onInputFocus = () => {
    if (blurTimer) { clearTimeout(blurTimer); blurTimer = null }
    isInputFocused.value = true
  }
  const onInputBlur = () => {
    blurTimer = setTimeout(() => { isInputFocused.value = false }, 100)
  }

  const formData = reactive({
    username: siteStore.isDemoMode() ? 'admin' : '',
    password: siteStore.isDemoMode() ? '123456' : '',
    rememberPassword: true
  })

  provide('loginAnimation', {
    isTyping: isInputFocused,
    hasSecret: computed(() => formData.password.length > 0),
    secretVisible: passwordVisible,
  })

  const rules = computed<FormRules>(() => ({
    username: [{ required: true, message: t('login.placeholder.username'), trigger: 'blur' }],
    password: [{ required: true, message: t('login.placeholder.password'), trigger: 'blur' }]
  }))

  const loading = ref(false)

  const handleSubmit = async () => {
    if (!formRef.value) return
    const valid = await formRef.value.validate().catch(() => false)
    if (!valid) return
    clickCaptchaRef.value?.open()
  }

  const onCaptchaSuccess = async (captchaId: string, captchaInfo: string) => {
    loading.value = true
    try {
      const { username, password } = formData
      const { accessToken, token, refreshToken } = await fetchLogin({
        userName: username,
        password,
        captchaId,
        captchaInfo,
      })

      const finalToken = accessToken || token
      if (!finalToken) throw new Error('Login failed - no token received')

      userStore.setToken(finalToken, refreshToken)
      userStore.setLoginStatus(true)

      const redirect = route.query.redirect as string
      let targetPath = `${ADMIN_BASE_PATH}/dashboard/console`
      if (redirect && !redirect.startsWith('/user') && redirect !== '/') {
        targetPath = redirect
      }

      router.push(targetPath)

      const unregister = router.afterEach(() => {
        unregister()
        loading.value = false
        showLoginSuccessNotice()
      })
    } catch (error) {
      loading.value = false
      if (!(error instanceof HttpError)) {
        console.error('[Login] Unexpected error:', error)
      }
    }
  }

  const showLoginSuccessNotice = () => {
    setTimeout(() => {
      const userInfo = userStore.getUserInfo
      const displayName = userInfo.nickname || userInfo.username || systemName.value
      ElNotification({
        title: t('login.success.title'),
        type: 'success',
        duration: 2500,
        zIndex: 10000,
        message: `${t('login.success.message')}, ${displayName}!`
      })
    }, 1000)
  }
</script>

<style scoped>
  @import './style.css';
</style>
