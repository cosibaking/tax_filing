import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { RISK_DISCLOSURE_ITEMS } from '@/lib/api/constants';

export default function HomePage() {
  return (
    <div>
      <section className="bg-gradient-to-b from-primary/5 to-background px-4 py-16 sm:py-24">
        <div className="mx-auto max-w-4xl text-center">
          <h1 className="text-3xl font-bold tracking-tight sm:text-4xl md:text-5xl">
            个人主播合规，从诊断到申报一站搞定
          </h1>
          <p className="mt-4 text-lg text-muted-foreground sm:text-xl">
            OPC 设立、收入费用台账、申报提醒与月度对账单，让合规可执行、可追踪。
          </p>
          <div className="mt-8 flex flex-col gap-3 sm:flex-row sm:justify-center">
            <Button asChild size="lg">
              <Link href="/diagnosis">免费合规诊断</Link>
            </Button>
            <Button asChild variant="outline" size="lg">
              <Link href="/pricing">查看服务套餐</Link>
            </Button>
          </div>
        </div>
      </section>

      <section className="mx-auto max-w-6xl px-4 py-12">
        <h2 className="mb-6 text-center text-2xl font-semibold">三步开启合规服务</h2>
        <div className="grid gap-6 sm:grid-cols-3">
          {[
            { step: '1', title: '免费诊断', desc: '填写问卷，了解税负对比与合规风险' },
            { step: '2', title: '签约 OPC', desc: '选择套餐，完成风险告知与电子签约' },
            { step: '3', title: '记账申报', desc: '录入流水，顾问协助申报与对账' },
          ].map((item) => (
            <Card key={item.step}>
              <CardHeader>
                <CardTitle className="text-primary">步骤 {item.step}</CardTitle>
                <CardDescription>{item.title}</CardDescription>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-muted-foreground">{item.desc}</p>
              </CardContent>
            </Card>
          ))}
        </div>
      </section>

      <section className="mx-auto max-w-4xl px-4 pb-16">
        <Card className="border-amber-200 bg-amber-50/50">
          <CardHeader>
            <CardTitle className="text-base">重要告知</CardTitle>
          </CardHeader>
          <CardContent>
            <ul className="list-inside list-disc space-y-1 text-sm text-muted-foreground">
              {RISK_DISCLOSURE_ITEMS.map((item) => (
                <li key={item}>{item}</li>
              ))}
            </ul>
          </CardContent>
        </Card>
      </section>
    </div>
  );
}
