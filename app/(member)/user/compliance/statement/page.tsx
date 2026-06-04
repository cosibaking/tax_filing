'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { ApiClientError, memberFetch } from '@/lib/api/client';

interface Statement {
  id: string;
  period: string;
  status: string;
  sentAt?: string;
}

export default function StatementPage() {
  const [items, setItems] = useState<Statement[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    memberFetch<Statement[]>('/api/member/compliance/statements')
      .then(setItems)
      .catch((err) => setError(err instanceof ApiClientError ? err.message : '加载失败'))
      .finally(() => setLoading(false));
  }, []);

  const downloadPdf = (id: string) => {
    window.open(`/api/member/compliance/statements/${id}/pdf`, '_blank');
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">月度对账单</h1>
      {loading && <p className="text-muted-foreground">加载中…</p>}
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      {!loading && items.length === 0 && (
        <Card><CardContent className="py-8 text-center text-muted-foreground">暂无对账单</CardContent></Card>
      )}
      <div className="space-y-3">
        {items.map((s) => (
          <Card key={s.id}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-base">{s.period} 对账单</CardTitle>
              <Badge variant={s.status === 'sent' ? 'success' : 'secondary'}>
                {s.status === 'sent' ? '已发送' : '生成中'}
              </Badge>
            </CardHeader>
            <CardContent className="flex items-center justify-between">
              {s.sentAt && (
                <span className="text-sm text-muted-foreground">
                  发送于 {new Date(s.sentAt).toLocaleDateString('zh-CN')}
                </span>
              )}
              <Button size="sm" variant="outline" onClick={() => downloadPdf(s.id)}>
                下载 PDF
              </Button>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
