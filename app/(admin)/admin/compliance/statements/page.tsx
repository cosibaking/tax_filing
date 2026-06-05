'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { withBasePath } from '@/lib/base-path';

interface StatementJob {
  memberId: string;
  phone: string;
  period: string;
  status: string;
}

async function adminFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const token = localStorage.getItem('admin_token');
  const res = await fetch(withBasePath(path), {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
  });
  const json = await res.json();
  if (json.code !== 0) throw new Error(json.message || '请求失败');
  return json.data as T;
}

export default function StatementsPage() {
  const [period, setPeriod] = useState('');
  const [items, setItems] = useState<StatementJob[]>([]);
  const [loading, setLoading] = useState(false);
  const [sending, setSending] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [message, setMessage] = useState<string | null>(null);

  const load = () => {
    setLoading(true);
    adminFetch<StatementJob[]>('/api/admin/compliance/statements')
      .then(setItems)
      .catch((err) => setError(err instanceof Error ? err.message : '加载失败'))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    load();
  }, []);

  const sendStatements = async () => {
    if (!period) {
      setError('请输入账期，如 2026-05');
      return;
    }
    setSending(true);
    setError(null);
    setMessage(null);
    try {
      await adminFetch('/api/admin/compliance/statements/send', {
        method: 'POST',
        body: JSON.stringify({ period }),
      });
      setMessage('对账单生成/发送任务已提交');
      load();
    } catch (err) {
      setError(err instanceof Error ? err.message : '发送失败');
    } finally {
      setSending(false);
    }
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">对账单管理</h1>
      <Card>
        <CardHeader>
          <CardTitle className="text-base">批量生成/发送</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col gap-4 sm:flex-row sm:items-end">
          <div className="flex-1">
            <Label htmlFor="period">账期</Label>
            <Input
              id="period"
              className="mt-1"
              placeholder="2026-05"
              value={period}
              onChange={(e) => setPeriod(e.target.value)}
            />
          </div>
          <Button onClick={sendStatements} disabled={sending}>
            {sending ? '处理中…' : '生成并发送'}
          </Button>
        </CardContent>
      </Card>
      {message && (
        <Alert variant="info">
          <AlertDescription>{message}</AlertDescription>
        </Alert>
      )}
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      {loading && <p className="text-muted-foreground">加载中…</p>}
      <div className="space-y-2">
        {items.map((s) => (
          <Card key={`${s.memberId}-${s.period}`}>
            <CardContent className="flex justify-between py-4">
              <div>
                <p className="font-medium">{s.period}</p>
                <p className="text-sm text-muted-foreground">{s.phone}</p>
              </div>
              <span className="text-sm text-muted-foreground">{s.status}</span>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
