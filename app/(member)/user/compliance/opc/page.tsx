'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { OpcMaterialsForm } from '@/components/compliance/opc-materials-form';
import { OpcFileUpload } from '@/components/compliance/opc-file-upload';
import { ApiClientError, memberFetch } from '@/lib/api/client';
import type { OpcMaterialsInput } from '@/lib/services/compliance/opc/opc-types';

interface OpcProgress {
  opcStatus?: string;
  statusLabel?: string;
  companyName?: string;
  creditCode?: string;
  steps?: { key: string; label: string; status: 'done' | 'current' | 'pending'; date?: string }[];
  canEditMaterials?: boolean;
  rejectNote?: string | null;
  bankOpeningChecklist?: string[];
  legalPersonName?: string;
  phone?: string;
  proposedNamePrimary?: string;
}

export default function OpcPage() {
  const [progress, setProgress] = useState<OpcProgress | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [bankName, setBankName] = useState('');
  const [bankReceiptFileId, setBankReceiptFileId] = useState('');

  const load = () => {
    setLoading(true);
    memberFetch<OpcProgress | null>('/api/member/compliance/opc/progress')
      .then((data) => setProgress(data ?? null))
      .catch((err) => {
        setError(err instanceof ApiClientError ? err.message : '加载失败');
      })
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    load();
  }, []);

  const submitMaterials = async (data: OpcMaterialsInput) => {
    setSubmitting(true);
    setError(null);
    try {
      await memberFetch('/api/member/compliance/opc/materials', {
        method: 'POST',
        body: JSON.stringify(data),
      });
      load();
    } catch (err) {
      setError(err instanceof ApiClientError ? err.message : '提交失败');
    } finally {
      setSubmitting(false);
    }
  };

  const submitBankReceipt = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!bankReceiptFileId) {
      setError('请上传开户回执');
      return;
    }
    setSubmitting(true);
    setError(null);
    try {
      await memberFetch('/api/member/compliance/opc/bank-receipt', {
        method: 'POST',
        body: JSON.stringify({ bankName, bankReceiptFileId }),
      });
      setBankReceiptFileId('');
      load();
    } catch (err) {
      setError(err instanceof ApiClientError ? err.message : '上传失败');
    } finally {
      setSubmitting(false);
    }
  };

  const steps = progress?.steps ?? [];
  const showBankSection = progress?.opcStatus === 'tax' || progress?.opcStatus === 'bank';

  if (loading) return <p className="text-muted-foreground">加载中…</p>;

  if (!progress) {
    return (
      <Alert>
        <AlertDescription>
          请先完成{' '}
          <Link href="/user/compliance/plan" className="underline">
            方案签约
          </Link>{' '}
          后再提交 OPC 注册资料。
        </AlertDescription>
      </Alert>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center gap-3">
        <h1 className="text-2xl font-bold">OPC 落地进度</h1>
        <Badge variant={progress.opcStatus === 'active' ? 'success' : 'warning'}>
          {progress.statusLabel ?? progress.opcStatus}
        </Badge>
        <span className="text-sm text-muted-foreground">预计 14 个工作日</span>
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {progress.rejectNote && (
        <Alert variant="destructive">
          <AlertDescription>资料被退回：{progress.rejectNote}</AlertDescription>
        </Alert>
      )}

      <Card>
        <CardHeader>
          <CardTitle>进度时间线</CardTitle>
        </CardHeader>
        <CardContent>
          <ol className="relative border-l border-muted pl-6">
            {steps.map((step, i) => (
              <li key={step.key} className={`mb-6 ${i === steps.length - 1 ? 'mb-0' : ''}`}>
                <span
                  className={`absolute -left-2 flex h-4 w-4 rounded-full ${
                    step.status === 'done'
                      ? 'bg-primary'
                      : step.status === 'current'
                        ? 'bg-amber-500'
                        : 'bg-muted'
                  }`}
                />
                <p className="font-medium">{step.label}</p>
                {step.date && <p className="text-xs text-muted-foreground">{step.date}</p>}
              </li>
            ))}
          </ol>
          {progress.companyName && (
            <p className="mt-4 text-sm">
              <span className="text-muted-foreground">核准名称：</span>
              {progress.companyName}
            </p>
          )}
          {progress.creditCode && (
            <p className="text-sm">
              <span className="text-muted-foreground">统一社会信用代码：</span>
              {progress.creditCode}
            </p>
          )}
        </CardContent>
      </Card>

      {progress.canEditMaterials && (
        <Card>
          <CardHeader>
            <CardTitle>注册资料</CardTitle>
          </CardHeader>
          <CardContent>
            <OpcMaterialsForm
              onSubmit={submitMaterials}
              submitting={submitting}
              onValidationError={setError}
            />
          </CardContent>
        </Card>
      )}

      {!progress.canEditMaterials && progress.opcStatus !== 'active' && (
        <Card>
          <CardHeader>
            <CardTitle>已提交资料</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2 text-sm">
            <p>
              <span className="text-muted-foreground">拟设名称：</span>
              {progress.proposedNamePrimary ?? '—'}
            </p>
            <p>
              <span className="text-muted-foreground">法人：</span>
              {progress.legalPersonName ?? '—'}
            </p>
            <p>
              <span className="text-muted-foreground">联系手机：</span>
              {progress.phone ?? '—'}
            </p>
            <p className="text-muted-foreground">资料审核中，请等待顾问处理。</p>
          </CardContent>
        </Card>
      )}

      {showBankSection && (
        <Card>
          <CardHeader>
            <CardTitle>银行开户</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <details className="text-sm">
              <summary className="cursor-pointer font-medium">开户材料清单</summary>
              <ul className="mt-2 list-inside list-disc text-muted-foreground">
                {(progress.bankOpeningChecklist ?? []).map((item) => (
                  <li key={item}>{item}</li>
                ))}
              </ul>
            </details>
            <form onSubmit={submitBankReceipt} className="space-y-4">
              <div>
                <Label htmlFor="bankName">开户银行（选填）</Label>
                <Input
                  id="bankName"
                  className="mt-1"
                  value={bankName}
                  onChange={(e) => setBankName(e.target.value)}
                />
              </div>
              <OpcFileUpload
                label="开户回执 / 基本存款账户信息"
                value={bankReceiptFileId}
                onChange={setBankReceiptFileId}
              />
              <Button type="submit" disabled={submitting}>
                {submitting ? '提交中…' : '上传开户回执'}
              </Button>
            </form>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
