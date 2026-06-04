'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Checkbox } from '@/components/ui/checkbox';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select } from '@/components/ui/select';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { PLATFORMS } from '@/lib/api/constants';
import { ApiClientError, publicFetch } from '@/lib/api/client';

const INCOME_RANGES = [
  { value: '0-2', label: '0–2 万/月' },
  { value: '2-5', label: '2–5 万/月' },
  { value: '5-15', label: '5–15 万/月' },
  { value: '15+', label: '15 万+/月' },
];

const ENTITIES = [
  { value: 'none', label: '无主体' },
  { value: 'individual', label: '个体户' },
  { value: 'company', label: '公司' },
  { value: 'other', label: '其他' },
];

export default function DiagnosisPage() {
  const router = useRouter();
  const [step, setStep] = useState(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [platforms, setPlatforms] = useState<string[]>([]);
  const [monthlyIncomeRange, setMonthlyIncomeRange] = useState('');
  const [annualCostEstimate, setAnnualCostEstimate] = useState('');
  const [existingEntity, setExistingEntity] = useState('none');
  const [hasFiledTax, setHasFiledTax] = useState('yes');
  const [taxBureauContact, setTaxBureauContact] = useState('no');
  const [notes, setNotes] = useState('');

  const togglePlatform = (value: string) => {
    setPlatforms((prev) =>
      prev.includes(value) ? prev.filter((p) => p !== value) : [...prev, value],
    );
  };

  const handleSubmit = async () => {
    setError(null);
    setLoading(true);
    try {
      const data = await publicFetch<{ id: string; taxComparison?: unknown }>(
        '/api/site/compliance/diagnosis',
        {
          method: 'POST',
          body: JSON.stringify({
            platforms,
            monthlyIncomeRange,
            annualCostEstimate: annualCostEstimate ? Number(annualCostEstimate) : 0,
            existingEntity,
            hasFiledTax: hasFiledTax === 'yes',
            taxBureauContact: taxBureauContact === 'yes',
            notes,
          }),
        },
      );
      sessionStorage.setItem('diagnosis_result', JSON.stringify(data));
      router.push(`/diagnosis/result?id=${data.id}`);
    } catch (err) {
      setError(err instanceof ApiClientError ? err.message : '提交失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  const canNext =
    (step === 1 && platforms.length > 0 && monthlyIncomeRange) ||
    (step === 2 && existingEntity && hasFiledTax) ||
    step === 3;

  return (
    <div className="mx-auto max-w-2xl px-4 py-8">
      <div className="mb-6">
        <h1 className="text-2xl font-bold">合规诊断问卷</h1>
        <p className="text-muted-foreground">约 10 分钟，了解您的税负与合规风险</p>
        <div className="mt-4 flex gap-2">
          {[1, 2, 3, 4].map((s) => (
            <div
              key={s}
              className={`h-2 flex-1 rounded-full ${s <= step ? 'bg-primary' : 'bg-muted'}`}
            />
          ))}
        </div>
        <p className="mt-2 text-sm text-muted-foreground">步骤 {step} / 4</p>
      </div>

      {error && (
        <Alert variant="destructive" className="mb-4">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <Card>
        <CardHeader>
          <CardTitle>
            {step === 1 && '平台与收入'}
            {step === 2 && '现有主体'}
            {step === 3 && '风险信号'}
            {step === 4 && '确认提交'}
          </CardTitle>
          {step === 4 && (
            <CardDescription>请确认信息无误后提交诊断</CardDescription>
          )}
        </CardHeader>
        <CardContent className="space-y-4">
          {step === 1 && (
            <>
              <div>
                <Label className="mb-2 block">主要平台（可多选）</Label>
                <div className="grid gap-3 sm:grid-cols-2">
                  {PLATFORMS.map((p) => (
                    <label key={p.value} className="flex items-center gap-2 text-sm">
                      <Checkbox
                        checked={platforms.includes(p.value)}
                        onCheckedChange={() => togglePlatform(p.value)}
                      />
                      {p.label}
                    </label>
                  ))}
                </div>
              </div>
              <div>
                <Label htmlFor="income">月收入区间</Label>
                <Select
                  id="income"
                  className="mt-1"
                  value={monthlyIncomeRange}
                  onChange={(e) => setMonthlyIncomeRange(e.target.value)}
                >
                  <option value="">请选择</option>
                  {INCOME_RANGES.map((r) => (
                    <option key={r.value} value={r.value}>
                      {r.label}
                    </option>
                  ))}
                </Select>
              </div>
              <div>
                <Label htmlFor="cost">年可扣除成本估算（元，选填）</Label>
                <Input
                  id="cost"
                  type="number"
                  className="mt-1"
                  placeholder="如设备、团队、场地等"
                  value={annualCostEstimate}
                  onChange={(e) => setAnnualCostEstimate(e.target.value)}
                />
              </div>
            </>
          )}

          {step === 2 && (
            <>
              <div>
                <Label htmlFor="entity">是否已有经营主体</Label>
                <Select
                  id="entity"
                  className="mt-1"
                  value={existingEntity}
                  onChange={(e) => setExistingEntity(e.target.value)}
                >
                  {ENTITIES.map((e) => (
                    <option key={e.value} value={e.value}>
                      {e.label}
                    </option>
                  ))}
                </Select>
              </div>
              <div>
                <Label htmlFor="filed">是否按时报税</Label>
                <Select
                  id="filed"
                  className="mt-1"
                  value={hasFiledTax}
                  onChange={(e) => setHasFiledTax(e.target.value)}
                >
                  <option value="yes">是</option>
                  <option value="no">否</option>
                  <option value="unknown">不确定</option>
                </Select>
              </div>
            </>
          )}

          {step === 3 && (
            <>
              <div>
                <Label htmlFor="contact">是否被税务局联系过</Label>
                <Select
                  id="contact"
                  className="mt-1"
                  value={taxBureauContact}
                  onChange={(e) => setTaxBureauContact(e.target.value)}
                >
                  <option value="no">否</option>
                  <option value="yes">是</option>
                </Select>
              </div>
              <div>
                <Label htmlFor="notes">补充说明（选填）</Label>
                <textarea
                  id="notes"
                  className="mt-1 flex min-h-[100px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  value={notes}
                  onChange={(e) => setNotes(e.target.value)}
                />
              </div>
            </>
          )}

          {step === 4 && (
            <dl className="space-y-2 text-sm">
              <div className="flex justify-between">
                <dt className="text-muted-foreground">平台</dt>
                <dd>{platforms.map((p) => PLATFORMS.find((x) => x.value === p)?.label ?? p).join('、')}</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-muted-foreground">月收入</dt>
                <dd>{INCOME_RANGES.find((r) => r.value === monthlyIncomeRange)?.label}</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-muted-foreground">现有主体</dt>
                <dd>{ENTITIES.find((e) => e.value === existingEntity)?.label}</dd>
              </div>
            </dl>
          )}

          <div className="flex gap-3 pt-4">
            {step > 1 && (
              <Button variant="outline" onClick={() => setStep((s) => s - 1)} disabled={loading}>
                上一步
              </Button>
            )}
            {step < 4 ? (
              <Button className="flex-1" disabled={!canNext} onClick={() => setStep((s) => s + 1)}>
                下一步
              </Button>
            ) : (
              <Button className="flex-1" onClick={handleSubmit} disabled={loading}>
                {loading ? '分析中…' : '查看诊断结果'}
              </Button>
            )}
          </div>
        </CardContent>
      </Card>

      <p className="mt-4 text-center text-xs text-muted-foreground">
        提交即表示您已阅读
        <Link href="/legal/privacy" className="underline">
          隐私政策
        </Link>
      </p>
    </div>
  );
}
