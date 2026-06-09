/**
 * 访客诊断本地缓存：未登录时暂存问卷与结果，登录后同步至服务端
 */
import type { DiagnosisSubmitParams, DiagnosisSubmitResult } from '@/api/frontend/compliance/diagnosis'

const CACHE_KEY = 'compliance_guest_diagnoses'
const SESSION_RESULT_PREFIX = 'compliance_diagnosis_result_'
/** 最多保留条数 */
const MAX_ENTRIES = 3
/** 过期时间：72 小时 */
const TTL_MS = 72 * 60 * 60 * 1000

export interface GuestDiagnosisRecord {
  guestId: string
  payload: DiagnosisSubmitParams
  result: DiagnosisSubmitResult
  cachedAt: number
}

export function isGuestDiagnosisId(id: number | string | undefined | null): boolean {
  return typeof id === 'string' && id.startsWith('guest-')
}

export function generateGuestDiagnosisId(): string {
  return `guest-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
}

function readAll(): GuestDiagnosisRecord[] {
  try {
    const raw = localStorage.getItem(CACHE_KEY)
    if (!raw) return []
    const parsed = JSON.parse(raw) as GuestDiagnosisRecord[]
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

function writeAll(records: GuestDiagnosisRecord[]) {
  localStorage.setItem(CACHE_KEY, JSON.stringify(records))
}

/** 清理过期与超量记录，并同步清理 sessionStorage 中的结果缓存 */
export function pruneGuestDiagnosisCache(): GuestDiagnosisRecord[] {
  const now = Date.now()
  let records = readAll().filter((r) => now - r.cachedAt <= TTL_MS)
  records.sort((a, b) => b.cachedAt - a.cachedAt)

  if (records.length > MAX_ENTRIES) {
    const removed = records.slice(MAX_ENTRIES)
    records = records.slice(0, MAX_ENTRIES)
    removed.forEach((r) => sessionStorage.removeItem(`${SESSION_RESULT_PREFIX}${r.guestId}`))
  }

  writeAll(records)
  return records
}

/** 保存访客诊断（含问卷与结果） */
export function saveGuestDiagnosis(
  guestId: string,
  payload: DiagnosisSubmitParams,
  result: DiagnosisSubmitResult
) {
  const records = pruneGuestDiagnosisCache().filter((r) => r.guestId !== guestId)
  records.unshift({
    guestId,
    payload,
    result: { ...result, id: guestId },
    cachedAt: Date.now()
  })
  if (records.length > MAX_ENTRIES) {
    const removed = records.splice(MAX_ENTRIES)
    removed.forEach((r) => sessionStorage.removeItem(`${SESSION_RESULT_PREFIX}${r.guestId}`))
  }
  writeAll(records)
}

/** 读取单条访客诊断 */
export function getGuestDiagnosis(guestId: string): GuestDiagnosisRecord | null {
  const record = pruneGuestDiagnosisCache().find((r) => r.guestId === guestId)
  return record ?? null
}

/** 列出有效访客诊断（按时间倒序） */
export function listGuestDiagnoses(): GuestDiagnosisRecord[] {
  return pruneGuestDiagnosisCache()
}

/** 清空访客诊断缓存 */
export function clearGuestDiagnosisCache() {
  const records = readAll()
  records.forEach((r) => sessionStorage.removeItem(`${SESSION_RESULT_PREFIX}${r.guestId}`))
  localStorage.removeItem(CACHE_KEY)
}

/** 收集 sessionStorage 中可能存在的匿名 DB 诊断 ID（兼容旧数据） */
export function collectAnonymousDiagnosisIds(): number[] {
  const ids = new Set<number>()
  for (let i = 0; i < sessionStorage.length; i++) {
    const key = sessionStorage.key(i)
    if (!key?.startsWith(SESSION_RESULT_PREFIX)) continue
    const suffix = key.slice(SESSION_RESULT_PREFIX.length)
    if (isGuestDiagnosisId(suffix)) continue
    const numId = Number(suffix)
    if (Number.isInteger(numId) && numId > 0) ids.add(numId)
  }
  return [...ids]
}
