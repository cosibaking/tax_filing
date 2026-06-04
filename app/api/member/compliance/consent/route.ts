import { jsonOk, jsonFail, handleApiError } from '@/lib/api/envelope';
import {
  ErrorCodes,
  CONSENT_TYPES,
  LEGAL_DOCUMENT_VERSIONS,
  ORDER_LINKED_CONSENT_TYPES,
  isDev,
  type ConsentType,
} from '@/lib/api/constants';
import { getMemberFromRequest } from '@/lib/auth/member';
import { recordConsent } from '@/lib/services/compliance/order/order-service';

function isConsentType(value: string): value is ConsentType {
  return (CONSENT_TYPES as readonly string[]).includes(value);
}

function resolveDocumentVersion(type: ConsentType, documentVersion?: string): string | null {
  const expected = LEGAL_DOCUMENT_VERSIONS[type];
  if (!documentVersion) return expected;
  return documentVersion === expected ? documentVersion : null;
}

export async function POST(request: Request) {
  try {
    const member = await getMemberFromRequest(request);
    const body = await request.json();
    const { type, orderId, documentVersion } = body as {
      type?: string;
      orderId?: string;
      documentVersion?: string;
    };

    let resolvedType: ConsentType;
    let resolvedVersion: string;

    if (isDev) {
      resolvedType = type && isConsentType(type) ? type : 'plan_confirm';
      resolvedVersion =
        resolveDocumentVersion(resolvedType, documentVersion) ??
        LEGAL_DOCUMENT_VERSIONS[resolvedType];
    } else {
      if (!type || !isConsentType(type)) {
        return jsonFail(ErrorCodes.VALIDATION, '同意类型无效');
      }
      if (!documentVersion) {
        return jsonFail(ErrorCodes.VALIDATION, '缺少文档版本号');
      }
      const version = resolveDocumentVersion(type, documentVersion);
      if (!version) {
        return jsonFail(ErrorCodes.VALIDATION, '文档版本号无效');
      }
      if (
        (ORDER_LINKED_CONSENT_TYPES as readonly string[]).includes(type) &&
        !orderId
      ) {
        return jsonFail(ErrorCodes.VALIDATION, '缺少订单号');
      }
      resolvedType = type;
      resolvedVersion = version;
    }

    await recordConsent(
      member.id,
      resolvedType,
      resolvedVersion,
      request,
      orderId ? BigInt(orderId) : undefined,
    );

    return jsonOk({ recorded: true });
  } catch (error) {
    return handleApiError(error);
  }
}
