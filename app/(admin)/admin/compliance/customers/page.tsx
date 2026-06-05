'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { withBasePath } from '@/lib/base-path';

interface Customer {
  id: string;
  phone: string;
  nickname?: string;
  opcStatus?: string;
  planName?: string;
}

async function adminFetch<T>(path: string): Promise<T> {
  const token = localStorage.getItem('admin_token');
  const res = await fetch(withBasePath(path), {
    headers: token ? { Authorization: `Bearer ${token}` } : {},
  });
  const json = await res.json();
  if (json.code !== 0) throw new Error(json.message || '请求失败');
  return json.data as T;
}

export default function CustomersPage() {
  const [items, setItems] = useState<Customer[]>([]);
  const [search, setSearch] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const q = search ? `?search=${encodeURIComponent(search)}` : '';
    setLoading(true);
    adminFetch<Customer[]>(`/api/admin/compliance/customers${q}`)
      .then(setItems)
      .catch((err) => setError(err instanceof Error ? err.message : '加载失败'))
      .finally(() => setLoading(false));
  }, [search]);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">合规客户列表</h1>
      <Input
        placeholder="搜索手机号或昵称"
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        className="max-w-sm"
      />
      {loading && <p className="text-muted-foreground">加载中…</p>}
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      <div className="space-y-2">
        {items.map((c) => (
          <Card key={c.id}>
            <CardContent className="flex flex-col gap-2 py-4 sm:flex-row sm:items-center sm:justify-between">
              <div>
                <p className="font-medium">{c.phone}</p>
                <p className="text-sm text-muted-foreground">{c.nickname ?? '—'}</p>
              </div>
              <div className="flex gap-2">
                <Badge>{c.opcStatus ?? '未知'}</Badge>
                {c.planName && <Badge variant="secondary">{c.planName}</Badge>}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
