'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { ApiClientError, memberFetch } from '@/lib/api/client';

interface Notice {
  id: string;
  title: string;
  content: string;
  type: string;
  read: boolean;
  createdAt: string;
}

export default function NotificationPage() {
  const [items, setItems] = useState<Notice[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    memberFetch<Notice[]>('/api/member/notices')
      .then(setItems)
      .catch((err) => {
        setError(err instanceof ApiClientError ? err.message : '加载失败');
      })
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">通知中心</h1>
      {loading && <p className="text-muted-foreground">加载中…</p>}
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      {!loading && !error && items.length === 0 && (
        <Card>
          <CardContent className="py-8 text-center text-muted-foreground">暂无通知</CardContent>
        </Card>
      )}
      <div className="space-y-3">
        {items.map((n) => (
          <Card key={n.id} className={!n.read ? 'border-primary/30' : ''}>
            <CardHeader className="flex flex-row items-start justify-between gap-2 pb-2">
              <CardTitle className="text-base">{n.title}</CardTitle>
              {!n.read && <Badge>未读</Badge>}
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">{n.content}</p>
              <p className="mt-2 text-xs text-muted-foreground">
                {new Date(n.createdAt).toLocaleString('zh-CN')}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
