'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { ApiClientError, getMemberToken, memberFetch } from '@/lib/api/client';

interface OverviewData {
  opcStatus?: string;
  pendingTasks?: { label: string; href: string; urgent?: boolean }[];
  nextTaxDeadline?: string;
  monthlySummary?: { income: number; expense: number; profit: number };
}

export default function OverviewPage() {
  const [data, setData] = useState<OverviewData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!getMemberToken()) {
      setLoading(false);
      setError('请先登录');
      return;
    }
    memberFetch<OverviewData>('/api/member/compliance/overview')
      .then(setData)
      .catch((err) => {
        setError(err instanceof ApiClientError ? err.message : '加载失败');
      })
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p className="text-muted-foreground">加载中…</p>;
  if (error) {
    return (
      <Alert variant="destructive">
        <AlertDescription>{error}</AlertDescription>
      </Alert>
    );
  }

  const tasks = data?.pendingTasks ?? [];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">合规仪表盘</h1>
        <p className="text-muted-foreground">您的合规服务概览</p>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>OPC 状态</CardDescription>
            <CardTitle className="text-lg">
              <Badge variant={data?.opcStatus === 'active' ? 'success' : 'warning'}>
                {data?.opcStatus === 'active' ? '已激活' : '设立中/未签约'}
              </Badge>
            </CardTitle>
          </CardHeader>
        </Card>
        {data?.nextTaxDeadline && (
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>下次申报截止</CardDescription>
              <CardTitle className="text-lg">{data.nextTaxDeadline}</CardTitle>
            </CardHeader>
          </Card>
        )}
        {data?.monthlySummary && (
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>本月利润（估）</CardDescription>
              <CardTitle className="text-lg">
                ¥{data.monthlySummary.profit.toLocaleString()}
              </CardTitle>
            </CardHeader>
            <CardContent className="text-xs text-muted-foreground">
              收入 ¥{data.monthlySummary.income.toLocaleString()} · 费用 ¥
              {data.monthlySummary.expense.toLocaleString()}
            </CardContent>
          </Card>
        )}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>待办事项</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          {tasks.length === 0 ? (
            <p className="py-4 text-center text-sm text-muted-foreground">暂无待办，继续保持</p>
          ) : (
            tasks.map((task) => (
              <div
                key={task.href + task.label}
                className={`flex items-center justify-between rounded-md border p-3 ${task.urgent ? 'border-destructive/50 bg-destructive/5' : ''}`}
              >
                <span className="text-sm">{task.label}</span>
                <Button asChild size="sm" variant={task.urgent ? 'default' : 'outline'}>
                  <Link href={task.href}>处理</Link>
                </Button>
              </div>
            ))
          )}
        </CardContent>
      </Card>

      <div className="flex flex-wrap gap-3">
        <Button asChild variant="outline">
          <Link href="/user/compliance/income">收入台账</Link>
        </Button>
        <Button asChild variant="outline">
          <Link href="/user/compliance/tax">申报日历</Link>
        </Button>
      </div>
    </div>
  );
}
