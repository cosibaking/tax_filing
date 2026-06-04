'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { ApiClientError, memberFetch } from '@/lib/api/client';

interface ExpenseRow {
  id: string;
  date: string;
  category: string;
  amount: number;
  description?: string;
}

export default function ExpensePage() {
  const [items, setItems] = useState<ExpenseRow[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ date: '', category: '', amount: '', description: '' });
  const [saving, setSaving] = useState(false);

  const load = () => {
    memberFetch<{ items: ExpenseRow[] }>('/api/member/compliance/expense')
      .then((d) => setItems(d.items ?? []))
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
      await memberFetch('/api/member/compliance/expense', {
        method: 'POST',
        body: JSON.stringify({ ...form, amount: Number(form.amount) }),
      });
      setShowForm(false);
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
          <CardHeader><CardTitle>新增费用</CardTitle></CardHeader>
          <CardContent>
            <form onSubmit={handleAdd} className="grid gap-4 sm:grid-cols-2">
              <div>
                <Label>日期</Label>
                <Input type="date" className="mt-1" value={form.date} onChange={(e) => setForm({ ...form, date: e.target.value })} required />
              </div>
              <div>
                <Label>类别</Label>
                <Input className="mt-1" value={form.category} onChange={(e) => setForm({ ...form, category: e.target.value })} required />
              </div>
              <div>
                <Label>金额（元）</Label>
                <Input type="number" className="mt-1" value={form.amount} onChange={(e) => setForm({ ...form, amount: e.target.value })} required />
              </div>
              <div className="sm:col-span-2">
                <Label>说明</Label>
                <Input className="mt-1" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} />
              </div>
              <Button type="submit" disabled={saving} className="sm:col-span-2">{saving ? '保存中…' : '保存'}</Button>
            </form>
          </CardContent>
        </Card>
      )}
      {loading ? (
        <p className="text-muted-foreground">加载中…</p>
      ) : items.length === 0 ? (
        <Card><CardContent className="py-8 text-center text-muted-foreground">暂无费用记录</CardContent></Card>
      ) : (
        <>
          <div className="space-y-2 md:hidden">
            {items.map((row) => (
              <Card key={row.id}>
                <CardContent className="py-4">
                  <p className="font-medium">¥{row.amount.toLocaleString()}</p>
                  <p className="text-xs text-muted-foreground">{row.date} · {row.category}</p>
                </CardContent>
              </Card>
            ))}
          </div>
          <div className="hidden overflow-x-auto md:block">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b text-muted-foreground">
                  <th className="p-2 text-left">日期</th>
                  <th className="p-2 text-left">类别</th>
                  <th className="p-2 text-right">金额</th>
                </tr>
              </thead>
              <tbody>
                {items.map((row) => (
                  <tr key={row.id} className="border-b">
                    <td className="p-2">{row.date}</td>
                    <td className="p-2">{row.category}</td>
                    <td className="p-2 text-right">¥{row.amount.toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </>
      )}
    </div>
  );
}
