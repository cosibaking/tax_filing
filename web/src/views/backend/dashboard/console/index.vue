<!-- 工作台页面 -->
<template>
  <div class="compliance-console">
    <ComplianceStats />
    <PendingWork />
  </div>
</template>

<script setup lang="ts">
import ComplianceStats from './modules/compliance-stats.vue'
import PendingWork from './modules/pending-work.vue'

defineOptions({ name: 'Console' })
</script>
