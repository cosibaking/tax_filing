'use client';

import { useCallback, useEffect, useState } from 'react';
import Link from 'next/link';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Select } from '@/components/ui/select';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { LEDGER_ACCOUNT_LABELS, type ProfitView } from '@/lib/api/constants';
import { ApiClientError, memberFetch } from '@/lib/api/client';

interface ProfitPeriodRow {
  key: string;
  label: string;
  month?: number;
  quarter?: number;
  revenue: number;
  cost: number;
  profit: number;
  cumulativeProfit: number;
}

interface ProfitReport {
  year: number;
  view: ProfitView;
  rows: ProfitPeriodRow[];
  yearTotal: { revenue: number; cost: number; profit: number };
  cumulativeProfit: number;
}

interface LedgerVoucher {
  id: string;
  period: string;
  debitAccount: string;
  creditAccount: string;
  amount: string | number;
  refType: string;
}

interface MonthDetailResponse {
  year: number;
  month: number;
  detail?: {
    revenue: number;
    cost: number;
    profit: number;
    cumulativeProfit: number;
  };
  vouchers?: LedgerVoucher[];
}

function formatMoney(n: number): string {
  return `¥${n.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
}

function accountLabel(code: string): string {
  return LEDGER_ACCOUNT_LABELS[code] ?? code;
}

function refTypeLabel(refType: string): string {
  if (refType === 'income') return '收入';
  if (refType === 'expense') return '费用';
  return refType;
}

const CURRENT_YEAR = new Date().getFullYear();
const YEAR_OPTIONS = [CURRENT_YEAR, CURRENT_YEAR - 1, CURRENT_YEAR - 2];

export default function LedgerPage() {
  const [year, setYear] = useState(CURRENT_YEAR);
  const [view, setView] = useState<ProfitView>('monthly');
  const [report, setReport] = useState<ProfitReport | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [detailMonth, setDetailMonth] = useState<number | null>(null);
  const [vouchers, setVouchers] = useState<LedgerVoucher[]>([]);
  const [detailLoading, setDetailLoading] = useState(false);

  const loadReport = useCallback(() => {
    setLoading(true);
    setDetailMonth(null);
    const params = new URLSearchParams({ year: String(year), view });
    memberFetch<ProfitReport>(`/api/member/compliance/ledger/profit?${params}`)
      .then((d) => {
        setReport(d);
        setError(null);
      })
      .catch((err) => setError(err instanceof ApiClientError ? err.message : '加载失败'))
      .finally(() => setLoading(false));
  }, [year, view]);

  useEffect(() => {
    loadReport();
  }, [loadReport]);

  const loadMonthDetail = (month: number) => {
    setDetailMonth(month);
    setDetailLoading(true);
    memberFetch<MonthDetailResponse>(
      `/api/member/compliance/ledger/profit?year=${year}&month=${month}`,
    )
      .then((d) => setVouchers(d.vouchers ?? []))
      .catch(() => setVouchers([]))
      .finally(() => setDetailLoading(false));
  };

  const showCumulative = view === 'monthly';

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold">利润表</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            收入为不含税口径（含税收入÷1.01），成本为费用台账合计；数据来源于收入与费用台账自动汇总
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm" asChild>
            <Link href="/user/compliance/income">收入台账</Link>
          </Button>
          <Button variant="outline" size="sm" asChild>
            <Link href="/user/compliance/expense">费用台账</Link>
          </Button>
        </div>
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <Card>
        <CardContent className="flex flex-col gap-4 pt-6 sm:flex-row sm:items-end">
          <div className="w-full sm:w-40">
            <Label htmlFor="profit-year">年度</Label>
            <Select
              id="profit-year"
              className="mt-1"
              value={String(year)}
              onChange={(e) => setYear(parseInt(e.target.value, 10))}
              disabled={loading}
            >
              {YEAR_OPTIONS.map((y) => (
                <option key={y} value={y}>
                  {y} 年
                </option>
              ))}
            </Select>
          </div>
        </CardContent>
      </Card>

      {report && !loading && (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>年度收入（不含税）</CardDescription>
              <CardTitle className="text-lg">{formatMoney(report.yearTotal.revenue)}</CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>年度成本费用</CardDescription>
              <CardTitle className="text-lg">{formatMoney(report.yearTotal.cost)}</CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>年度利润</CardDescription>
              <CardTitle className="text-lg text-primary">{formatMoney(report.yearTotal.profit)}</CardTitle>
            </CardHeader>
          </Card>
          <Card>
            <CardHeader className="pb-2">
              <CardDescription>累计年度利润</CardDescription>
              <CardTitle className="text-lg">{formatMoney(report.cumulativeProfit)}</CardTitle>
            </CardHeader>
          </Card>
        </div>
      )}

      <Tabs value={view} onValueChange={(v) => setView(v as ProfitView)}>
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="monthly">月度</TabsTrigger>
          <TabsTrigger value="quarterly">季度</TabsTrigger>
          <TabsTrigger value="yearly">年度</TabsTrigger>
        </TabsList>
        <TabsContent value={view} className="mt-4 space-y-4">
          {loading ? (
            <p className="text-muted-foreground">加载中…</p>
          ) : !report || report.rows.length === 0 ? (
            <Card>
              <CardContent className="py-8 text-center text-muted-foreground">
                {year} 年暂无利润数据，请先在
                <Link href="/user/compliance/income" className="mx-1 text-primary underline">
                  收入台账
                </Link>
                与
                <Link href="/user/compliance/expense" className="mx-1 text-primary underline">
                  费用台账
                </Link>
                录入数据
              </CardContent>
            </Card>
          ) : (
            <>
              <div className="space-y-2 md:hidden">
                {report.rows.map((row) => (
                  <Card key={row.key}>
                    <CardHeader className="pb-2">
                      <CardTitle className="text-base">{row.label}</CardTitle>
                    </CardHeader>
                    <CardContent className="grid grid-cols-2 gap-3 text-sm">
                      <div>
                        <p className="text-muted-foreground">收入（不含税）</p>
                        <p>{formatMoney(row.revenue)}</p>
                      </div>
                      <div>
                        <p className="text-muted-foreground">成本费用</p>
                        <p>{formatMoney(row.cost)}</p>
                      </div>
                      <div>
                        <p className="text-muted-foreground">利润</p>
                        <p className="font-semibold text-primary">{formatMoney(row.profit)}</p>
                      </div>
                      {showCumulative && (
                        <div>
                          <p className="text-muted-foreground">累计利润</p>
                          <p>{formatMoney(row.cumulativeProfit)}</p>
                        </div>
                      )}
                      {view === 'monthly' && row.month && (
                        <div className="col-span-2">
                          <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            onClick={() => loadMonthDetail(row.month!)}
                          >
                            查看分录
                          </Button>
                        </div>
                      )}
                    </CardContent>
                  </Card>
                ))}
              </div>
              <div className="hidden overflow-x-auto md:block">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b text-left text-muted-foreground">
                      <th className="p-2">期间</th>
                      <th className="p-2 text-right">收入（不含税）</th>
                      <th className="p-2 text-right">成本费用</th>
                      <th className="p-2 text-right">利润</th>
                      {showCumulative && <th className="p-2 text-right">累计年度利润</th>}
                      {view === 'monthly' && <th className="p-2" />}
                    </tr>
                  </thead>
                  <tbody>
                    {report.rows.map((row) => (
                      <tr key={row.key} className="border-b">
                        <td className="p-2 font-medium">{row.label}</td>
                        <td className="p-2 text-right">{formatMoney(row.revenue)}</td>
                        <td className="p-2 text-right">{formatMoney(row.cost)}</td>
                        <td className="p-2 text-right font-medium text-primary">
                          {formatMoney(row.profit)}
                        </td>
                        {showCumulative && (
                          <td className="p-2 text-right">{formatMoney(row.cumulativeProfit)}</td>
                        )}
                        {view === 'monthly' && row.month && (
                          <td className="p-2 text-right">
                            <Button
                              type="button"
                              variant="ghost"
                              size="sm"
                              onClick={() => loadMonthDetail(row.month!)}
                            >
                              分录
                            </Button>
                          </td>
                        )}
                      </tr>
                    ))}
                  </tbody>
                  <tfoot>
                    <tr className="border-t bg-muted/40 font-medium">
                      <td className="p-2">合计</td>
                      <td className="p-2 text-right">{formatMoney(report.yearTotal.revenue)}</td>
                      <td className="p-2 text-right">{formatMoney(report.yearTotal.cost)}</td>
                      <td className="p-2 text-right text-primary">
                        {formatMoney(report.yearTotal.profit)}
                      </td>
                      {showCumulative && (
                        <td className="p-2 text-right">{formatMoney(report.cumulativeProfit)}</td>
                      )}
                      {view === 'monthly' && <td className="p-2" />}
                    </tr>
                  </tfoot>
                </table>
              </div>
            </>
          )}
        </TabsContent>
      </Tabs>

      {detailMonth != null && (
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-base">
                {year}年{detailMonth}月 会计分录
              </CardTitle>
              <CardDescription>收入/费用确认后自动生成的账套凭证（F-45）</CardDescription>
            </div>
            <Button type="button" variant="ghost" size="sm" onClick={() => setDetailMonth(null)}>
              关闭
            </Button>
          </CardHeader>
          <CardContent>
            {detailLoading ? (
              <p className="text-muted-foreground">加载中…</p>
            ) : vouchers.length === 0 ? (
              <p className="text-sm text-muted-foreground">该月暂无会计分录</p>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b text-left text-muted-foreground">
                      <th className="p-2">类型</th>
                      <th className="p-2">借方</th>
                      <th className="p-2">贷方</th>
                      <th className="p-2 text-right">金额</th>
                    </tr>
                  </thead>
                  <tbody>
                    {vouchers.map((v) => (
                      <tr key={v.id} className="border-b">
                        <td className="p-2">{refTypeLabel(v.refType)}</td>
                        <td className="p-2">
                          {accountLabel(v.debitAccount)}
                          <span className="ml-1 text-xs text-muted-foreground">({v.debitAccount})</span>
                        </td>
                        <td className="p-2">
                          {accountLabel(v.creditAccount)}
                          <span className="ml-1 text-xs text-muted-foreground">({v.creditAccount})</span>
                        </td>
                        <td className="p-2 text-right">
                          {formatMoney(typeof v.amount === 'number' ? v.amount : parseFloat(v.amount))}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  );
}
