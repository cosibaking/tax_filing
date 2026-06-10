<!-- 服务 SLA 承诺面板 -->
<template>
  <div class="sla-panel">
    <h3 v-if="title" class="sla-panel__title">{{ title }}</h3>
    <p v-if="subtitle" class="sla-panel__subtitle">{{ subtitle }}</p>
    <ul class="sla-list">
      <li v-for="item in items" :key="item.phase" class="sla-item">
        <ArtSvgIcon :icon="item.icon" class="sla-item__icon" />
        <div class="sla-item__body">
          <span class="sla-item__phase">{{ item.phase }}</span>
          <span class="sla-item__deadline">{{ item.deadline }}</span>
          <span v-if="item.note" class="sla-item__note">{{ item.note }}</span>
        </div>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import type { SlaItem } from '@/data/frontend/complianceSla'

defineProps<{
  items: SlaItem[]
  title?: string
  subtitle?: string
}>()
</script>

<style lang="scss" scoped>
.sla-panel__title {
  margin: 0 0 6px;
  font-size: 16px;
  font-weight: 700;
  color: #1a1f36;
}

.sla-panel__subtitle {
  margin: 0 0 16px;
  font-size: 13px;
  color: #94a3b8;
}

.sla-list {
  margin: 0;
  padding: 0;
  list-style: none;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.sla-item {
  display: flex;
  gap: 10px;
  padding: 12px;
  border: 1px solid #e8edf3;
  border-radius: 8px;
  background: #f8fafc;
}

.sla-item__icon {
  flex-shrink: 0;
  font-size: 18px;
  color: #2563eb;
  margin-top: 2px;
}

.sla-item__body {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.sla-item__phase {
  font-size: 13px;
  font-weight: 700;
  color: #334155;
}

.sla-item__deadline {
  font-size: 12px;
  font-weight: 600;
  color: #2563eb;
}

.sla-item__note {
  font-size: 11px;
  color: #94a3b8;
}

@media (max-width: 768px) {
  .sla-list {
    grid-template-columns: 1fr;
  }
}
</style>
