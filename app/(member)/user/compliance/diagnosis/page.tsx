'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { ApiClientError, memberFetch } from '@/lib/api/client';

interface DiagnosisRecord {
  id: string;
  recommendedPlan: string;
  createdAt: string;
}

export default function DiagnosisHistoryPage() {
  const [items, setItems] = useState<DiagnosisRecord[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    memberFetch<DiagnosisRecord[]>('/api/member/compliance/diagnosis')
      .then(setItems)
      .catch((err) => {
        setError(err instanceof ApiClientError ? err.message : '加载失败');
      })
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <h1 className="text-2xl font-bold">诊断历史</h1>
        <Button asChild variant="outline">
          <Link href="/diagnosis">新建诊断</Link>
        </Button>
      </div>
      {loading && <p className="text-muted-foreground">加载中…</p>}
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      {!loading && !error && items.length === 0 && (
        <Card>
          <CardContent className="py-8 text-center text-muted-foreground">
            暂无诊断记录，请先完成
            <Link href="/diagnosis" className="text-primary underline">
              免费诊断
            </Link>
          </CardContent>
        </Card>
      )}
      <div className="space-y-3">
        {items.map((d) => (
          <Card key={d.id}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-base">诊断 #{d.id}</CardTitle>
              <Badge>{d.recommendedPlan}</Badge>
            </CardHeader>
            <CardContent className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">
                {new Date(d.createdAt).toLocaleString('zh-CN')}
              </span>
              <Button asChild size="sm" variant="outline">
                <Link href={`/diagnosis/result?id=${d.id}`}>查看结果</Link>
              </Button>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
