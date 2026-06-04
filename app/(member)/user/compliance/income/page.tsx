'use client';

import { useEffect, useMemo, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select } from '@/components/ui/select';
import { Alert, AlertDescription } from '@/components/ui/alert';
import {
  DEFAULT_INCOME_SORT,
  INCOME_CATEGORIES,
  INCOME_SORT_OPTIONS,
  PLATFORMS,
  type IncomeSort,
} from '@/lib/api/constants';
import { ApiClientError, memberFetch } from '@/lib/api/client';

interface ApiIncomeEntry {
  id: string;
  platform: string;
  category: string;
  grossAmount: string | number;
  platformFee: string | number;
  netAmount: string | number;
  occurredAt: string;
  note?: string;
}

interface IncomeRow {
  id: string;
  date: string;
  platform: string;
  platformLabel: string;
  category: string;
  categoryLabel: string;
  grossAmount: number;
  platformFee: number;
  netAmount: number;
  note?: string;
}

const PAGE_SIZE = 10;

const EMPTY_FORM = {
  date: '',
  platform: '',
  category: '',
  grossAmount: '',
  platformFee: '',
  note: '',
};

interface IncomeFilters {
  dateFrom: string;
  dateTo: string;
  platform: string;
  category: string;
}

interface IncomeSummary {
  grossAmount: number;
  platformFee: number;
  netAmount: number;
}

interface IncomeListResponse {
  items: ApiIncomeEntry[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
  sort?: IncomeSort;
  summary?: IncomeSummary;
}

function formatMoney(n: number): string {
  return `¥${n.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
}

function toNumber(v: string | number): number {
  return typeof v === 'number' ? v : parseFloat(v);
}

function formatDate(iso: string): string {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return iso;
  return d.toLocaleDateString('zh-CN');
}

function labelFor<T extends { value: string; label: string }>(options: readonly T[], value: string) {
  return options.find((o) => o.value === value)?.label ?? value;
}

function mapEntry(raw: ApiIncomeEntry): IncomeRow {
  return {
    id: raw.id,
    date: formatDate(raw.occurredAt),
    platform: raw.platform,
    platformLabel: labelFor(PLATFORMS, raw.platform),
    category: raw.category,
    categoryLabel: labelFor(INCOME_CATEGORIES, raw.category),
    grossAmount: toNumber(raw.grossAmount),
    platformFee: toNumber(raw.platformFee),
    netAmount: toNumber(raw.netAmount),
    note: raw.note,
  };
}

export default function IncomePage() {
  const [items, setItems] = useState<IncomeRow[]>([]);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState(EMPTY_FORM);
  const [saving, setSaving] = useState(false);
  const [sort, setSort] = useState<IncomeSort>(DEFAULT_INCOME_SORT);
  const [filters, setFilters] = useState<IncomeFilters>({
    dateFrom: '',
    dateTo: '',
    platform: '',
    category: '',
  });
  const [summary, setSummary] = useState<IncomeSummary>({
    grossAmount: 0,
    platformFee: 0,
    netAmount: 0,
  });

  const netPreview = useMemo(() => {
    const gross = parseFloat(form.grossAmount);
    const fee = parseFloat(form.platformFee || '0');
    if (Number.isNaN(gross) || Number.isNaN(fee)) return null;
    return gross - fee;
  }, [form.grossAmount, form.platformFee]);

  const load = (
    targetPage: number,
    targetSort: IncomeSort = sort,
    targetFilters: IncomeFilters = filters,
  ) => {
    const hasDateFrom = Boolean(targetFilters.dateFrom);
    const hasDateTo = Boolean(targetFilters.dateTo);
    if (hasDateFrom !== hasDateTo) {
      setError('请同时选择开始与结束日期');
      setLoading(false);
      return;
    }
    if (hasDateFrom && hasDateTo && targetFilters.dateFrom > targetFilters.dateTo) {
      setError('开始日期不能晚于结束日期');
      setLoading(false);
      return;
    }
    setLoading(true);
    const params = new URLSearchParams({
      page: String(targetPage),
      pageSize: String(PAGE_SIZE),
      sort: targetSort,
    });
    if (hasDateFrom && hasDateTo) {
      params.set('dateFrom', targetFilters.dateFrom);
      params.set('dateTo', targetFilters.dateTo);
    }
    if (targetFilters.platform) params.set('platform', targetFilters.platform);
    if (targetFilters.category) params.set('category', targetFilters.category);

    memberFetch<IncomeListResponse>(`/api/member/compliance/income?${params}`)
      .then((d) => {
        setItems((d.items ?? []).map(mapEntry));
        setTotal(d.total ?? 0);
        setTotalPages(d.totalPages ?? 0);
        setPage(d.page ?? targetPage);
        if (d.sort) setSort(d.sort);
        setSummary(
          d.summary ?? { grossAmount: 0, platformFee: 0, netAmount: 0 },
        );
        setError(null);
      })
      .catch((err) => setError(err instanceof ApiClientError ? err.message : '加载失败'))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    load(page, sort, filters);
  }, [page, sort, filters.dateFrom, filters.dateTo, filters.platform, filters.category]);

  const resetPageIfNeeded = () => {
    if (page !== 1) setPage(1);
  };

  const handleSortChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const next = e.target.value as IncomeSort;
    if (next === sort) return;
    setSort(next);
    resetPageIfNeeded();
  };

  const updateFilter = <K extends keyof IncomeFilters>(key: K, value: IncomeFilters[K]) => {
    setFilters((prev) => ({ ...prev, [key]: value }));
    resetPageIfNeeded();
  };

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    const grossAmount = parseFloat(form.grossAmount);
    const platformFee = parseFloat(form.platformFee || '0');
    if (Number.isNaN(grossAmount) || grossAmount < 0) {
      setError('请填写有效的含税收入');
      return;
    }
    if (Number.isNaN(platformFee) || platformFee < 0) {
      setError('请填写有效的平台服务费');
      return;
    }
    if (platformFee > grossAmount) {
      setError('平台服务费不能大于含税收入');
      return;
    }

    setSaving(true);
    try {
      await memberFetch('/api/member/compliance/income', {
        method: 'POST',
        body: JSON.stringify({
          platform: form.platform,
          category: form.category,
          grossAmount,
          platformFee,
          occurredAt: form.date,
          note: form.note || undefined,
        }),
      });
      setShowForm(false);
      setForm(EMPTY_FORM);
      if (page === 1) load(1);
      else setPage(1);
    } catch (err) {
      setError(err instanceof ApiClientError ? err.message : '保存失败');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <h1 className="text-2xl font-bold">收入台账</h1>
        <Button onClick={() => setShowForm(!showForm)}>{showForm ? '取消' : '记一笔收入'}</Button>
      </div>
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      {showForm && (
        <Card>
          <CardHeader>
            <CardTitle>新增收入</CardTitle>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleAdd} className="grid gap-4 sm:grid-cols-2">
              <div>
                <Label>日期</Label>
                <Input
                  type="date"
                  className="mt-1"
                  value={form.date}
                  onChange={(e) => setForm({ ...form, date: e.target.value })}
                  required
                />
              </div>
              <div>
                <Label>平台</Label>
                <Select
                  className="mt-1"
                  value={form.platform}
                  onChange={(e) => setForm({ ...form, platform: e.target.value })}
                  required
                >
                  <option value="">选择</option>
                  {PLATFORMS.map((p) => (
                    <option key={p.value} value={p.value}>
                      {p.label}
                    </option>
                  ))}
                </Select>
              </div>
              <div>
                <Label>类别</Label>
                <Select
                  className="mt-1"
                  value={form.category}
                  onChange={(e) => setForm({ ...form, category: e.target.value })}
                  required
                >
                  <option value="">选择</option>
                  {INCOME_CATEGORIES.map((c) => (
                    <option key={c.value} value={c.value}>
                      {c.label}
                    </option>
                  ))}
                </Select>
              </div>
              <div>
                <Label>含税收入（元）</Label>
                <Input
                  type="number"
                  min="0"
                  step="0.01"
                  className="mt-1"
                  value={form.grossAmount}
                  onChange={(e) => setForm({ ...form, grossAmount: e.target.value })}
                  required
                />
              </div>
              <div>
                <Label>平台服务费（元）</Label>
                <Input
                  type="number"
                  min="0"
                  step="0.01"
                  className="mt-1"
                  placeholder="无则填 0"
                  value={form.platformFee}
                  onChange={(e) => setForm({ ...form, platformFee: e.target.value })}
                />
              </div>
              <div>
                <Label>实收金额（元）</Label>
                <Input
                  type="text"
                  className="mt-1 bg-muted"
                  readOnly
                  value={netPreview != null && !Number.isNaN(netPreview) ? netPreview.toFixed(2) : '—'}
                />
                <p className="mt-1 text-xs text-muted-foreground">实收 = 含税收入 − 平台服务费</p>
              </div>
              <div className="sm:col-span-2">
                <Label>备注</Label>
                <Input
                  className="mt-1"
                  value={form.note}
                  onChange={(e) => setForm({ ...form, note: e.target.value })}
                />
              </div>
              <Button type="submit" disabled={saving} className="sm:col-span-2">
                {saving ? '保存中…' : '保存'}
              </Button>
            </form>
          </CardContent>
        </Card>
      )}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">筛选与合计</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <div>
              <Label htmlFor="filter-date-from">开始日期</Label>
              <Input
                id="filter-date-from"
                type="date"
                className="mt-1"
                value={filters.dateFrom}
                onChange={(e) => updateFilter('dateFrom', e.target.value)}
                disabled={loading}
              />
              <p className="mt-1 text-xs text-muted-foreground">留空表示不限，需与结束日期成对填写</p>
            </div>
            <div>
              <Label htmlFor="filter-date-to">结束日期</Label>
              <Input
                id="filter-date-to"
                type="date"
                className="mt-1"
                value={filters.dateTo}
                onChange={(e) => updateFilter('dateTo', e.target.value)}
                disabled={loading}
              />
            </div>
            <div>
              <Label htmlFor="filter-platform">平台</Label>
              <Select
                id="filter-platform"
                className="mt-1"
                value={filters.platform}
                onChange={(e) => updateFilter('platform', e.target.value)}
                disabled={loading}
              >
                <option value="">全部平台</option>
                {PLATFORMS.map((p) => (
                  <option key={p.value} value={p.value}>
                    {p.label}
                  </option>
                ))}
              </Select>
            </div>
            <div>
              <Label htmlFor="filter-category">类别</Label>
              <Select
                id="filter-category"
                className="mt-1"
                value={filters.category}
                onChange={(e) => updateFilter('category', e.target.value)}
                disabled={loading}
              >
                <option value="">全部类别</option>
                {INCOME_CATEGORIES.map((c) => (
                  <option key={c.value} value={c.value}>
                    {c.label}
                  </option>
                ))}
              </Select>
            </div>
          </div>
          <div className="grid gap-4 border-t pt-4 sm:grid-cols-3">
            <p className="text-xs text-muted-foreground sm:col-span-3">
              合计为当前筛选条件下全部记录之和（非仅本页）；未选日期时统计全部收入
            </p>
            <div>
              <Label>合计含税收入</Label>
              <Input
                type="text"
                className="mt-1 bg-muted"
                readOnly
                value={loading ? '—' : formatMoney(summary.grossAmount)}
              />
            </div>
            <div>
              <Label>合计服务费</Label>
              <Input
                type="text"
                className="mt-1 bg-muted"
                readOnly
                value={loading ? '—' : formatMoney(summary.platformFee)}
              />
            </div>
            <div>
              <Label>合计实收</Label>
              <Input
                type="text"
                className="mt-1 bg-muted"
                readOnly
                value={loading ? '—' : formatMoney(summary.netAmount)}
              />
            </div>
          </div>
        </CardContent>
      </Card>
      <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-end">
        <Label htmlFor="income-sort" className="shrink-0 text-sm text-muted-foreground">
          排序方式
        </Label>
        <Select
          id="income-sort"
          className="w-full sm:w-56"
          value={sort}
          onChange={handleSortChange}
          disabled={loading}
        >
          {INCOME_SORT_OPTIONS.map((o) => (
            <option key={o.value} value={o.value}>
              {o.label}
            </option>
          ))}
        </Select>
      </div>
      {loading ? (
        <p className="text-muted-foreground">加载中…</p>
      ) : total === 0 ? (
        <Card>
          <CardContent className="py-8 text-center text-muted-foreground">暂无收入记录</CardContent>
        </Card>
      ) : (
        <>
          <div className="space-y-2 md:hidden">
            {items.map((row) => (
              <Card key={row.id}>
                <CardContent className="py-4">
                  <p className="font-medium">实收 ¥{row.netAmount.toLocaleString()}</p>
                  <p className="text-xs text-muted-foreground">
                    {row.date} · {row.platformLabel} · {row.categoryLabel}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    含税 ¥{row.grossAmount.toLocaleString()}，服务费 ¥{row.platformFee.toLocaleString()}
                  </p>
                </CardContent>
              </Card>
            ))}
          </div>
          <div className="hidden overflow-x-auto md:block">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b text-left text-muted-foreground">
                  <th className="p-2">日期</th>
                  <th className="p-2">平台</th>
                  <th className="p-2">类别</th>
                  <th className="p-2 text-right">含税收入</th>
                  <th className="p-2 text-right">服务费</th>
                  <th className="p-2 text-right">实收</th>
                </tr>
              </thead>
              <tbody>
                {items.map((row) => (
                  <tr key={row.id} className="border-b">
                    <td className="p-2">{row.date}</td>
                    <td className="p-2">{row.platformLabel}</td>
                    <td className="p-2">{row.categoryLabel}</td>
                    <td className="p-2 text-right">¥{row.grossAmount.toLocaleString()}</td>
                    <td className="p-2 text-right">¥{row.platformFee.toLocaleString()}</td>
                    <td className="p-2 text-right">¥{row.netAmount.toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {totalPages > 1 && (
            <div className="flex flex-col items-center justify-between gap-3 sm:flex-row">
              <p className="text-sm text-muted-foreground">
                共 {total} 条，第 {page} / {totalPages} 页
              </p>
              <div className="flex gap-2">
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  disabled={page <= 1 || loading}
                  onClick={() => setPage((p) => p - 1)}
                >
                  上一页
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  disabled={page >= totalPages || loading}
                  onClick={() => setPage((p) => p + 1)}
                >
                  下一页
                </Button>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
}
