'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select } from '@/components/ui/select';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { INCOME_CATEGORIES, PLATFORMS } from '@/lib/api/constants';
import { ApiClientError, memberFetch } from '@/lib/api/client';

interface IncomeRow {
  id: string;
  date: string;
  platform: string;
  category: string;
  amount: number;
  note?: string;
}

export default function IncomePage() {
  const [items, setItems] = useState<IncomeRow[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ date: '', platform: '', category: '', amount: '', note: '' });
  const [saving, setSaving] = useState(false);

  const load = () => {
    memberFetch<{ items: IncomeRow[] }>('/api/member/compliance/income')
      .then((d) => setItems(d.items ?? (Array.isArray(d) ? (d as unknown as IncomeRow[]) : [])))
      .catch((err) => setError(err instanceof ApiClientError ? err.message : '加载失败'))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    load();
  }, []);

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      await memberFetch('/api/member/compliance/income', {
        method: 'POST',
        body: JSON.stringify({
          ...form,
          amount: Number(form.amount),
        }),
      });
      setShowForm(false);
      setForm({ date: '', platform: '', category: '', amount: '', note: '' });
      load();
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
                <Input type="date" className="mt-1" value={form.date} onChange={(e) => setForm({ ...form, date: e.target.value })} required />
              </div>
              <div>
                <Label>平台</Label>
                <Select className="mt-1" value={form.platform} onChange={(e) => setForm({ ...form, platform: e.target.value })} required>
                  <option value="">选择</option>
                  {PLATFORMS.map((p) => (
                    <option key={p.value} value={p.value}>{p.label}</option>
                  ))}
                </Select>
              </div>
              <div>
                <Label>类别</Label>
                <Select className="mt-1" value={form.category} onChange={(e) => setForm({ ...form, category: e.target.value })} required>
                  <option value="">选择</option>
                  {INCOME_CATEGORIES.map((c) => (
                    <option key={c.value} value={c.value}>{c.label}</option>
                  ))}
                </Select>
              </div>
              <div>
                <Label>金额（元）</Label>
                <Input type="number" className="mt-1" value={form.amount} onChange={(e) => setForm({ ...form, amount: e.target.value })} required />
              </div>
              <Button type="submit" disabled={saving} className="sm:col-span-2">
                {saving ? '保存中…' : '保存'}
              </Button>
            </form>
          </CardContent>
        </Card>
      )}
      {loading ? (
        <p className="text-muted-foreground">加载中…</p>
      ) : items.length === 0 ? (
        <Card><CardContent className="py-8 text-center text-muted-foreground">暂无收入记录</CardContent></Card>
      ) : (
        <div className="space-y-2 md:hidden">
          {items.map((row) => (
            <Card key={row.id}>
              <CardContent className="flex justify-between py-4">
                <div>
                  <p className="font-medium">¥{row.amount.toLocaleString()}</p>
                  <p className="text-xs text-muted-foreground">{row.date} · {row.category}</p>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
      {!loading && items.length > 0 && (
        <div className="hidden overflow-x-auto md:block">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b text-left text-muted-foreground">
                <th className="p-2">日期</th>
                <th className="p-2">平台</th>
                <th className="p-2">类别</th>
                <th className="p-2 text-right">金额</th>
              </tr>
            </thead>
            <tbody>
              {items.map((row) => (
                <tr key={row.id} className="border-b">
                  <td className="p-2">{row.date}</td>
                  <td className="p-2">{row.platform}</td>
                  <td className="p-2">{row.category}</td>
                  <td className="p-2 text-right">¥{row.amount.toLocaleString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
