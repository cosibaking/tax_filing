'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Select } from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import { Checkbox } from '@/components/ui/checkbox';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { TAX_CHECKLIST_ITEMS } from '@/lib/api/constants';

interface FilingTask {
  id: string;
  memberPhone: string;
  taxType: string;
  period: string;
  status: string;
}

async function adminFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const token = localStorage.getItem('admin_token');
  const res = await fetch(path, {
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

export default function FilingPage() {
  const [tasks, setTasks] = useState<FilingTask[]>([]);
  const [statusFilter, setStatusFilter] = useState('pending');
  const [checklist, setChecklist] = useState<boolean[]>(TAX_CHECKLIST_ITEMS.map(() => false));
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = () => {
    setLoading(true);
    adminFetch<FilingTask[]>(`/api/admin/compliance/filing?status=${statusFilter}`)
      .then(setTasks)
      .catch((err) => setError(err instanceof Error ? err.message : '加载失败'))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    load();
  }, [statusFilter]);

  const allChecked = checklist.every(Boolean);

  const markFiled = async (id: string) => {
    if (!allChecked) {
      setError('请先完成全部 9 项自查');
      return;
    }
    try {
      await adminFetch(`/api/admin/compliance/filing/${id}`, {
        method: 'PATCH',
        body: JSON.stringify({ status: 'filed' }),
      });
      load();
    } catch (err) {
      setError(err instanceof Error ? err.message : '更新失败');
    }
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">申报工作台</h1>
      <Select value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)} className="max-w-xs">
        <option value="pending">待申报</option>
        <option value="filed">已申报</option>
        <option value="overdue">已逾期</option>
      </Select>
      {loading && <p className="text-muted-foreground">加载中…</p>}
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="text-base">申报前自查（9 项）</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            {TAX_CHECKLIST_ITEMS.map((item, i) => (
              <label key={item} className="flex items-start gap-2 text-sm">
                <Checkbox
                  checked={checklist[i]}
                  onCheckedChange={(v) => {
                    const next = [...checklist];
                    next[i] = !!v;
                    setChecklist(next);
                  }}
                />
                {item}
              </label>
            ))}
          </CardContent>
        </Card>
        <div className="space-y-3">
          {tasks.map((t) => (
            <Card key={t.id}>
              <CardContent className="flex flex-col gap-3 py-4 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <p className="font-medium">{t.taxType} · {t.period}</p>
                  <p className="text-sm text-muted-foreground">{t.memberPhone}</p>
                </div>
                <div className="flex items-center gap-2">
                  <Badge>{t.status}</Badge>
                  {t.status === 'pending' && (
                    <Button size="sm" disabled={!allChecked} onClick={() => markFiled(t.id)}>
                      标记已申报
                    </Button>
                  )}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </div>
  );
}
