/**
 * OPC 资料真实性校验开关（OCR / 三要素 / 手机实名等）
 * MVP 默认 Mock，生产接入第三方后将 VITE_COMPLIANCE_VERIFY_MOCK 设为 false
 */
export function isComplianceVerifyMock(): boolean {
  const flag = import.meta.env.VITE_COMPLIANCE_VERIFY_MOCK
  if (flag === 'false') return false
  if (flag === 'true') return true
  return import.meta.env.DEV
}

/** 身份证 Mock 模式：仅校验 18 位格式，跳过校验位 */
export const ID_CARD_FORMAT_PATTERN =
  /^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$/

/** 手机号 Mock 模式：大陆 11 位即可 */
export const PHONE_MOCK_PATTERN = /^1\d{10}$/

/** 信用代码 Mock 模式：非空且 ≤18 位（未对接工商公开库） */
export function validateCreditCodeMock(code: string): boolean {
  const trimmed = code.trim()
  return trimmed.length > 0 && trimmed.length <= 18
}

const CREDIT_CODE_FORMAT_PATTERN =
  /^[0-9A-HJ-NPQRTUWXY]{2}\d{6}[0-9A-HJ-NPQRTUWXY]{10}$/
const CREDIT_CODE_CHARSET = '0123456789ABCDEFGHJKLMNPQRTUWXY'
const CREDIT_CODE_WEIGHTS = [1, 3, 9, 27, 19, 26, 16, 17, 20, 29, 25, 13, 8, 24, 10, 30, 28]

/** 信用代码生产模式：GB 32100 格式 + 校验位 */
export function validateCreditCode(code: string): boolean {
  const c = code.trim().toUpperCase()
  if (!CREDIT_CODE_FORMAT_PATTERN.test(c)) return false
  let sum = 0
  for (let i = 0; i < 17; i++) {
    const idx = CREDIT_CODE_CHARSET.indexOf(c[i])
    if (idx < 0) return false
    sum += idx * CREDIT_CODE_WEIGHTS[i]
  }
  const checkIdx = (31 - (sum % 31)) % 31
  return CREDIT_CODE_CHARSET[checkIdx] === c[17]
}
