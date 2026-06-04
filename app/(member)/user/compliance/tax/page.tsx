'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { TAX_CHECKLIST_ITEMS } from '@/lib/api/constants';
import { ApiClientError, memberFetch } from '@/lib/api/client';

interface TaxTask {
  id: string;
  taxType: string;
  period: string;
  dueDate: string;
  status: string;
}

export default function TaxPage() {
  const [tasks, setTasks] = useState<TaxTask[]>([]);
  const [checklist, setChecklist] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    Promise.all([
      memberFetch<{ tasks: TaxTask[] }>('/api/member/compliance/tax/calendar'),
      memberFetch<{ items: string[] }>('/api/member/compliance/tax/checklist').catch(() => null),
    ])
      .then(([cal, chk]) => {
        setTasks(cal.tasks ?? []);
        setChecklist(chk?.items ?? [...TAX_CHECKLIST_ITEMS]);
      })
      .catch((err) => {
        setError(err instanceof ApiClientError ? err.message : '加载失败');
      })
      .finally(() => setLoading(false));
  }, []);

  const statusLabel: Record<string, string> = {
    pending: '待申报',
    filed: '已申报',
    overdue: '已逾期',
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">申报日历</h1>
      {loading && <p className="text-muted-foreground">加载中…</p>}
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      {!loading && (
        <div className="grid gap-6 lg:grid-cols-2">
          <div className="space-y-3">
            <h2 className="font-semibold">申报任务</h2>
            {tasks.length === 0 ? (
              <Card><CardContent className="py-6 text-center text-muted-foreground">暂无申报任务</CardContent></Card>
            ) : (
              tasks.map((t) => (
                <Card key={t.id}>
                  <CardContent className="flex items-center justify-between py-4">
                    <div>
                      <p className="font-medium">{t.taxType}</p>
                      <p className="text-sm text-muted-foreground">{t.period} · 截止 {t.dueDate}</p>
                    </div>
                    <Badge variant={t.status === 'overdue' ? 'destructive' : t.status === 'filed' ? 'success' : 'warning'}>
                      {statusLabel[t.status] ?? t.status}
                    </Badge>
                  </CardContent>
                </Card>
              ))
            )}
          </div>
          <Card>
            <CardHeader>
              <CardTitle className="text-base">申报前自查清单</CardTitle>
            </CardHeader>
            <CardContent>
              <ul className="space-y-2 text-sm">
                {checklist.map((item) => (
                  <li key={item} className="flex gap-2">
                    <span className="text-primary">□</span>
                    {item}
                  </li>
                ))}
              </ul>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  );
}
