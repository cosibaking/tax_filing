'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Select } from '@/components/ui/select';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';

interface OpcTask {
  id: string;
  status: string;
  statusLabel: string;
  companyName?: string | null;
  proposedName?: string | null;
  legalPersonName?: string | null;
  memberPhone: string;
  memberName?: string | null;
  materialsSubmittedAt?: string | null;
  slaDay?: number | null;
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
    if (search) params.set('q', search);
    adminFetch<OpcTask[]>(`/api/admin/compliance/opc-tasks?${params}`)
      .then(setTasks)
      .catch((err) => setError(err instanceof Error ? err.message : '加载失败'))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    load();
  }, [statusFilter, search]);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">OPC 注册任务</h1>
      <div className="flex flex-col gap-3 sm:flex-row">
        <Input
          placeholder="搜索手机号/公司名/法人"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="max-w-xs"
        />
        <Select value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)} className="max-w-xs">
          <option value="">全部状态</option>
          <option value="pending">待提交资料</option>
          <option value="materials">资料待补正</option>
          <option value="materials_review">资料审核中</option>
          <option value="registering">工商注册中</option>
          <option value="tax">税务登记中</option>
          <option value="bank">银行开户中</option>
          <option value="active">已激活</option>
        </Select>
        <Button variant="outline" onClick={load}>
          刷新
        </Button>
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
              <CardTitle className="text-base">
                {t.companyName ?? t.proposedName ?? t.legalPersonName ?? t.memberPhone}
              </CardTitle>
              <Badge>{t.statusLabel}</Badge>
            </CardHeader>
            <CardContent className="flex flex-wrap items-center gap-3 text-sm text-muted-foreground">
              <span>主播：{t.memberName ?? t.memberPhone}</span>
              {t.materialsSubmittedAt && (
                <span>提交：{t.materialsSubmittedAt.slice(0, 10)}</span>
              )}
              {t.slaDay != null && <span>第 {t.slaDay} 天</span>}
              <Button size="sm" asChild>
                <Link href={`/admin/compliance/opc-tasks/${t.id}`}>查看详情</Link>
              </Button>
            </CardContent>
          </Card>
        ))}
        {!loading && tasks.length === 0 && (
          <p className="text-muted-foreground">暂无任务</p>
        )}
      </div>
    </div>
  );
}
