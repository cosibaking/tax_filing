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
  <Teleport to="body">
    <!-- 遮罩层 -->
    <Transition name="fade">
      <div v-if="visible" class="click-captcha-overlay" @click="handleClose"></div>
    </Transition>
    <!-- 验证码弹窗 -->
    <Transition name="pop">
      <div v-if="visible" class="click-captcha-dialog">
        <div class="click-captcha-card">
          <!-- Loading -->
          <div v-if="loading" class="captcha-loading" :style="{ width: `${captchaData.width}px`, height: `${captchaData.height}px` }">
            <ArtSvgIcon icon="ri:loader-4-line" class="text-2xl animate-spin text-clay-accent" />
            <span class="ml-2 text-sm text-clay-muted font-bold">加载中...</span>
          </div>
          <!-- 图片区 -->
          <div v-else class="captcha-img-box">
            <img
              ref="captchaImgRef"
              class="captcha-img"
              :src="captchaData.base64"
              :style="{ width: `${captchaData.width}px`, height: `${captchaData.height}px` }"
              alt="验证码"
              @click.prevent="onImageClick"
            />
            <!-- 点击标记 -->
            <span
              v-for="(pos, index) in clickPoints"
              :key="index"
              class="captcha-step"
              :style="{ left: `${pos.x - 13}px`, top: `${pos.y - 13}px` }"
              @click.stop="onCancelClick(index)"
            >
              {{ index + 1 }}
            </span>
          </div>
          <!-- 提示文字 -->
          <div class="captcha-prompt">
            <template v-if="tipMsg">
              <span :class="tipSuccess ? 'text-emerald-500' : 'text-red-500'" class="font-bold">{{ tipMsg }}</span>
            </template>
            <template v-else>
              <span class="text-clay-muted">请依次点击</span>
              <span
                v-for="(text, index) in captchaData.text"
                :key="index"
                class="captcha-text-hint"
                :class="clickPoints.length > index ? 'is-clicked' : ''"
              >
                {{ text }}
              </span>
            </template>
          </div>
          <!-- 刷新按钮 -->
          <div class="captcha-refresh">
            <div class="captcha-refresh-line"></div>
            <div class="captcha-refresh-btn" @click="loadCaptcha" title="刷新验证码">
              <ArtSvgIcon icon="ri:refresh-line" class="text-xl text-clay-muted hover:text-clay-accent transition-colors" />
            </div>
            <div class="captcha-refresh-line"></div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { getClickCaptcha } from '@/api/common/captcha'
import { HttpError } from '@/utils/http/error'

defineOptions({ name: 'ArtClickCaptcha' })

interface CaptchaState {
  id: string
  text: string[]
  base64: string
  width: number
  height: number
}

const props = withDefaults(defineProps<{
  /** 验证成功回调，返回 captchaId + captchaInfo 供业务接口使用 */
  onSuccess?: (captchaId: string, captchaInfo: string) => void
}>(), {})

const visible = ref(false)
const loading = ref(true)
const tipMsg = ref('')
const tipSuccess = ref(false)
const captchaImgRef = ref<HTMLImageElement>()

const captchaData = reactive<CaptchaState>({
  id: '',
  text: [],
  base64: '',
  width: 350,
  height: 200,
})

const clickPoints = ref<{ x: number; y: number }[]>([])

// 加载验证码
const loadCaptcha = async () => {
  loading.value = true
  tipMsg.value = ''
  clickPoints.value = []
  try {
    const data = await getClickCaptcha()
    captchaData.id = data.id
    captchaData.text = data.text
    captchaData.base64 = data.base64
    captchaData.width = data.width
    captchaData.height = data.height
  } catch (error) {
    tipMsg.value = error instanceof HttpError ? error.message : '验证码加载失败'
    tipSuccess.value = false
  } finally {
    loading.value = false
  }
}

// 图片点击 — 只收集坐标，不调校验API（验证码与登录接口强关联）
const onImageClick = (event: MouseEvent) => {
  if (clickPoints.value.length >= captchaData.text.length) return

  const x = event.offsetX
  const y = event.offsetY
  clickPoints.value.push({ x, y })

  // 达到要求的点击数量后，拼接 captchaInfo 回调给父组件
  if (clickPoints.value.length === captchaData.text.length) {
    const img = captchaImgRef.value
    if (!img) return

    const captchaInfo = [
      clickPoints.value.map(p => `${p.x},${p.y}`).join('-'),
      img.clientWidth,
      img.clientHeight,
    ].join(';')

    // 直接回调（由业务登录接口统一校验）
    tipMsg.value = '提交中...'
    tipSuccess.value = true
    setTimeout(() => {
      props.onSuccess?.(captchaData.id, captchaInfo)
      handleClose()
    }, 300)
  }
}

// 取消某个点击
const onCancelClick = (index: number) => {
  clickPoints.value.splice(index, 1)
}

// 关闭
const handleClose = () => {
  visible.value = false
}

// 对外暴露打开方法
const open = () => {
  visible.value = true
  loadCaptcha()
}

defineExpose({ open, close: handleClose })
</script>

<style lang="scss" scoped>
.click-captcha-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  z-index: 9998;
}

.click-captcha-dialog {
  position: fixed;
  z-index: 9999;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
}

.click-captcha-card {
  padding: 12px;
  border: 1px solid var(--el-border-color-extra-light);
  background-color: var(--el-color-white);
  border-radius: 12px;
  box-shadow: 0 0 0 1px hsla(0, 0%, 100%, 0.3) inset, 0 0.5em 1em rgba(0, 0, 0, 0.6);
}

.captcha-loading {
  display: flex;
  align-items: center;
  justify-content: center;
}

.captcha-img-box {
  position: relative;

  .captcha-img {
    display: block;
    border: none;
    cursor: pointer;
    border-radius: 6px;
  }

  .captcha-step {
    box-sizing: border-box;
    position: absolute;
    width: 26px;
    height: 26px;
    line-height: 26px;
    font-size: 12px;
    font-weight: bold;
    text-align: center;
    color: #fff;
    border: 2px solid rgba(255, 255, 255, 0.8);
    background: var(--el-color-primary);
    border-radius: 50%;
    box-shadow: 0 0 10px rgba(255, 255, 255, 0.5);
    user-select: none;
    cursor: pointer;
    transition: transform 0.2s;

    &:hover {
      transform: scale(1.15);
    }
  }
}

.captcha-prompt {
  height: 40px;
  line-height: 40px;
  font-size: 14px;
  text-align: center;

  .captcha-text-hint {
    margin-left: 8px;
    font-size: 16px;
    font-weight: bold;
    color: var(--el-color-danger);

    &.is-clicked {
      color: var(--el-color-primary);
    }
  }
}

.captcha-refresh {
  display: flex;
  align-items: center;
  margin-top: 8px;

  .captcha-refresh-line {
    flex: 1;
    height: 1px;
    background: #e0e0e0;
  }

  .captcha-refresh-btn {
    cursor: pointer;
    padding: 4px 12px;
    display: flex;
    align-items: center;
  }
}

.text-clay-accent { color: #5a8dee; }
.text-clay-muted { color: #8898aa; }

/* 过渡动画 */
.fade-enter-active, .fade-leave-active { transition: opacity 0.25s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

.pop-enter-active { transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1); }
.pop-leave-active { transition: all 0.2s ease; }
.pop-enter-from { opacity: 0; transform: translate(-50%, -50%) scale(0.8); }
.pop-leave-to { opacity: 0; transform: translate(-50%, -50%) scale(0.9); }
</style>
