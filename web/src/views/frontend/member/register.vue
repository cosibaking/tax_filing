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
      <div
        class="bg-white rounded-xl shadow-clay-deep border border-[#e8edf3] p-8 md:p-10 text-center"
      >
        <div
          class="w-16 h-16 rounded-lg bg-[#f8fafc] border border-[#e8edf3] flex items-center justify-center mx-auto mb-6"
        >
          <ArtSvgIcon icon="ri:lock-2-line" class="text-[36px] text-clay-muted" />
        </div>
        <h2 class="font-heading font-black text-2xl text-clay-foreground mb-3">会员中心已关闭</h2>
        <p class="text-clay-muted font-medium leading-relaxed"
          >会员中心已禁用，请联系网站管理员开启。</p
        >
        <RouterLink
          to="/"
          class="inline-block mt-8 px-8 py-3 rounded-lg bg-white border border-[#d8dee9] hover:bg-[#f8fafc] font-bold text-clay-foreground transition-all"
        >
          返回首页
        </RouterLink>
      </div>
    </div>

    <div v-else class="w-full max-w-md relative">
      <div
        class="bg-white rounded-xl shadow-clay-deep border border-[#e8edf3] p-8 md:p-10 relative z-10"
      >
        <!-- Header -->
        <div class="text-center mb-8">
          <img
            v-if="siteStore.getLogo()"
            :src="siteStore.getLogo()"
            alt="logo"
            class="w-14 h-14 rounded-lg border border-[#e8edf3] mb-6 mx-auto object-contain bg-[#2563eb]"
          />
          <div
            v-else
            class="inline-flex w-14 h-14 rounded-lg bg-[#2563eb] items-center justify-center text-white text-2xl font-black mb-6"
          >
            {{ siteName.charAt(0) }}
          </div>
          <h1 class="font-heading font-black text-3xl text-clay-foreground mb-2">创建账号</h1>
          <p class="text-clay-muted font-medium">注册 {{ siteName }} 会员</p>
        </div>

        <!-- 表单 -->
        <ElForm
          ref="formRef"
          :model="formData"
          :rules="rules"
          class="register-form"
          @keyup.enter="handleSubmit"
        >
          <ElFormItem prop="username" class="register-form-item">
            <label class="block text-sm font-bold text-clay-foreground mb-3 ml-1">用户名</label>
            <ElInput
              v-model.trim="formData.username"
              placeholder="4-20位字母数字"
              size="large"
              class="clay-input"
            >
              <template #prefix>
                <ArtSvgIcon icon="ri:user-line" class="text-lg text-clay-muted" />
              </template>
            </ElInput>
          </ElFormItem>

          <ElFormItem prop="password" class="register-form-item">
            <label class="block text-sm font-bold text-clay-foreground mb-3 ml-1">密码</label>
            <ElInput
              v-model.trim="formData.password"
              placeholder="不少于6位"
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

          <ElFormItem prop="confirmPassword" class="register-form-item">
            <label class="block text-sm font-bold text-clay-foreground mb-3 ml-1">确认密码</label>
            <ElInput
              v-model.trim="formData.confirmPassword"
              placeholder="再次输入密码"
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

          <ElFormItem prop="mobile" class="register-form-item">
            <label class="block text-sm font-bold text-clay-foreground mb-3 ml-1">手机号</label>
            <ElInput
              v-model.trim="formData.mobile"
              placeholder="请输入手机号"
              size="large"
              maxlength="11"
              class="clay-input"
            >
              <template #prefix>
                <ArtSvgIcon icon="ri:phone-line" class="text-lg text-clay-muted" />
              </template>
            </ElInput>
          </ElFormItem>

          <!-- 协议 -->
          <div class="register-agreements">
            <ElFormItem prop="agreeTerms" class="register-agreement-item">
              <ElCheckbox v-model="formData.agreeTerms" class="!text-clay-muted">
                <span class="text-xs font-bold text-clay-muted">
                  我已阅读并同意
                  <RouterLink
                    to="/legal/terms"
                    target="_blank"
                    class="text-clay-accent hover:underline"
                    >用户协议</RouterLink
                  >
                </span>
              </ElCheckbox>
            </ElFormItem>
            <ElFormItem prop="agreePrivacy" class="register-agreement-item">
              <ElCheckbox v-model="formData.agreePrivacy" class="!text-clay-muted">
                <span class="text-xs font-bold text-clay-muted">
                  我已阅读并同意
                  <RouterLink
                    to="/legal/privacy"
                    target="_blank"
                    class="text-clay-accent hover:underline"
                    >隐私政策</RouterLink
                  >
                </span>
              </ElCheckbox>
            </ElFormItem>
          </div>

          <!-- 提交按钮 -->
          <button
            type="button"
            class="register-submit w-full h-12 rounded-lg bg-[#2563eb] text-white font-bold text-base shadow-clay-btn hover:bg-[#1d4ed8] hover:-translate-y-0.5 active:translate-y-0 transition-all duration-200 flex items-center justify-center gap-2"
            :disabled="loading"
            @click="handleSubmit"
          >
            <ArtSvgIcon v-if="loading" icon="ri:loader-4-line" class="text-xl animate-spin" />
            {{ loading ? '注册中...' : '立即注册' }}
          </button>
        </ElForm>

        <!-- Footer -->
        <div class="mt-10 text-center">
          <p class="text-sm text-clay-muted font-medium">
            已有账号？
            <RouterLink to="/user/login" class="text-clay-accent font-black hover:underline ml-1"
              >返回登录</RouterLink
            >
          </p>
        </div>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
  import { useSiteStore } from '@/store/modules/site'
  import { useMemberStore } from '@/store/modules/member'
  import { useMemberMenuStore } from '@/store/modules/memberMenu'
  import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
  import { memberRegister, getMemberInfo } from '@/api/frontend'
  import { syncGuestDiagnosisAfterLogin } from '@/utils/compliance/syncGuestDiagnosis'

  defineOptions({ name: 'UserRegister' })

  const router = useRouter()
  const route = useRoute()
  const siteStore = useSiteStore()
  const memberStore = useMemberStore()
  const memberMenuStore = useMemberMenuStore()
  const siteName = computed(() => siteStore.getSiteName())
  const memberCenterOpen = computed(() => siteStore.isUserCenterEnabled())

  const formRef = ref<FormInstance>()
  const formData = reactive({
    username: '',
    password: '',
    confirmPassword: '',
    mobile: '',
    agreeTerms: false,
    agreePrivacy: false
  })

  const rules: FormRules = {
    username: [
      { required: true, message: '请输入用户名', trigger: 'blur' },
      { min: 4, max: 20, message: '用户名长度 4-20 位', trigger: 'blur' }
    ],
    password: [
      { required: true, message: '请输入密码', trigger: 'blur' },
      { min: 6, message: '密码不少于 6 位', trigger: 'blur' }
    ],
    confirmPassword: [
      { required: true, message: '请再次输入密码', trigger: 'blur' },
      {
        validator: (_rule: any, value: string, callback: any) => {
          if (value !== formData.password) {
            callback(new Error('两次密码不一致'))
          } else {
            callback()
          }
        },
        trigger: 'blur'
      }
    ],
    mobile: [
      { required: true, message: '请输入手机号', trigger: 'blur' },
      { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }
    ],
    agreeTerms: [
      {
        validator: (_rule: any, value: boolean, callback: any) => {
          if (value) {
            callback()
            return
          }
          callback(new Error('请先同意用户协议'))
        },
        trigger: 'change'
      }
    ],
    agreePrivacy: [
      {
        validator: (_rule: any, value: boolean, callback: any) => {
          if (value) {
            callback()
            return
          }
          callback(new Error('请先同意隐私政策'))
        },
        trigger: 'change'
      }
    ]
  }

  const loading = ref(false)

  const handleSubmit = async () => {
    if (!formRef.value) return
    const valid = await formRef.value.validate().catch(() => false)
    if (!valid) return

    loading.value = true
    try {
      const res = await memberRegister({
        username: formData.username,
        password: formData.password,
        mobile: formData.mobile,
        agreeTerms: formData.agreeTerms,
        agreePrivacy: formData.agreePrivacy
      })

      if (!res.token) {
        throw new Error('注册成功但未获取登录凭证')
      }

      memberStore.setToken(res.token, res.refreshToken)
      const info = await getMemberInfo()
      memberStore.setMemberInfo(info)
      await memberMenuStore.fetchMenus()

      await syncGuestDiagnosisAfterLogin()

      ElMessage.success('注册成功，已自动登录')
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
  .text-clay-foreground {
    color: #32325d;
  }

  .text-clay-muted {
    color: #8898aa;
  }

  .text-clay-accent {
    color: #5a8dee;
  }

  .font-heading {
    font-family: Nunito, 'PingFang SC', sans-serif;
  }

  .shadow-clay-deep {
    box-shadow: 0 12px 32px rgb(15 23 42 / 8%);
  }

  .shadow-clay-btn {
    box-shadow: 0 4px 14px rgb(37 99 235 / 28%);
  }

  .shadow-clay-btn-hover {
    box-shadow: 0 6px 18px rgb(37 99 235 / 32%);
  }

  .shadow-clay-pressed {
    box-shadow: inset 0 2px 4px rgb(15 23 42 / 8%);
  }

  .register-form {
    display: flex;
    flex-direction: column;
    gap: 24px;
  }

  .register-form-item,
  .register-agreement-item {
    margin-bottom: 0;
  }

  .register-agreements {
    display: flex;
    flex-direction: column;
    gap: 14px;
    padding: 2px 4px 0;
  }

  .register-submit {
    margin-top: 2px;
  }

  :deep(.register-form-item .el-form-item__content),
  :deep(.register-agreement-item .el-form-item__content) {
    display: flex;
    flex-direction: column;
    align-items: stretch;
  }

  :deep(.register-form .el-form-item__error) {
    position: static;
    width: 100%;
    padding-top: 6px;
    line-height: 1.45;
  }

  :deep(.clay-input) {
    .el-input__wrapper {
      height: 48px;
      padding: 0 16px;
      background: #fff;
      border: 1px solid #d8dee9;
      border-radius: 8px;
      box-shadow: none;
      transition: all 0.2s ease;

      &.is-focus {
        border-color: #2563eb;
        box-shadow: 0 0 0 3px rgb(37 99 235 / 12%);
      }
    }

    .el-input__inner {
      font-weight: 500;
      color: #32325d;
    }
  }
</style>
