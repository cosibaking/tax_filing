'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select } from '@/components/ui/select';
import { Alert, AlertDescription } from '@/components/ui/alert';
import {
  DEFAULT_EXPENSE_SORT,
  ErrorCodes,
  EXPENSE_CATEGORIES,
  EXPENSE_SORT_OPTIONS,
  INVOICE_TYPES,
  type ExpenseSort,
} from '@/lib/api/constants';
import { ApiClientError, memberFetch } from '@/lib/api/client';

interface ApiExpenseEntry {
  id: string;
  category: string;
  amount: string | number;
  invoiceType: string;
  description?: string;
  warningFlag?: boolean;
  occurredAt: string;
}

interface ExpenseRow {
  id: string;
  date: string;
  category: string;
  categoryLabel: string;
  invoiceType: string;
  invoiceTypeLabel: string;
  amount: number;
  description?: string;
  warningFlag: boolean;
}

const PAGE_SIZE = 10;

const EMPTY_FORM = {
  date: '',
  category: '',
  amount: '',
  invoiceType: 'general',
  description: '',
};

interface ExpenseFilters {
  dateFrom: string;
  dateTo: string;
  category: string;
}

interface ExpenseSummary {
  totalAmount: number;
}

interface ExpenseListResponse {
  items: ApiExpenseEntry[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
  sort?: ExpenseSort;
  summary?: ExpenseSummary;
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

function mapEntry(raw: ApiExpenseEntry): ExpenseRow {
  return {
    id: raw.id,
    date: formatDate(raw.occurredAt),
    category: raw.category,
    categoryLabel: labelFor(EXPENSE_CATEGORIES, raw.category),
    invoiceType: raw.invoiceType,
    invoiceTypeLabel: labelFor(INVOICE_TYPES, raw.invoiceType),
    amount: toNumber(raw.amount),
    description: raw.description,
    warningFlag: Boolean(raw.warningFlag),
  };
}

export default function ExpensePage() {
  const [items, setItems] = useState<ExpenseRow[]>([]);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState(EMPTY_FORM);
  const [saving, setSaving] = useState(false);
  const [sort, setSort] = useState<ExpenseSort>(DEFAULT_EXPENSE_SORT);
  const [filters, setFilters] = useState<ExpenseFilters>({
    dateFrom: '',
    dateTo: '',
    category: '',
  });
  const [summary, setSummary] = useState<ExpenseSummary>({ totalAmount: 0 });

  const load = (
    targetPage: number,
    targetSort: ExpenseSort = sort,
    targetFilters: ExpenseFilters = filters,
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
    if (targetFilters.category) params.set('category', targetFilters.category);

    memberFetch<ExpenseListResponse>(`/api/member/compliance/expense?${params}`)
      .then((d) => {
        setItems((d.items ?? []).map(mapEntry));
        setTotal(d.total ?? 0);
        setTotalPages(d.totalPages ?? 0);
        setPage(d.page ?? targetPage);
        if (d.sort) setSort(d.sort);
        setSummary(d.summary ?? { totalAmount: 0 });
        setError(null);
      })
      .catch((err) => setError(err instanceof ApiClientError ? err.message : '加载失败'))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    load(page, sort, filters);
  }, [page, sort, filters.dateFrom, filters.dateTo, filters.category]);

  const resetPageIfNeeded = () => {
    if (page !== 1) setPage(1);
  };

  const handleSortChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const next = e.target.value as ExpenseSort;
    if (next === sort) return;
    setSort(next);
    resetPageIfNeeded();
  };

  const updateFilter = <K extends keyof ExpenseFilters>(key: K, value: ExpenseFilters[K]) => {
    setFilters((prev) => ({ ...prev, [key]: value }));
    resetPageIfNeeded();
  };

  const submitExpense = async (forceConfirm = false) => {
    const amount = parseFloat(form.amount);
    if (Number.isNaN(amount) || amount < 0) {
      setError('请填写有效的金额');
      return;
    }
    await memberFetch('/api/member/compliance/expense', {
      method: 'POST',
      body: JSON.stringify({
        category: form.category,
        amount,
        invoiceType: form.invoiceType,
        occurredAt: form.date,
        description: form.description || undefined,
        forceConfirm,
      }),
    });
    setShowForm(false);
    setForm(EMPTY_FORM);
    if (page === 1) load(1);
    else setPage(1);
  };

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSaving(true);
    try {
      await submitExpense(false);
    } catch (err) {
      if (err instanceof ApiClientError && err.code === ErrorCodes.COST_RATIO_WARNING) {
        if (window.confirm(`${err.message}\n\n是否仍要保存？`)) {
          try {
            await submitExpense(true);
            return;
          } catch (retryErr) {
            setError(retryErr instanceof ApiClientError ? retryErr.message : '保存失败');
          }
        } else {
          setError(err.message);
        }
      } else {
        setError(err instanceof ApiClientError ? err.message : '保存失败');
      }
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <h1 className="text-2xl font-bold">费用台账</h1>
        <Button onClick={() => setShowForm(!showForm)}>{showForm ? '取消' : '记一笔费用'}</Button>
      </div>
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      {showForm && (
        <Card>
          <CardHeader>
            <CardTitle>新增费用</CardTitle>
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
                <Label>类别</Label>
                <Select
                  className="mt-1"
                  value={form.category}
                  onChange={(e) => setForm({ ...form, category: e.target.value })}
                  required
                >
                  <option value="">选择</option>
                  {EXPENSE_CATEGORIES.map((c) => (
                    <option key={c.value} value={c.value}>
                      {c.label}
                    </option>
                  ))}
                </Select>
              </div>
              <div>
                <Label>金额（元）</Label>
                <Input
                  type="number"
                  min="0"
                  step="0.01"
                  className="mt-1"
                  value={form.amount}
                  onChange={(e) => setForm({ ...form, amount: e.target.value })}
                  required
                />
              </div>
              <div>
                <Label>发票类型</Label>
                <Select
                  className="mt-1"
                  value={form.invoiceType}
                  onChange={(e) => setForm({ ...form, invoiceType: e.target.value })}
                  required
                >
                  {INVOICE_TYPES.map((t) => (
                    <option key={t.value} value={t.value}>
                      {t.label}
                    </option>
                  ))}
                </Select>
              </div>
              <div className="sm:col-span-2">
                <Label>说明</Label>
                <Input
                  className="mt-1"
                  value={form.description}
                  onChange={(e) => setForm({ ...form, description: e.target.value })}
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
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
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
              <Label htmlFor="filter-category">类别</Label>
              <Select
                id="filter-category"
                className="mt-1"
                value={filters.category}
                onChange={(e) => updateFilter('category', e.target.value)}
                disabled={loading}
              >
                <option value="">全部类别</option>
                {EXPENSE_CATEGORIES.map((c) => (
                  <option key={c.value} value={c.value}>
                    {c.label}
                  </option>
                ))}
              </Select>
            </div>
          </div>
          <div className="border-t pt-4">
            <p className="mb-3 text-xs text-muted-foreground">
              合计为当前筛选条件下全部记录之和（非仅本页）；未选日期时统计全部费用
            </p>
            <div className="max-w-sm">
              <Label>合计金额</Label>
              <Input
                type="text"
                className="mt-1 bg-muted"
                readOnly
                value={loading ? '—' : formatMoney(summary.totalAmount)}
              />
            </div>
          </div>
        </CardContent>
      </Card>
      <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-end">
        <Label htmlFor="expense-sort" className="shrink-0 text-sm text-muted-foreground">
          排序方式
        </Label>
        <Select
          id="expense-sort"
          className="w-full sm:w-56"
          value={sort}
          onChange={handleSortChange}
          disabled={loading}
        >
          {EXPENSE_SORT_OPTIONS.map((o) => (
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
          <CardContent className="py-8 text-center text-muted-foreground">暂无费用记录</CardContent>
        </Card>
      ) : (
        <>
          <div className="space-y-2 md:hidden">
            {items.map((row) => (
              <Card key={row.id}>
                <CardContent className="py-4">
                  <p className="font-medium">{formatMoney(row.amount)}</p>
                  <p className="text-xs text-muted-foreground">
                    {row.date} · {row.categoryLabel} · {row.invoiceTypeLabel}
                  </p>
                  {row.description && (
                    <p className="text-xs text-muted-foreground">{row.description}</p>
                  )}
                  {row.warningFlag && (
                    <p className="text-xs text-amber-600">成本占比预警</p>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
          <div className="hidden overflow-x-auto md:block">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b text-left text-muted-foreground">
                  <th className="p-2">日期</th>
                  <th className="p-2">类别</th>
                  <th className="p-2">发票类型</th>
                  <th className="p-2 text-right">金额</th>
                  <th className="p-2">说明</th>
                </tr>
              </thead>
              <tbody>
                {items.map((row) => (
                  <tr key={row.id} className="border-b">
                    <td className="p-2">{row.date}</td>
                    <td className="p-2">
                      {row.categoryLabel}
                      {row.warningFlag && (
                        <span className="ml-1 text-xs text-amber-600">预警</span>
                      )}
                    </td>
                    <td className="p-2">{row.invoiceTypeLabel}</td>
                    <td className="p-2 text-right">{formatMoney(row.amount)}</td>
                    <td className="p-2 text-muted-foreground">{row.description ?? '—'}</td>
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
