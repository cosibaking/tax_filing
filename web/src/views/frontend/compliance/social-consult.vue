<!-- 社保咨询 P-SC -->
<template>
  <section class="space-y-6">
    <div class="bg-white/70 backdrop-blur-xl rounded-[48px] shadow-clay-card border border-[#d1d9e6]/40 p-8 md:p-10">
      <h2 class="font-heading font-black text-2xl text-clay-foreground mb-2">社保咨询</h2>
      <p class="text-sm text-clay-muted mb-6">
        进阶版及以上套餐可提交咨询，顾问将在工作时间 4 小时内首次回复。
      </p>

      <ElForm label-position="top" class="max-w-xl">
        <ElFormItem label="咨询类型">
          <ElSelect v-model="form.category" class="w-full">
            <ElOption label="创始人本人参保" value="founder" />
            <ElOption label="雇员社保" value="employee" />
            <ElOption label="其他" value="other" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="注册地/属地（选填）">
          <ElInput v-model="form.regionCode" placeholder="如：浙江省杭州市" />
        </ElFormItem>
        <ElFormItem label="您的问题" required>
          <ElInput
            v-model="form.question"
            type="textarea"
            :rows="4"
            placeholder="请描述您的社保参保疑问，至少 5 个字"
          />
        </ElFormItem>
        <ElButton type="primary" :loading="submitting" @click="handleSubmit">提交咨询</ElButton>
      </ElForm>
    </div>

    <div class="bg-white/70 backdrop-blur-xl rounded-[48px] shadow-clay-card border border-[#d1d9e6]/40 p-8 md:p-10">
      <h3 class="font-heading font-black text-lg text-clay-foreground mb-4">我的咨询记录</h3>
      <div v-if="listLoading" class="py-8 text-center">
        <ArtSvgIcon icon="ri:loader-4-line" class="text-2xl animate-spin mx-auto text-clay-accent" />
      </div>
      <div v-else-if="consults.length === 0" class="text-sm text-clay-muted py-4">暂无咨询记录</div>
      <div v-else class="space-y-4">
        <div
          v-for="item in consults"
          :key="item.id"
          class="p-5 rounded-2xl bg-[#f0f3f8] shadow-clay-pressed"
        >
          <div class="flex flex-wrap items-center gap-2 mb-2">
            <span class="text-xs font-bold px-2 py-0.5 rounded-full bg-blue-100 text-clay-accent">
              {{ SOCIAL_CONSULT_STATUS_LABELS[item.status] }}
            </span>
            <span class="text-xs text-clay-muted">{{ item.createdAt }}</span>
          </div>
          <p class="text-sm font-medium text-clay-foreground whitespace-pre-wrap">{{ item.question }}</p>
          <div v-if="item.reply" class="mt-3 pt-3 border-t border-[#d1d9e6]/40">
            <p class="text-xs font-black text-clay-muted mb-1">顾问回复</p>
            <p class="text-sm text-clay-foreground whitespace-pre-wrap">{{ item.reply }}</p>
            <p v-if="item.repliedAt" class="text-xs text-clay-muted mt-1">{{ item.repliedAt }}</p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import {
  createSocialConsult,
  getSocialConsults,
  SOCIAL_CONSULT_STATUS_LABELS,
  type SocialConsultItem
} from '@/api/frontend/compliance/social'
import { ElMessage } from 'element-plus'

defineOptions({ name: 'ComplianceSocialConsult' })

const form = reactive({
  category: 'founder',
  regionCode: '',
  question: ''
})
const submitting = ref(false)
const listLoading = ref(true)
const consults = ref<SocialConsultItem[]>([])

async function loadList() {
  listLoading.value = true
  try {
    const res = await getSocialConsults()
    consults.value = res.list || []
  } finally {
    listLoading.value = false
  }
}

async function handleSubmit() {
  if (form.question.trim().length < 5) {
    ElMessage.warning('请至少输入 5 个字的咨询内容')
    return
  }
  submitting.value = true
  try {
    await createSocialConsult({
      category: form.category,
      question: form.question.trim(),
      regionCode: form.regionCode.trim() || undefined
    })
    ElMessage.success('咨询已提交')
    form.question = ''
    await loadList()
  } catch (e: any) {
    ElMessage.error(e?.message || '提交失败')
  } finally {
    submitting.value = false
  }
}

onMounted(loadList)
</script>
