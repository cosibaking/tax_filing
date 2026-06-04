'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Select } from '@/components/ui/select';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';

interface OpcTask {
  id: string;
  phone: string;
  companyName?: string;
  status: string;
  currentStep: string;
}

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

export default function OpcTasksPage() {
  const [tasks, setTasks] = useState<OpcTask[]>([]);
  const [statusFilter, setStatusFilter] = useState('');
  const [search, setSearch] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = () => {
    setLoading(true);
    const params = new URLSearchParams();
    if (statusFilter) params.set('status', statusFilter);
    if (search) params.set('search', search);
    adminFetch<OpcTask[]>(`/api/admin/compliance/opc-tasks?${params}`)
      .then(setTasks)
      .catch((err) => setError(err instanceof Error ? err.message : '加载失败'))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    load();
  }, [statusFilter, search]);

  const updateStep = async (id: string, step: string) => {
    try {
      await adminFetch(`/api/admin/compliance/opc-tasks`, {
        method: 'PATCH',
        body: JSON.stringify({ id, currentStep: step }),
      });
      load();
    } catch (err) {
      setError(err instanceof Error ? err.message : '更新失败');
    }
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">OPC 注册任务</h1>
      <div className="flex flex-col gap-3 sm:flex-row">
        <Input placeholder="搜索手机号/公司名" value={search} onChange={(e) => setSearch(e.target.value)} className="max-w-xs" />
        <Select value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)} className="max-w-xs">
          <option value="">全部状态</option>
          <option value="materials">材料收集</option>
          <option value="business">工商注册</option>
          <option value="tax">税务登记</option>
          <option value="bank">银行开户</option>
          <option value="active">已激活</option>
        </Select>
      </div>
      {loading && <p className="text-muted-foreground">加载中…</p>}
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      <div className="space-y-3">
        {tasks.map((t) => (
          <Card key={t.id}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-base">{t.companyName ?? t.phone}</CardTitle>
              <Badge>{t.status}</Badge>
            </CardHeader>
            <CardContent className="flex flex-wrap items-center gap-2">
              <span className="text-sm text-muted-foreground">当前节点：{t.currentStep}</span>
              <Button size="sm" variant="outline" onClick={() => updateStep(t.id, 'business')}>
                推进至工商
              </Button>
              <Button size="sm" variant="outline" onClick={() => updateStep(t.id, 'active')}>
                标记激活
              </Button>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
