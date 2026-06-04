'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { ApiClientError, memberFetch } from '@/lib/api/client';

interface ProfitRow {
  month: string;
  income: number;
  expense: number;
  profit: number;
}

export default function LedgerPage() {
  const [rows, setRows] = useState<ProfitRow[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    memberFetch<{ months: ProfitRow[] }>('/api/member/compliance/ledger/profit')
      .then((d) => setRows(d.months ?? []))
      .catch((err) => setError(err instanceof ApiClientError ? err.message : '加载失败'))
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">账套与利润表</h1>
      {loading && <p className="text-muted-foreground">加载中…</p>}
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      {!loading && !error && (
        <>
          <div className="grid gap-4 sm:hidden">
            {rows.map((r) => (
              <Card key={r.month}>
                <CardHeader className="pb-2">
                  <CardTitle className="text-base">{r.month}</CardTitle>
                </CardHeader>
                <CardContent className="grid grid-cols-3 gap-2 text-sm">
                  <div>
                    <p className="text-muted-foreground">收入</p>
                    <p>¥{r.income.toLocaleString()}</p>
                  </div>
                  <div>
                    <p className="text-muted-foreground">费用</p>
                    <p>¥{r.expense.toLocaleString()}</p>
                  </div>
                  <div>
                    <p className="text-muted-foreground">利润</p>
                    <p className="font-semibold text-primary">¥{r.profit.toLocaleString()}</p>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
          <div className="hidden overflow-x-auto md:block">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b text-muted-foreground">
                  <th className="p-2 text-left">月份</th>
                  <th className="p-2 text-right">收入</th>
                  <th className="p-2 text-right">费用</th>
                  <th className="p-2 text-right">利润</th>
                </tr>
              </thead>
              <tbody>
                {rows.map((r) => (
                  <tr key={r.month} className="border-b">
                    <td className="p-2">{r.month}</td>
                    <td className="p-2 text-right">¥{r.income.toLocaleString()}</td>
                    <td className="p-2 text-right">¥{r.expense.toLocaleString()}</td>
                    <td className="p-2 text-right font-medium">¥{r.profit.toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {rows.length === 0 && (
            <Card><CardContent className="py-8 text-center text-muted-foreground">暂无利润数据</CardContent></Card>
          )}
        </>
      )}
    </div>
  );
}
