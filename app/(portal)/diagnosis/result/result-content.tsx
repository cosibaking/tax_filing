'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { useSearchParams } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { ApiClientError, publicFetch } from '@/lib/api/client';

interface TaxComparisonItem {
  mode: string;
  label: string;
  annualTax: number;
  effectiveRate: string;
  highlight?: boolean;
  warning?: string;
}

interface RawTaxComparisonItem {
  plan?: string;
  mode?: string;
  label: string;
  taxAmount?: number;
  annualTax?: number;
  effectiveRate?: number | string;
  warning?: string;
  highlight?: boolean;
}

interface DiagnosisResult {
  id: string;
  recommendedPlan?: string;
  recommendedLabel?: string;
  taxComparison?: RawTaxComparisonItem[];
  warning?: string;
}

const PLAN_LABELS: Record<string, string> = {
  opc: 'OPC 小微公司',
  individual: '个体工商户',
  transitional: '过渡方案',
  labor: '劳务报酬（现状）',
  none: '不报税',
};

const DEFAULT_COMPARISON: TaxComparisonItem[] = [
  { mode: 'labor', label: '劳务报酬', annualTax: 270000, effectiveRate: '27%', highlight: false },
  { mode: 'individual', label: '个体户', annualTax: 180000, effectiveRate: '18%', highlight: false },
  { mode: 'opc', label: 'OPC', annualTax: 120000, effectiveRate: '12%', highlight: true },
];

function formatEffectiveRate(rate: number | string | undefined): string {
  if (typeof rate === 'number') return `${(rate * 100).toFixed(1)}%`;
  if (typeof rate === 'string' && rate.length > 0) return rate;
  return '—';
}

function normalizeComparison(
  items: RawTaxComparisonItem[] | undefined,
  recommendedPlan: string,
): TaxComparisonItem[] {
  if (!items?.length) return DEFAULT_COMPARISON;

  return items.map((item, index) => {
    const mode = item.plan ?? item.mode ?? `item-${index}`;
    const annualTax = item.annualTax ?? item.taxAmount ?? 0;
    return {
      mode,
      label: item.label,
      annualTax,
      effectiveRate: formatEffectiveRate(item.effectiveRate),
      highlight: item.highlight ?? mode === recommendedPlan,
      warning: item.warning,
    };
  });
}

export function DiagnosisResultContent() {
  const searchParams = useSearchParams();
  const id = searchParams.get('id');
  const [data, setData] = useState<DiagnosisResult | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const cached = sessionStorage.getItem('diagnosis_result');
    if (cached) {
      try {
        setData(JSON.parse(cached) as DiagnosisResult);
        setLoading(false);
        return;
      } catch {
        /* fall through */
      }
    }
    if (!id) {
      setError('未找到诊断记录');
      setLoading(false);
      return;
    }
    publicFetch<DiagnosisResult>(`/api/site/compliance/diagnosis/${id}`)
      .then(setData)
      .catch((err) => {
        setError(err instanceof ApiClientError ? err.message : '加载失败');
      })
      .finally(() => setLoading(false));
  }, [id]);

  if (loading) {
    return (
      <div className="mx-auto max-w-4xl px-4 py-16 text-center text-muted-foreground">
        正在加载诊断结果…
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="mx-auto max-w-2xl px-4 py-16">
        <Alert variant="destructive">
          <AlertDescription>{error ?? '暂无结果'}</AlertDescription>
        </Alert>
        <Button asChild className="mt-4">
          <Link href="/diagnosis">重新诊断</Link>
        </Button>
      </div>
    );
  }

  const recommended = data.recommendedPlan ?? 'opc';
  const comparison = normalizeComparison(data.taxComparison, recommended);

  return (
    <div className="mx-auto max-w-4xl px-4 py-8">
      <div className="mb-8 text-center">
        <Badge className="mb-2">诊断完成</Badge>
        <h1 className="text-2xl font-bold sm:text-3xl">您的合规诊断结果</h1>
        <p className="mt-2 text-muted-foreground">
          推荐方案：
          <span className="font-semibold text-primary">
            {data.recommendedLabel ?? PLAN_LABELS[recommended] ?? recommended}
          </span>
        </p>
      </div>

      {data.warning && (
        <Alert variant="warning" className="mb-6">
          <AlertDescription>{data.warning}</AlertDescription>
        </Alert>
      )}

      <h2 className="mb-4 text-lg font-semibold">税负对比（年估）</h2>
      <div className="grid gap-4 sm:grid-cols-3">
        {comparison.map((item) => (
          <Card
            key={item.mode}
            className={item.highlight ? 'border-primary ring-2 ring-primary/20' : ''}
          >
            <CardHeader>
              <CardTitle className="text-lg">{item.label}</CardTitle>
              <CardDescription>有效税率约 {item.effectiveRate}</CardDescription>
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-bold">
                ¥{(item.annualTax / 10000).toFixed(1)}
                <span className="text-base font-normal text-muted-foreground"> 万/年</span>
              </p>
              {item.warning && (
                <p className="mt-2 text-xs text-destructive">{item.warning}</p>
              )}
              {item.highlight && <Badge variant="success" className="mt-2">推荐</Badge>}
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="mt-8 flex flex-col gap-3 sm:flex-row sm:justify-center">
        <Button asChild size="lg">
          <Link href="/pricing">查看服务套餐</Link>
        </Button>
        <Button asChild variant="outline" size="lg">
          <Link href="/user/register">注册并签约</Link>
        </Button>
      </div>

      <p className="mt-6 text-center text-xs text-muted-foreground">
        以上为估算结果，实际税负取决于真实收入与成本凭证。本服务为合规方案，非逃税方案。
      </p>
    </div>
  );
}
