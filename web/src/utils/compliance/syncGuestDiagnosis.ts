/**
 * 登录/注册成功后，将访客诊断缓存同步至会员账户
 */
import { bindDiagnosisRecords, syncGuestDiagnoses } from '@/api/frontend/compliance/member'
import {
  clearGuestDiagnosisCache,
  collectAnonymousDiagnosisIds,
  listGuestDiagnoses
} from './guestDiagnosisCache'

/** 同步访客诊断到数据库，返回成功同步条数 */
export async function syncGuestDiagnosisAfterLogin(): Promise<number> {
  let synced = 0

  const guestRecords = listGuestDiagnoses()
  if (guestRecords.length > 0) {
    try {
      const res = await syncGuestDiagnoses(guestRecords.map((r) => r.payload))
      synced += res?.count ?? 0
    } catch {
      // 同步失败保留缓存，下次登录重试
      return synced
    }
  }

  const anonymousIds = collectAnonymousDiagnosisIds()
  if (anonymousIds.length > 0) {
    try {
      const res = await bindDiagnosisRecords(anonymousIds)
      synced += res?.boundCount ?? 0
    } catch {
      if (synced === 0) return 0
    }
  }

  clearGuestDiagnosisCache()
  return synced
}
