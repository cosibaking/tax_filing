'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { ApiClientError, publicFetch } from '@/lib/api/client';

interface Plan {
  id: string;
  name: string;
  price: number;
  period: string;
  features: string[];
  recommended?: boolean;
}

const FALLBACK_PLANS: Plan[] = [
  {
    id: 'basic',
    name: '基础合规',
    price: 299,
    period: '月',
    features: ['合规诊断', '政策解读', '申报提醒'],
  },
  {
    id: 'standard',
    name: '标准 OPC',
    price: 899,
    period: '月',
    features: ['OPC 设立代办', '收入费用台账', '季度申报协助', '月度对账单'],
    recommended: true,
  },
  {
    id: 'premium',
    name: '全程托管',
    price: 1999,
    period: '月',
    features: ['标准版全部', '专属顾问', '材料催收', '税局沟通协助'],
  },
];

export default function PricingPage() {
  const [plans, setPlans] = useState<Plan[]>(FALLBACK_PLANS);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    publicFetch<Plan[]>('/api/site/compliance/plans')
      .then((data) => {
        if (Array.isArray(data) && data.length > 0) setPlans(data);
      })
      .catch((err) => {
        if (err instanceof ApiClientError) setError(err.message);
      })
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="mx-auto max-w-6xl px-4 py-12">
      <div className="mb-10 text-center">
        <h1 className="text-3xl font-bold">服务套餐</h1>
        <p className="mt-2 text-muted-foreground">选择适合您的合规服务方案</p>
      </div>

      {error && (
        <Alert variant="warning" className="mb-6 max-w-2xl mx-auto">
          <AlertDescription>无法加载最新套餐，展示默认方案。{error}</AlertDescription>
        </Alert>
      )}

      {loading ? (
        <p className="text-center text-muted-foreground">加载中…</p>
      ) : (
        <div className="grid gap-6 md:grid-cols-3">
          {plans.map((plan) => (
            <Card
              key={plan.id}
              className={plan.recommended ? 'border-primary shadow-md md:-mt-2 md:mb-2' : ''}
            >
              <CardHeader>
                {plan.recommended && <Badge className="mb-2 w-fit">推荐</Badge>}
                <CardTitle>{plan.name}</CardTitle>
                <CardDescription>
                  <span className="text-3xl font-bold text-foreground">¥{plan.price}</span>
                  <span className="text-muted-foreground"> /{plan.period}</span>
                </CardDescription>
              </CardHeader>
              <CardContent>
                <ul className="space-y-2 text-sm">
                  {plan.features.map((f) => (
                    <li key={f} className="flex items-start gap-2">
                      <span className="text-primary">✓</span>
                      {f}
                    </li>
                  ))}
                </ul>
              </CardContent>
              <CardFooter>
                <Button asChild className="w-full" variant={plan.recommended ? 'default' : 'outline'}>
                  <Link href="/user/compliance/plan">选择方案</Link>
                </Button>
              </CardFooter>
            </Card>
          ))}
        </div>
      )}

      <p className="mt-8 text-center text-sm text-muted-foreground">
        不确定选哪个？
        <Link href="/diagnosis" className="text-primary underline">
          先做免费诊断
        </Link>
      </p>
    </div>
  );
}
