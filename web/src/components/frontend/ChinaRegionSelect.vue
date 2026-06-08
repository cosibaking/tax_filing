<!-- 中国省市区三级联动（支持模糊搜索，含港澳台） -->
<template>
  <div class="grid md:grid-cols-3 gap-4">
    <ElFormItem label="省" :prop="provinceProp">
      <ElSelect
        v-model="province"
        filterable
        clearable
        placeholder="请选择省份"
        size="large"
        class="w-full clay-select"
        :disabled="disabled"
        @change="handleProvinceChange"
      >
        <ElOption
          v-for="item in provinces"
          :key="item.code"
          :label="item.name"
          :value="item.name"
        />
      </ElSelect>
    </ElFormItem>
    <ElFormItem label="市" :prop="cityProp">
      <ElSelect
        v-model="city"
        filterable
        clearable
        placeholder="请选择城市"
        size="large"
        class="w-full clay-select"
        :disabled="disabled || !province"
        @change="handleCityChange"
      >
        <ElOption
          v-for="item in cities"
          :key="item.code"
          :label="item.name"
          :value="item.name"
        />
      </ElSelect>
    </ElFormItem>
    <ElFormItem label="区" :prop="districtProp">
      <ElSelect
        v-model="district"
        filterable
        clearable
        placeholder="请选择区县"
        size="large"
        class="w-full clay-select"
        :disabled="disabled || !city"
      >
        <ElOption
          v-for="item in districts"
          :key="item.code"
          :label="item.name"
          :value="item.name"
        />
      </ElSelect>
    </ElFormItem>
  </div>
</template>

<script setup lang="ts">
import { CHINA_REGION_TREE, findProvince, type RegionNode } from '@/data/chinaRegions'

defineOptions({ name: 'ChinaRegionSelect' })

withDefaults(
  defineProps<{
    provinceProp?: string
    cityProp?: string
    districtProp?: string
    disabled?: boolean
  }>(),
  {
    provinceProp: 'registerProvince',
    cityProp: 'registerCity',
    districtProp: 'registerDistrict',
    disabled: false,
  }
)

const province = defineModel<string>('province', { default: '' })
const city = defineModel<string>('city', { default: '' })
const district = defineModel<string>('district', { default: '' })

const provinces = CHINA_REGION_TREE

const cities = computed<RegionNode[]>(() => {
  if (!province.value) return []
  return findProvince(province.value)?.children ?? []
})

const districts = computed<RegionNode[]>(() => {
  if (!province.value || !city.value) return []
  const cityNode = cities.value.find((c) => c.name === city.value)
  return cityNode?.children ?? []
})

function handleProvinceChange() {
  city.value = ''
  district.value = ''
}

function handleCityChange() {
  district.value = ''
}
</script>
