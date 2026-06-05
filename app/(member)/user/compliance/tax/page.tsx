'use client';

import { useCallback, useEffect, useMemo, useState } from 'react';
import {
  eachDayOfInterval,
  endOfMonth,
  endOfWeek,
  format,
  isSameMonth,
  startOfMonth,
  startOfWeek,
} from 'date-fns';
import { zhCN } from 'date-fns/locale';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { TAX_FILING_STATUS_LABELS } from '@/lib/api/constants';
import { ApiClientError, getMemberToken, memberFetch } from '@/lib/api/client';
import { withBasePath } from '@/lib/base-path';
import {
  formatNextDueSummary,
  shiftCalendarMonth,
  type TaxCalendarDto,
  type TaxTaskDetailDto,
} from '@/lib/compliance/tax-calendar-shared';

interface ChecklistItem {
  id: number;
  text: string;
  checked: boolean;
}

interface ChecklistResponse {
  items: ChecklistItem[];
  allChecked: boolean;
  source: string | null;
}

const WEEKDAY_LABELS = ['一', '二', '三', '四', '五', '六', '日'];

function formatMoney(n: number): string {
  return `¥${n.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
}

function statusBadgeVariant(status: string): 'destructive' | 'success' | 'warning' | 'secondary' {
  if (status === 'overdue') return 'destructive';
  if (status === 'filed') return 'success';
  if (status === 'pending') return 'warning';
  return 'secondary';
}

const now = new Date();

export default function TaxPage() {
  const [year, setYear] = useState(now.getFullYear());
  const [month, setMonth] = useState(now.getMonth() + 1);
  const [calendar, setCalendar] = useState<TaxCalendarDto | null>(null);
  const [checklist, setChecklist] = useState<ChecklistResponse | null>(null);
  const [selectedTaskId, setSelectedTaskId] = useState<string | null>(null);
  const [taskDetail, setTaskDetail] = useState<TaxTaskDetailDto | null>(null);
  const [detailLoading, setDetailLoading] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(() => {
    setLoading(true);
    setError(null);
    const params = new URLSearchParams({ year: String(year), month: String(month) });
    Promise.all([
      memberFetch<TaxCalendarDto>(`/api/member/compliance/tax/calendar?${params}`),
      memberFetch<ChecklistResponse>('/api/member/compliance/tax/checklist').catch(() => null),
    ])
      .then(([cal, chk]) => {
        setCalendar(cal);
        setChecklist(chk);
      })
      .catch((err) => {
        setError(err instanceof ApiClientError ? err.message : '加载失败');
      })
      .finally(() => setLoading(false));
  }, [year, month]);

  useEffect(() => {
    load();
  }, [load]);

  useEffect(() => {
    if (!selectedTaskId) {
      setTaskDetail(null);
      return;
    }
    setDetailLoading(true);
    memberFetch<TaxTaskDetailDto>(`/api/member/compliance/tax/tasks/${selectedTaskId}`)
      .then(setTaskDetail)
      .catch(() => setTaskDetail(null))
      .finally(() => setDetailLoading(false));
  }, [selectedTaskId]);

  const monthAnchor = useMemo(() => new Date(year, month - 1, 1), [year, month]);

  const calendarCells = useMemo(() => {
    const start = startOfWeek(startOfMonth(monthAnchor), { weekStartsOn: 1 });
    const end = endOfWeek(endOfMonth(monthAnchor), { weekStartsOn: 1 });
    return eachDayOfInterval({ start, end });
  }, [monthAnchor]);

  const dueDaySet = useMemo(() => new Set(calendar?.dueDays ?? []), [calendar?.dueDays]);

  const changeMonth = (delta: number) => {
    setSelectedTaskId(null);
    const next = shiftCalendarMonth(year, month, delta);
    setYear(next.year);
    setMonth(next.month);
  };

  const downloadReceipt = async (fileId: string) => {
    const token = getMemberToken();
    const res = await fetch(withBasePath(`/api/upload/${fileId}`), {
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    });
    if (!res.ok) {
      setError('回执下载失败');
      return;
    }
    const blob = await res.blob();
    const url = URL.createObjectURL(blob);
    window.open(url, '_blank');
    setTimeout(() => URL.revokeObjectURL(url), 60_000);
  };

  const nextDueText = calendar?.nextDue ? formatNextDueSummary(calendar.nextDue) : null;

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">申报日历</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          查看增值税、附加税、企业所得税等申报截止日与任务状态（US-E5-08）
        </p>
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <Card>
        <CardContent className="flex flex-col gap-3 pt-6 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex items-center gap-2">
            <Button type="button" variant="outline" size="icon" onClick={() => changeMonth(-1)} aria-label="上个月">
              <ChevronLeft className="h-4 w-4" />
            </Button>
            <span className="min-w-[7rem] text-center font-semibold">
              {format(monthAnchor, 'yyyy年M月', { locale: zhCN })}
            </span>
            <Button type="button" variant="outline" size="icon" onClick={() => changeMonth(1)} aria-label="下个月">
              <ChevronRight className="h-4 w-4" />
            </Button>
          </div>
          {nextDueText && (
            <p className="text-sm">
              <span className="text-muted-foreground">下次截止：</span>
              <span className="font-medium text-primary">{nextDueText}</span>
            </p>
          )}
        </CardContent>
      </Card>

      {loading ? (
        <p className="text-muted-foreground">加载中…</p>
      ) : (
        <>
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-base">月历</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-7 gap-1 text-center text-xs text-muted-foreground">
                {WEEKDAY_LABELS.map((d) => (
                  <div key={d} className="py-1 font-medium">
                    {d}
                  </div>
                ))}
              </div>
              <div className="grid grid-cols-7 gap-1">
                {calendarCells.map((day) => {
                  const inMonth = isSameMonth(day, monthAnchor);
                  const dayNum = day.getDate();
                  const isDue = inMonth && dueDaySet.has(dayNum);
                  return (
                    <div
                      key={day.toISOString()}
                      className={`relative flex min-h-[2.5rem] items-start justify-center rounded-md p-1 text-sm ${
                        inMonth ? 'text-foreground' : 'text-muted-foreground/40'
                      } ${isDue ? 'bg-destructive/10 font-semibold' : ''}`}
                    >
                      {dayNum}
                      {isDue && (
                        <span
                          className="absolute bottom-1 left-1/2 h-1.5 w-1.5 -translate-x-1/2 rounded-full bg-destructive"
                          title="申报截止日"
                        />
                      )}
                    </div>
                  );
                })}
              </div>
              <p className="mt-3 text-xs text-muted-foreground">
                <span className="mr-1 inline-block h-1.5 w-1.5 rounded-full bg-destructive align-middle" />
                红点为申报截止日（通常为每月 15 日）
              </p>
            </CardContent>
          </Card>

          <div className="grid gap-6 lg:grid-cols-2">
            <div className="space-y-3">
              <h2 className="font-semibold">本月申报任务</h2>
              {!calendar?.tasks.length ? (
                <Card>
                  <CardContent className="py-6 text-center text-muted-foreground">
                    本月暂无申报任务，请确认收入台账已录入
                  </CardContent>
                </Card>
              ) : (
                calendar.tasks.map((t) => (
                  <Card
                    key={t.id}
                    className={selectedTaskId === t.id ? 'ring-2 ring-primary' : undefined}
                  >
                    <CardContent className="flex flex-col gap-3 py-4 sm:flex-row sm:items-center sm:justify-between">
                      <div>
                        <p className="font-medium">{t.taxTypeLabel}</p>
                        <p className="text-sm text-muted-foreground">
                          {t.period} · 截止 {t.dueDate}
                        </p>
                        <p className="mt-1 text-sm">
                          {t.status === 'filed' && t.filedAmount != null
                            ? `已申报 ${formatMoney(t.filedAmount)}`
                            : `预估 ${formatMoney(t.calculatedAmount)}`}
                        </p>
                      </div>
                      <div className="flex flex-wrap items-center gap-2">
                        <Badge variant={statusBadgeVariant(t.status)}>
                          {TAX_FILING_STATUS_LABELS[t.status] ?? t.status}
                        </Badge>
                        <Button
                          type="button"
                          size="sm"
                          variant="outline"
                          onClick={() => setSelectedTaskId(t.id)}
                        >
                          详情
                        </Button>
                        {t.status === 'filed' && t.receiptFileId && (
                          <Button
                            type="button"
                            size="sm"
                            variant="ghost"
                            onClick={() => downloadReceipt(t.receiptFileId!)}
                          >
                            回执
                          </Button>
                        )}
                      </div>
                    </CardContent>
                  </Card>
                ))
              )}
            </div>

            <Card>
              <CardHeader>
                <CardTitle className="text-base">申报前自查清单</CardTitle>
                {checklist?.source && (
                  <p className="text-xs text-muted-foreground">
                    状态来自最近申报：{checklist.source}
                  </p>
                )}
              </CardHeader>
              <CardContent>
                <ul className="space-y-2 text-sm">
                  {(checklist?.items ?? []).map((item) => (
                    <li
                      key={item.id}
                      className={`flex gap-2 rounded-md px-2 py-1.5 ${
                        item.checked ? 'text-muted-foreground' : 'bg-amber-50 text-amber-950 dark:bg-amber-950/30 dark:text-amber-100'
                      }`}
                    >
                      <span className="shrink-0">{item.checked ? '☑' : '□'}</span>
                      <span>{item.text}</span>
                    </li>
                  ))}
                </ul>
                {checklist && !checklist.allChecked && (
                  <p className="mt-3 text-xs text-amber-700 dark:text-amber-300">
                    未完成项已高亮；顾问完成全部 9 项自查并标记申报后，此处将显示为已勾选。
                  </p>
                )}
              </CardContent>
            </Card>
          </div>

          {selectedTaskId && (
            <Card>
              <CardHeader className="flex flex-row items-center justify-between">
                <CardTitle className="text-base">任务详情</CardTitle>
                <Button type="button" variant="ghost" size="sm" onClick={() => setSelectedTaskId(null)}>
                  关闭
                </Button>
              </CardHeader>
              <CardContent className="space-y-4">
                {detailLoading || !taskDetail ? (
                  <p className="text-muted-foreground">{detailLoading ? '加载中…' : '无法加载详情'}</p>
                ) : (
                  <>
                    <div className="flex flex-wrap items-center gap-2">
                      <span className="font-medium">{taskDetail.taxTypeLabel}</span>
                      <Badge variant={statusBadgeVariant(taskDetail.status)}>
                        {TAX_FILING_STATUS_LABELS[taskDetail.status] ?? taskDetail.status}
                      </Badge>
                    </div>
                    <dl className="grid gap-2 text-sm sm:grid-cols-2">
                      <div>
                        <dt className="text-muted-foreground">所属期间</dt>
                        <dd>{taskDetail.period}</dd>
                      </div>
                      <div>
                        <dt className="text-muted-foreground">截止日期</dt>
                        <dd>{taskDetail.dueDate}</dd>
                      </div>
                      <div>
                        <dt className="text-muted-foreground">系统测算税额</dt>
                        <dd>{formatMoney(taskDetail.calculatedAmount)}</dd>
                      </div>
                      {taskDetail.filedAmount != null && (
                        <div>
                          <dt className="text-muted-foreground">实际申报额</dt>
                          <dd>{formatMoney(taskDetail.filedAmount)}</dd>
                        </div>
                      )}
                    </dl>
                    <div>
                      <h3 className="mb-2 text-sm font-semibold">计算明细</h3>
                      <ul className="space-y-1 text-sm">
                        {taskDetail.breakdown.map((line) => (
                          <li key={line.label} className="flex justify-between gap-4 border-b py-1">
                            <span className="text-muted-foreground">{line.label}</span>
                            <span>{formatMoney(line.amount)}</span>
                          </li>
                        ))}
                      </ul>
                    </div>
                    {taskDetail.receiptFileId && (
                      <Button type="button" variant="outline" onClick={() => downloadReceipt(taskDetail.receiptFileId!)}>
                        下载申报回执
                      </Button>
                    )}
                  </>
                )}
              </CardContent>
            </Card>
          )}
        </>
      )}
    </div>
  );
}
