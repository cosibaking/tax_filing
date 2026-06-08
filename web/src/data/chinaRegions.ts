import levelData from 'province-city-china/dist/level.json'
import taiwanRegions from './taiwanRegions.json'

export interface RegionNode {
  code: string
  name: string
  children?: RegionNode[]
}

/** 直辖市：下级直接为区县 */
const MUNICIPALITY_CODES = new Set(['110000', '120000', '310000', '500000'])
/** 港澳：下级直接为区县，无地级 */
const SAR_CODES = new Set(['810000', '820000'])

function patchTaiwan(tree: RegionNode[]): RegionNode[] {
  return tree.map((node) => {
    if (node.code !== '710000') return node
    return { ...node, children: taiwanRegions as RegionNode[] }
  })
}

function normalizeTree(raw: RegionNode[]): RegionNode[] {
  return patchTaiwan(raw).map((province) => {
    if (MUNICIPALITY_CODES.has(province.code) || SAR_CODES.has(province.code)) {
      const cityName = SAR_CODES.has(province.code)
        ? province.name
        : province.name.replace(/市$/, '') + '市'
      return {
        ...province,
        children: [
          {
            code: `${province.code.slice(0, 2)}0100`,
            name: cityName,
            children: province.children ?? [],
          },
        ],
      }
    }
    return province
  })
}

export const CHINA_REGION_TREE: RegionNode[] = normalizeTree(levelData as RegionNode[])

export function findProvince(name: string): RegionNode | undefined {
  return CHINA_REGION_TREE.find((p) => p.name === name)
}

export function findCity(provinceName: string, cityName: string): RegionNode | undefined {
  return findProvince(provinceName)?.children?.find((c) => c.name === cityName)
}

export function isDirectRegion(provinceCode: string): boolean {
  return MUNICIPALITY_CODES.has(provinceCode) || SAR_CODES.has(provinceCode)
}
