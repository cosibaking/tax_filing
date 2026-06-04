'use client';

import { useCallback, useEffect, useState } from 'react';
import Link from 'next/link';
import { useParams } from 'next/navigation';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import type { OpcAdminAction } from '@/lib/services/compliance/opc/opc-types';

type AttachmentBrief = { id: string; fileName: string; downloadUrl: string };

type TaskDetail = Record<string, unknown> & {
  statusLabel?: string;
  allowedActions?: OpcAdminAction[];
  hasMaterials?: boolean;
  materialsSubmittedAt?: string | null;
  member?: { phone?: string; name?: string | null };
  plan?: { name?: string } | null;
  attachments?: {
    addressProof?: AttachmentBrief | null;
    idCardFront?: AttachmentBrief | null;
    idCardBack?: AttachmentBrief | null;
    license?: AttachmentBrief | null;
    bankReceipt?: AttachmentBrief | null;
  };
};

async function adminFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const token = localStorage.getItem('admin_token');
  const res = await fetch(path, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...init?.headers,
    },
  });
  const json = await res.json();
  if (json.code !== 0) throw new Error(json.message || '请求失败');
  return json.data as T;
}

async function adminUpload(file: File): Promise<string> {
  const form = new FormData();
  form.append('file', file);
  const res = await fetch('/api/upload', { method: 'POST', body: form });
  const json = await res.json();
  if (json.code !== 0) throw new Error(json.message || '上传失败');
  return String(json.data.id);
}

function Field({ label, value }: { label: string; value?: string | null }) {
  return (
    <p className="text-sm">
      <span className="text-muted-foreground">{label}：</span>
      {value ?? '—'}
    </p>
  );
}

export default function OpcTaskDetailPage() {
  const params = useParams();
  const opcId = params.opcId as string;
  const [detail, setDetail] = useState<TaskDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [note, setNote] = useState('');
  const [companyName, setCompanyName] = useState('');
  const [creditCode, setCreditCode] = useState('');
  const [establishedAt, setEstablishedAt] = useState('');
  const [licenseFileId, setLicenseFileId] = useState('');
  const [taxActivatedAt, setTaxActivatedAt] = useState('');
  const [bankAccount, setBankAccount] = useState('');
  const [acting, setActing] = useState(false);

  const load = useCallback(() => {
    setLoading(true);
    adminFetch<TaskDetail>(`/api/admin/compliance/opc-tasks/${opcId}`)
      .then((d) => {
        setDetail(d);
        setCompanyName((d.companyName as string) ?? '');
        setCreditCode((d.creditCode as string) ?? '');
      })
      .catch((err) => setError(err instanceof Error ? err.message : '加载失败'))
      .finally(() => setLoading(false));
  }, [opcId]);

  useEffect(() => {
    load();
  }, [load]);

  const runAction = async (action: OpcAdminAction, payload: Record<string, string> = {}) => {
    setActing(true);
    setError(null);
    try {
      await adminFetch('/api/admin/compliance/opc-tasks', {
        method: 'PATCH',
        body: JSON.stringify({ opcId, action, note, payload }),
      });
      setNote('');
      load();
    } catch (err) {
      setError(err instanceof Error ? err.message : '操作失败');
    } finally {
      setActing(false);
    }
  };

  if (loading) return <p className="text-muted-foreground">加载中…</p>;
  if (!detail) return null;

  const allowed = detail.allowedActions ?? [];
  const proposedNames = detail.proposedNames as string[] | undefined;
  const att = detail.attachments;

  const downloadAttachment = async (file: AttachmentBrief) => {
    const token = localStorage.getItem('admin_token');
    const res = await fetch(file.downloadUrl, {
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    });
    if (!res.ok) throw new Error('下载失败，请重新登录管理端');
    const blob = await res.blob();
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = file.fileName;
    a.click();
    URL.revokeObjectURL(url);
  };

  const AttachmentLink = ({ file }: { file?: AttachmentBrief | null }) =>
    file ? (
      <button
        type="button"
        onClick={() => downloadAttachment(file).catch((e) => setError(e instanceof Error ? e.message : '下载失败'))}
        className="text-sm text-primary underline"
      >
        {file.fileName}
      </button>
    ) : (
      <span className="text-sm text-muted-foreground">—</span>
    );

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Button variant="ghost" size="sm" asChild>
          <Link href="/admin/compliance/opc-tasks">← 返回列表</Link>
        </Button>
        <h1 className="text-2xl font-bold">OPC 任务详情</h1>
        <Badge>{detail.statusLabel as string}</Badge>
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {!detail.hasMaterials && (
        <Alert>
          <AlertDescription>
            该客户尚未成功提交注册资料（可能提交失败或未上传附件）。请通知主播在「OPC 落地进度」页重新填写并提交，提交成功后会显示提交时间。
          </AlertDescription>
        </Alert>
      )}

      <Card>
        <CardHeader>
          <CardTitle>客户信息</CardTitle>
        </CardHeader>
        <CardContent className="space-y-1">
          <Field label="会员手机" value={detail.member?.phone} />
          <Field label="会员姓名" value={detail.member?.name} />
          <Field label="套餐" value={detail.plan?.name} />
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>拟设公司</CardTitle>
        </CardHeader>
        <CardContent className="space-y-1">
          <Field label="备选名称" value={proposedNames?.join(' / ')} />
          <Field label="注册资本(万)" value={detail.registeredCapital != null ? String(detail.registeredCapital) : null} />
          <Field label="认缴期限(年)" value={detail.capitalTermYears != null ? String(detail.capitalTermYears) : null} />
          <Field label="经营范围" value={detail.businessScope as string} />
          <Field
            label="注册地址"
            value={
              [detail.registerProvince, detail.registerCity, detail.registerDistrict, detail.registerAddress]
                .filter(Boolean)
                .join(' ') || null
            }
          />
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>法人信息（脱敏）</CardTitle>
        </CardHeader>
        <CardContent className="space-y-1">
          <Field label="法人姓名" value={detail.legalPersonName as string} />
          <Field label="身份证" value={detail.idCardMasked as string} />
          <Field
            label="证件有效期"
            value={
              detail.idCardValidFrom && detail.idCardValidTo
                ? `${detail.idCardValidFrom} 至 ${detail.idCardValidTo}`
                : null
            }
          />
          <Field label="手机" value={detail.phone as string} />
          <Field label="邮箱" value={detail.email as string} />
          <Field label="户籍地址" value={detail.householdAddress as string} />
          <Field label="现居地址" value={detail.residentialAddress as string} />
          <Field
            label="资料提交时间"
            value={
              detail.materialsSubmittedAt
                ? detail.materialsSubmittedAt.slice(0, 19).replace('T', ' ')
                : null
            }
          />
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>附件材料</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2 text-sm">
          <p>
            <span className="text-muted-foreground">地址证明：</span>
            <AttachmentLink file={att?.addressProof} />
          </p>
          <p>
            <span className="text-muted-foreground">身份证正面：</span>
            <AttachmentLink file={att?.idCardFront} />
          </p>
          <p>
            <span className="text-muted-foreground">身份证反面：</span>
            <AttachmentLink file={att?.idCardBack} />
          </p>
          <p>
            <span className="text-muted-foreground">营业执照：</span>
            <AttachmentLink file={att?.license} />
          </p>
          <p>
            <span className="text-muted-foreground">开户回执：</span>
            <AttachmentLink file={att?.bankReceipt} />
          </p>
        </CardContent>
      </Card>

      {Boolean(detail.companyName || detail.creditCode) && (
        <Card>
          <CardHeader>
            <CardTitle>工商结果</CardTitle>
          </CardHeader>
          <CardContent className="space-y-1">
            <Field label="核准名称" value={detail.companyName as string} />
            <Field label="信用代码" value={detail.creditCode as string} />
          </CardContent>
        </Card>
      )}

      {allowed.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>推进进度</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <Label>备注 / 驳回原因</Label>
              <Input className="mt-1" value={note} onChange={(e) => setNote(e.target.value)} placeholder="驳回时至少 10 字" />
            </div>

            {allowed.includes('approve_materials') && (
              <Button disabled={acting} onClick={() => runAction('approve_materials')}>
                通过资料审核
              </Button>
            )}

            {allowed.includes('reject_materials') && (
              <Button variant="destructive" disabled={acting} onClick={() => runAction('reject_materials')}>
                驳回资料
              </Button>
            )}

            {allowed.includes('issue_license') && (
              <div className="space-y-3 rounded-md border p-4">
                <p className="text-sm font-medium">执照已下发</p>
                <div>
                  <Label>核准公司名称</Label>
                  <Input className="mt-1" value={companyName} onChange={(e) => setCompanyName(e.target.value)} />
                </div>
                <div>
                  <Label>统一社会信用代码</Label>
                  <Input className="mt-1" value={creditCode} onChange={(e) => setCreditCode(e.target.value)} />
                </div>
                <div>
                  <Label>成立日期</Label>
                  <Input type="date" className="mt-1" value={establishedAt} onChange={(e) => setEstablishedAt(e.target.value)} />
                </div>
                <div>
                  <Label>营业执照附件</Label>
                  <Input
                    type="file"
                    className="mt-1"
                    accept="image/*,.pdf"
                    onChange={async (e) => {
                      const f = e.target.files?.[0];
                      if (!f) return;
                      try {
                        const id = await adminUpload(f);
                        setLicenseFileId(id);
                      } catch (err) {
                        setError(err instanceof Error ? err.message : '上传失败');
                      }
                    }}
                  />
                  {licenseFileId && <p className="text-xs text-muted-foreground">已上传 #{licenseFileId}</p>}
                </div>
                <Button
                  disabled={acting}
                  onClick={() =>
                    runAction('issue_license', {
                      companyName,
                      creditCode,
                      establishedAt,
                      licenseFileId,
                    })
                  }
                >
                  确认执照已下发
                </Button>
              </div>
            )}

            {allowed.includes('complete_tax') && (
              <div className="space-y-3 rounded-md border p-4">
                <p className="text-sm font-medium">税务登记完成</p>
                <div>
                  <Label>电子税务局激活日期</Label>
                  <Input type="date" className="mt-1" value={taxActivatedAt} onChange={(e) => setTaxActivatedAt(e.target.value)} />
                </div>
                <Button
                  disabled={acting}
                  onClick={() => runAction('complete_tax', { taxActivatedAt })}
                >
                  税务登记完成
                </Button>
              </div>
            )}

            {allowed.includes('complete_bank') && (
              <div className="space-y-3 rounded-md border p-4">
                <p className="text-sm font-medium">银行开户完成</p>
                <div>
                  <Label>对公账号</Label>
                  <Input className="mt-1" value={bankAccount} onChange={(e) => setBankAccount(e.target.value)} />
                </div>
                <Button
                  disabled={acting}
                  onClick={() => runAction('complete_bank', { bankAccount })}
                >
                  银行开户完成
                </Button>
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  );
}
