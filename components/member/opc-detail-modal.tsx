'use client';

import { useEffect, useState } from 'react';
import { X } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { ProtectedImage } from '@/components/media/protected-image';
import { ApiClientError, memberFetch } from '@/lib/api/client';

type AttachmentBrief = {
  id: string;
  fileName: string;
  mimeType?: string;
  mediaUrl?: string;
  downloadUrl?: string;
};

interface OpcFullDetail {
  companyName?: string | null;
  proposedNames?: string[] | null;
  creditCode?: string | null;
  opcStatus?: string;
  statusLabel?: string;
  registeredCapital?: string | null;
  capitalTermYears?: number | null;
  businessTermType?: string | null;
  businessTermEnd?: string | null;
  businessScope?: string | null;
  registerProvince?: string | null;
  registerCity?: string | null;
  registerDistrict?: string | null;
  registerAddress?: string | null;
  legalPersonName?: string | null;
  idCard?: string | null;
  idCardValidFrom?: string | null;
  idCardValidTo?: string | null;
  ethnicity?: string | null;
  householdAddress?: string | null;
  residentialAddress?: string | null;
  phone?: string | null;
  email?: string | null;
  establishedAt?: string | null;
  taxActivatedAt?: string | null;
  bankName?: string | null;
  bankAccount?: string | null;
  rejectNote?: string | null;
  attachments?: {
    addressProof?: AttachmentBrief | null;
    idCardFront?: AttachmentBrief | null;
    idCardBack?: AttachmentBrief | null;
    license?: AttachmentBrief | null;
    bankReceipt?: AttachmentBrief | null;
  };
}

function DetailRow({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div className="grid grid-cols-[7rem_1fr] gap-2 border-b py-2 text-sm last:border-0">
      <span className="text-muted-foreground">{label}</span>
      <span className="break-all">{value ?? '—'}</span>
    </div>
  );
}

function AttachmentImage({ label, file }: { label: string; file?: AttachmentBrief | null }) {
  if (!file?.mediaUrl) return null;
  return (
    <div className="space-y-1">
      <p className="text-sm text-muted-foreground">{label}</p>
      <ProtectedImage
        src={file.mediaUrl}
        alt={label}
        className="max-h-40 rounded border object-contain"
      />
    </div>
  );
}

export function OpcDetailModal({ open, onClose }: { open: boolean; onClose: () => void }) {
  const [detail, setDetail] = useState<OpcFullDetail | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!open) return;
    setLoading(true);
    setError(null);
    memberFetch<OpcFullDetail>('/api/member/compliance/opc/detail')
      .then(setDetail)
      .catch((err) => {
        setError(err instanceof ApiClientError ? err.message : '加载失败');
      })
      .finally(() => setLoading(false));
  }, [open]);

  if (!open) return null;

  const registerAddr = [detail?.registerProvince, detail?.registerCity, detail?.registerDistrict, detail?.registerAddress]
    .filter(Boolean)
    .join(' ');

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/50" onClick={onClose} aria-hidden />
      <div className="relative z-10 flex max-h-[90vh] w-full max-w-2xl flex-col rounded-lg border bg-background shadow-lg">
        <div className="flex items-center justify-between border-b px-4 py-3">
          <h2 className="text-lg font-semibold">OPC 主体详情</h2>
          <Button variant="ghost" size="icon" onClick={onClose} aria-label="关闭">
            <X className="h-5 w-5" />
          </Button>
        </div>
        <div className="overflow-y-auto px-4 py-3">
          {loading && <p className="text-sm text-muted-foreground">加载中…</p>}
          {error && <p className="text-sm text-destructive">{error}</p>}
          {detail && !loading && (
            <div className="space-y-4">
              <div>
                <div className="mb-2 flex items-center gap-2">
                  <span className="font-medium">{detail.companyName ?? detail.proposedNames?.[0] ?? '设立中'}</span>
                  <Badge variant={detail.opcStatus === 'active' ? 'success' : 'warning'}>
                    {detail.statusLabel}
                  </Badge>
                </div>
                <DetailRow label="统一社会信用代码" value={detail.creditCode} />
                <DetailRow label="备选名称" value={detail.proposedNames?.join(' / ')} />
                <DetailRow label="注册资本（万元）" value={detail.registeredCapital} />
                <DetailRow label="认缴期限（年）" value={detail.capitalTermYears} />
                <DetailRow label="经营范围" value={detail.businessScope} />
                <DetailRow label="注册地址" value={registerAddr || null} />
                <DetailRow label="核准日期" value={detail.establishedAt} />
              </div>

              <div>
                <p className="mb-1 text-sm font-medium">法人信息</p>
                <DetailRow label="法人姓名" value={detail.legalPersonName} />
                <DetailRow label="身份证号" value={detail.idCard} />
                <DetailRow
                  label="身份证有效期"
                  value={
                    detail.idCardValidFrom && detail.idCardValidTo
                      ? `${detail.idCardValidFrom} ~ ${detail.idCardValidTo}`
                      : null
                  }
                />
                <DetailRow label="民族" value={detail.ethnicity} />
                <DetailRow label="户籍地址" value={detail.householdAddress} />
                <DetailRow label="现居住地址" value={detail.residentialAddress} />
                <DetailRow label="手机号" value={detail.phone} />
                <DetailRow label="邮箱" value={detail.email} />
              </div>

              <div>
                <p className="mb-1 text-sm font-medium">税务与银行</p>
                <DetailRow label="税务激活日期" value={detail.taxActivatedAt} />
                <DetailRow label="开户银行" value={detail.bankName} />
                <DetailRow label="对公账户" value={detail.bankAccount} />
              </div>

              {detail.rejectNote && <DetailRow label="驳回原因" value={detail.rejectNote} />}

              <div className="grid gap-3 sm:grid-cols-2">
                <AttachmentImage label="地址证明" file={detail.attachments?.addressProof} />
                <AttachmentImage label="身份证正面" file={detail.attachments?.idCardFront} />
                <AttachmentImage label="身份证反面" file={detail.attachments?.idCardBack} />
                <AttachmentImage label="营业执照" file={detail.attachments?.license} />
                <AttachmentImage label="开户回执" file={detail.attachments?.bankReceipt} />
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
