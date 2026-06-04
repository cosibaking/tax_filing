'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { ApiClientError, memberFetch } from '@/lib/api/client';

interface OpcProgress {
  opcStatus?: string;
  steps?: { key: string; label: string; status: 'done' | 'current' | 'pending'; date?: string }[];
}

export default function OpcPage() {
  const [progress, setProgress] = useState<OpcProgress | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [materialNote, setMaterialNote] = useState('');

  const load = () => {
    setLoading(true);
    memberFetch<OpcProgress>('/api/member/compliance/opc/progress')
      .then(setProgress)
      .catch((err) => {
        setError(err instanceof ApiClientError ? err.message : '加载失败');
      })
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    load();
  }, []);

  const submitMaterials = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError(null);
    try {
      await memberFetch('/api/member/compliance/opc/materials', {
        method: 'POST',
        body: JSON.stringify({ note: materialNote }),
      });
      setMaterialNote('');
      load();
    } catch (err) {
      setError(err instanceof ApiClientError ? err.message : '提交失败');
    } finally {
      setSubmitting(false);
    }
  };

  const steps = progress?.steps ?? [
    { key: 'materials', label: '材料收集', status: 'done' as const },
    { key: 'business', label: '工商注册', status: 'current' as const },
    { key: 'tax', label: '税务登记', status: 'pending' as const },
    { key: 'bank', label: '银行开户', status: 'pending' as const },
    { key: 'active', label: '激活完成', status: 'pending' as const },
  ];

  if (loading) return <p className="text-muted-foreground">加载中…</p>;

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <h1 className="text-2xl font-bold">OPC 落地进度</h1>
        <Badge variant={progress?.opcStatus === 'active' ? 'success' : 'warning'}>
          {progress?.opcStatus === 'active' ? '已激活' : '设立中'}
        </Badge>
      </div>
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
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
                  className={`absolute -left-2 flex h-4 w-4 items-center justify-center rounded-full ${
                    step.status === 'done'
                      ? 'bg-primary'
                      : step.status === 'current'
                        ? 'bg-amber-500'
                        : 'bg-muted'
                  }`}
                />
                <p className="font-medium">{step.label}</p>
                {step.date && (
                  <p className="text-xs text-muted-foreground">{step.date}</p>
                )}
              </li>
            ))}
          </ol>
        </CardContent>
      </Card>

      {progress?.opcStatus !== 'active' && (
        <Card>
          <CardHeader>
            <CardTitle>补充材料</CardTitle>
          </CardHeader>
          <CardContent>
            <form onSubmit={submitMaterials} className="space-y-4">
              <div>
                <Label htmlFor="note">材料说明</Label>
                <Input
                  id="note"
                  className="mt-1"
                  placeholder="如：身份证已上传、待补充地址证明"
                  value={materialNote}
                  onChange={(e) => setMaterialNote(e.target.value)}
                />
              </div>
              <Button type="submit" disabled={submitting}>
                {submitting ? '提交中…' : '提交材料说明'}
              </Button>
            </form>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
