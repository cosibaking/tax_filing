'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { RISK_DISCLOSURE_ITEMS, LEGAL_DOCUMENT_VERSIONS, isDev } from '@/lib/api/constants';
import { ApiClientError, getMemberToken, memberFetch } from '@/lib/api/client';
import { withBasePath } from '@/lib/base-path';

export default function PlanPage() {
  const [step, setStep] = useState(1);
  const [plans, setPlans] = useState<{ id: string; name: string; price: number }[]>([]);
  const [selectedPlan, setSelectedPlan] = useState('');
  const [orderId, setOrderId] = useState('');
  const [signerName, setSignerName] = useState('');
  const [riskAck, setRiskAck] = useState(false);
  const [consentAck, setConsentAck] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  useEffect(() => {
    if (!getMemberToken()) return;
    memberFetch<{ name?: string | null }>('/api/member/profile')
      .then((profile) => {
        if (profile.name) setSignerName(profile.name);
      })
      .catch(() => {
        /* 未登录或资料未填时不预填 */
      });
  }, []);

  useEffect(() => {
    fetch(withBasePath('/api/site/compliance/plans'))
      .then((r) => r.json())
      .then((json) => {
        if (json.data?.length) setPlans(json.data);
        else
          setPlans([
            { id: 'basic', name: '基础合规', price: 299 },
            { id: 'standard', name: '标准 OPC', price: 899 },
            { id: 'premium', name: '全程托管', price: 1999 },
          ]);
      })
      .catch(() => {
        setPlans([
          { id: 'standard', name: '标准 OPC', price: 899 },
        ]);
      });
  }, []);

  const createOrder = async () => {
    setLoading(true);
    setError(null);
    try {
      const order = await memberFetch<{ id: string }>('/api/member/compliance/order', {
        method: 'POST',
        body: JSON.stringify({ planId: selectedPlan }),
      });
      setOrderId(order.id);
      setStep(2);
    } catch (err) {
      setError(err instanceof ApiClientError ? err.message : '创建订单失败');
    } finally {
      setLoading(false);
    }
  };

  const submitConsent = async () => {
    setLoading(true);
    setError(null);
    try {
      const consents = [
        { type: 'risk_disclosure', documentVersion: LEGAL_DOCUMENT_VERSIONS.risk_disclosure },
        { type: 'plan_confirm', documentVersion: LEGAL_DOCUMENT_VERSIONS.plan_confirm },
      ] as const;

      for (const consent of consents) {
        await memberFetch('/api/member/compliance/consent', {
          method: 'POST',
          body: JSON.stringify({
            type: consent.type,
            documentVersion: consent.documentVersion,
            orderId: orderId || undefined,
          }),
        });
      }
      setStep(3);
    } catch (err) {
      setError(err instanceof ApiClientError ? err.message : '确认失败');
    } finally {
      setLoading(false);
    }
  };

  const sign = async () => {
    setLoading(true);
    setError(null);
    try {
      await memberFetch('/api/member/compliance/sign', {
        method: 'POST',
        body: JSON.stringify({
          orderId,
          signerName: signerName.trim(),
        }),
      });
      setSuccess(true);
    } catch (err) {
      setError(err instanceof ApiClientError ? err.message : '签约失败');
    } finally {
      setLoading(false);
    }
  };

  const step2Ready = isDev || (riskAck && consentAck);
  const step3Ready = isDev || (!!orderId && signerName.trim().length > 0);

  if (success) {
    return (
      <Alert variant="info">
        <AlertTitle>签约成功</AlertTitle>
        <AlertDescription>
          请前往 OPC 进度页查看设立状态。
          <a href="/user/compliance/opc" className="ml-1 underline">
            查看进度
          </a>
        </AlertDescription>
      </Alert>
    );
  }

  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <h1 className="text-2xl font-bold">方案确认与签约</h1>
      <div className="flex gap-2">
        {[1, 2, 3].map((s) => (
          <div key={s} className={`h-2 flex-1 rounded-full ${s <= step ? 'bg-primary' : 'bg-muted'}`} />
        ))}
      </div>
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {step === 1 && (
        <Card>
          <CardHeader>
            <CardTitle>步骤 1：选择套餐</CardTitle>
            <CardDescription>选择适合您的合规服务方案</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            {plans.map((p) => (
              <label
                key={p.id}
                className={`flex cursor-pointer items-center justify-between rounded-lg border p-4 ${selectedPlan === p.id ? 'border-primary bg-primary/5' : ''}`}
              >
                <span>
                  <input
                    type="radio"
                    name="plan"
                    className="mr-3"
                    checked={selectedPlan === p.id}
                    onChange={() => setSelectedPlan(p.id)}
                  />
                  {p.name}
                </span>
                <span className="font-semibold">¥{p.price}/月</span>
              </label>
            ))}
            <Button className="w-full" disabled={!selectedPlan || loading} onClick={createOrder}>
              {loading ? '处理中…' : '确认套餐'}
            </Button>
          </CardContent>
        </Card>
      )}

      {step === 2 && (
        <Card>
          <CardHeader>
            <CardTitle>步骤 2：风险告知</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <ul className="list-inside list-disc space-y-1 text-sm text-muted-foreground">
              {RISK_DISCLOSURE_ITEMS.map((item) => (
                <li key={item}>{item}</li>
              ))}
            </ul>
            <label className="flex items-start gap-2 text-sm">
              <Checkbox checked={riskAck} onCheckedChange={(v) => setRiskAck(!!v)} />
              我已阅读并理解上述风险告知
            </label>
            <label className="flex items-start gap-2 text-sm">
              <Checkbox checked={consentAck} onCheckedChange={(v) => setConsentAck(!!v)} />
              我同意服务范围与合规边界说明
            </label>
            <Button
              className="w-full"
              disabled={!step2Ready || loading}
              onClick={submitConsent}
            >
              确认并继续
            </Button>
          </CardContent>
        </Card>
      )}

      {step === 3 && (
        <Card>
          <CardHeader>
            <CardTitle>步骤 3：电子签约</CardTitle>
            <CardDescription>确认签约即表示您同意服务协议与订单条款</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <Label htmlFor="order-id">订单号</Label>
              <Input
                id="order-id"
                className="mt-1 bg-muted"
                value={orderId}
                readOnly
                placeholder="请先完成步骤 1 创建订单"
              />
            </div>
            <div>
              <Label htmlFor="signer-name">签署人姓名</Label>
              <Input
                id="signer-name"
                className="mt-1"
                value={signerName}
                onChange={(e) => setSignerName(e.target.value)}
                placeholder="请输入与身份证一致的姓名"
                autoComplete="name"
              />
              <p className="mt-1 text-xs text-muted-foreground">
                姓名将写入电子合同，请确保与实名信息一致。
              </p>
            </div>
            <Button className="w-full" disabled={!step3Ready || loading} onClick={sign}>
              {loading ? '签约中…' : '完成电子签约'}
            </Button>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
