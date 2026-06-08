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
