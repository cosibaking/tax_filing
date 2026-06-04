'use client';

import { useCallback, useEffect, useMemo, useState } from 'react';
import { CheckCircle2, Download, FileText } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { ApiClientError, getMemberToken, memberFetch } from '@/lib/api/client';
import {
  formatStatementMoney,
  statementStatusBadgeVariant,
  type MemberStatementDto,
} from '@/lib/compliance/statement-shared';

function SummaryRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex justify-between gap-4 py-2 text-sm">
      <span className="text-muted-foreground">{label}</span>
      <span className="font-medium tabular-nums">{value}</span>
    </div>
  );
}

export default function StatementPage() {
  const [items, setItems] = useState<MemberStatementDto[]>([]);
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [downloading, setDownloading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(() => {
    setLoading(true);
    setError(null);
    memberFetch<MemberStatementDto[]>('/api/member/compliance/statements')
      .then((list) => {
        setItems(list);
        setSelectedId((prev) => {
          if (prev && list.some((s) => s.id === prev)) return prev;
          return list[0]?.id ?? null;
        });
      })
      .catch((err) => setError(err instanceof ApiClientError ? err.message : '加载失败'))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const selected = useMemo(
    () => items.find((s) => s.id === selectedId) ?? null,
    [items, selectedId],
  );

  const downloadPdf = async (id: string) => {
    setDownloading(true);
    setError(null);
    try {
      const token = getMemberToken();
      const res = await fetch(`/api/member/compliance/statements/${id}/pdf`, {
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      });
      if (!res.ok) {
        const json = await res.json().catch(() => null);
        throw new Error((json as { message?: string })?.message ?? 'PDF 下载失败');
      }
      const blob = await res.blob();
      const url = URL.createObjectURL(blob);
      window.open(url, '_blank');
      setTimeout(() => URL.revokeObjectURL(url), 60_000);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'PDF 下载失败');
    } finally {
      setDownloading(false);
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">月度对账单</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          查看历史对账单与税务服务完成状态；对账单于次月 16–20 日生成并推送通知（F-61）
        </p>
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {loading && <p className="text-muted-foreground">加载中…</p>}

      {!loading && items.length === 0 && (
        <Card>
          <CardContent className="flex flex-col items-center gap-3 py-12 text-center">
            <FileText className="h-10 w-10 text-muted-foreground" />
            <p className="font-medium">暂无对账单</p>
            <p className="max-w-md text-sm text-muted-foreground">
              完成上月收入、费用录入与申报后，系统将在次月 16–20 日自动生成对账单并发送通知。
            </p>
          </CardContent>
        </Card>
      )}

      {!loading && items.length > 0 && (
        <>
          <Card>
            <CardHeader className="pb-3">
              <CardTitle className="text-base">对账单列表</CardTitle>
              <CardDescription>点击行查看当月服务完成摘要</CardDescription>
            </CardHeader>
            <CardContent className="p-0">
              <div className="hidden overflow-x-auto md:block">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b bg-muted/40 text-left">
                      <th className="p-3 font-medium">账期</th>
                      <th className="p-3 font-medium text-right">收入</th>
                      <th className="p-3 font-medium text-right">成本</th>
                      <th className="p-3 font-medium text-right">利润</th>
                      <th className="p-3 font-medium text-right">预缴税额</th>
                      <th className="p-3 font-medium">状态</th>
                      <th className="p-3 font-medium">操作</th>
                    </tr>
                  </thead>
                  <tbody>
                    {items.map((s) => {
                      const active = s.id === selectedId;
                      return (
                        <tr
                          key={s.id}
                          className={`cursor-pointer border-b transition-colors hover:bg-muted/30 ${active ? 'bg-muted/50' : ''}`}
                          onClick={() => setSelectedId(s.id)}
                        >
                          <td className="p-3 font-medium">{s.period}</td>
                          <td className="p-3 text-right tabular-nums">
                            {formatStatementMoney(s.summary.revenue)}
                          </td>
                          <td className="p-3 text-right tabular-nums">
                            {formatStatementMoney(s.summary.cost)}
                          </td>
                          <td className="p-3 text-right tabular-nums">
                            {formatStatementMoney(s.summary.profit)}
                          </td>
                          <td className="p-3 text-right tabular-nums">
                            {formatStatementMoney(s.summary.taxTotal)}
                          </td>
                          <td className="p-3">
                            <Badge
                              variant={statementStatusBadgeVariant(s.status, s.serviceComplete)}
                            >
                              {s.serviceComplete ? '已完成' : s.statusLabel}
                            </Badge>
                          </td>
                          <td className="p-3">
                            {s.hasPdf && (
                              <Button
                                size="sm"
                                variant="outline"
                                onClick={(e) => {
                                  e.stopPropagation();
                                  void downloadPdf(s.id);
                                }}
                                disabled={downloading}
                              >
                                <Download className="mr-1 h-3.5 w-3.5" />
                                PDF
                              </Button>
                            )}
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>

              <div className="divide-y md:hidden">
                {items.map((s) => {
                  const active = s.id === selectedId;
                  return (
                    <button
                      key={s.id}
                      type="button"
                      className={`w-full px-4 py-3 text-left transition-colors ${active ? 'bg-muted/50' : ''}`}
                      onClick={() => setSelectedId(s.id)}
                    >
                      <div className="flex items-center justify-between gap-2">
                        <span className="font-medium">{s.period}</span>
                        <Badge
                          variant={statementStatusBadgeVariant(s.status, s.serviceComplete)}
                        >
                          {s.serviceComplete ? '已完成' : s.statusLabel}
                        </Badge>
                      </div>
                      <p className="mt-1 text-xs text-muted-foreground">
                        利润 {formatStatementMoney(s.summary.profit)} · 税额{' '}
                        {formatStatementMoney(s.summary.taxTotal)}
                      </p>
                    </button>
                  );
                })}
              </div>
            </CardContent>
          </Card>

          {selected && (
            <Card className="border-primary/20">
              <CardHeader>
                <div className="flex flex-wrap items-start justify-between gap-3">
                  <div>
                    <CardTitle className="flex items-center gap-2 text-lg">
                      {selected.year}年{selected.month}月 · 税务服务完成报告
                      {selected.serviceComplete && (
                        <CheckCircle2 className="h-5 w-5 text-emerald-600" aria-hidden />
                      )}
                    </CardTitle>
                    <CardDescription className="mt-1">
                      {selected.sentAt
                        ? `已于 ${new Date(selected.sentAt).toLocaleString('zh-CN')} 发送通知`
                        : '对账单已生成，通知发送后将在此显示时间'}
                    </CardDescription>
                  </div>
                  {selected.hasPdf && (
                    <Button
                      variant="default"
                      size="sm"
                      disabled={downloading}
                      onClick={() => void downloadPdf(selected.id)}
                    >
                      <Download className="mr-1.5 h-4 w-4" />
                      {downloading ? '下载中…' : '下载 PDF'}
                    </Button>
                  )}
                </div>
              </CardHeader>
              <CardContent className="divide-y rounded-lg border bg-muted/20">
                <SummaryRow label="本月收入" value={formatStatementMoney(selected.summary.revenue)} />
                <SummaryRow label="本月成本" value={formatStatementMoney(selected.summary.cost)} />
                <SummaryRow label="本月利润" value={formatStatementMoney(selected.summary.profit)} />
                <SummaryRow
                  label="本月预缴税额"
                  value={formatStatementMoney(selected.summary.taxTotal)}
                />
                <SummaryRow
                  label="累计年度利润"
                  value={formatStatementMoney(selected.summary.cumulativeProfit)}
                />
                <SummaryRow label="申报状态" value={selected.summary.filingStatus} />
              </CardContent>
            </Card>
          )}
        </>
      )}
    </div>
  );
}
