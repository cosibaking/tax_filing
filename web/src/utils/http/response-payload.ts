import type { AxiosResponse } from 'axios'

export class BlobPayloadParseError extends Error {
  readonly response: AxiosResponse
  readonly originalError?: unknown

  constructor(response: AxiosResponse, cause: unknown, originalError?: unknown) {
    super('响应数据解析失败', { cause })
    this.name = 'BlobPayloadParseError'
    this.response = response
    this.originalError = originalError
  }
}

export function isJsonMediaType(contentType: string): boolean {
  const mediaType = contentType.split(';', 1)[0].trim().toLowerCase()
  return (
    mediaType === 'application/json' ||
    mediaType === 'text/json' ||
    (mediaType.startsWith('application/') && mediaType.endsWith('+json'))
  )
}

export function shouldHandleUnauthorized(
  payloadCode: number | undefined,
  httpStatus: number | undefined,
  unauthorizedCode: number
): boolean {
  return httpStatus === 401 || httpStatus === unauthorizedCode || payloadCode === unauthorizedCode
}

export async function normalizeResponsePayload(response: AxiosResponse, originalError?: unknown) {
  const contentType = response.headers['content-type'] || ''
  const payload = response.data as unknown
  if (!isJsonMediaType(contentType) || typeof Blob === 'undefined' || !(payload instanceof Blob)) {
    return payload
  }

  try {
    response.data = JSON.parse(await payload.text())
  } catch (cause) {
    throw new BlobPayloadParseError(response, cause, originalError)
  }
  return response.data
}
