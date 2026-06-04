'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Checkbox } from '@/components/ui/checkbox';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { ApiClientError, publicFetch, setMemberTokens } from '@/lib/api/client';

export default function RegisterPage() {
  const router = useRouter();
  const [phone, setPhone] = useState('');
  const [password, setPassword] = useState('');
  const [otp, setOtp] = useState('');
  const [agreeTerms, setAgreeTerms] = useState(false);
  const [agreePrivacy, setAgreePrivacy] = useState(false);
  const [agreeRisk, setAgreeRisk] = useState(false);
  const [loading, setLoading] = useState(false);
  const [otpSending, setOtpSending] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const sendOtp = async () => {
    if (!phone) return setError('请输入手机号');
    setOtpSending(true);
    try {
      await publicFetch<null>('/api/member/auth/otp/send', {
        method: 'POST',
        body: JSON.stringify({ phone }),
      });
    } catch (err) {
      setError(err instanceof ApiClientError ? err.message : '发送失败');
    } finally {
      setOtpSending(false);
    }
  };

  const handleRegister = async () => {
    if (!agreeTerms || !agreePrivacy || !agreeRisk) {
      setError('请勾选全部协议');
      return;
    }
    setLoading(true);
    setError(null);
    try {
      const data = await publicFetch<{
        accessToken: string;
        refreshToken: string;
      }>('/api/member/auth/register', {
        method: 'POST',
        body: JSON.stringify({ phone, password, otp, agreeTerms, agreePrivacy, agreeRisk }),
      });
      setMemberTokens(data.accessToken, data.refreshToken);
      router.push('/user/overview');
    } catch (err) {
      setError(err instanceof ApiClientError ? err.message : '注册失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="mx-auto max-w-md py-8">
      <Card>
        <CardHeader>
          <CardTitle>注册账号</CardTitle>
          <CardDescription>创建会员账号，开启合规服务</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {error && (
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}
          <div>
            <Label htmlFor="phone">手机号</Label>
            <Input id="phone" type="tel" className="mt-1" value={phone} onChange={(e) => setPhone(e.target.value)} />
          </div>
          <div className="flex gap-2">
            <div className="flex-1">
              <Label htmlFor="otp">验证码</Label>
              <Input id="otp" className="mt-1" value={otp} onChange={(e) => setOtp(e.target.value)} />
            </div>
            <Button variant="outline" className="mt-auto" disabled={otpSending} onClick={sendOtp}>
              获取验证码
            </Button>
          </div>
          <div>
            <Label htmlFor="password">设置密码</Label>
            <Input
              id="password"
              type="password"
              className="mt-1"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>
          <div className="space-y-3 rounded-md border p-3 text-sm">
            <label className="flex items-start gap-2">
              <Checkbox checked={agreeTerms} onCheckedChange={(v) => setAgreeTerms(!!v)} />
              <span>
                我已阅读并同意
                <Link href="/legal/terms" className="text-primary underline" target="_blank">
                  用户协议
                </Link>
              </span>
            </label>
            <label className="flex items-start gap-2">
              <Checkbox checked={agreePrivacy} onCheckedChange={(v) => setAgreePrivacy(!!v)} />
              <span>
                我已阅读并同意
                <Link href="/legal/privacy" className="text-primary underline" target="_blank">
                  隐私政策
                </Link>
              </span>
            </label>
            <label className="flex items-start gap-2">
              <Checkbox checked={agreeRisk} onCheckedChange={(v) => setAgreeRisk(!!v)} />
              <span>我理解本服务为合规方案，非逃税方案，不承诺包不被查</span>
            </label>
          </div>
          <Button className="w-full" disabled={loading} onClick={handleRegister}>
            {loading ? '注册中…' : '注册'}
          </Button>
          <p className="text-center text-sm text-muted-foreground">
            已有账号？
            <Link href="/user/login" className="text-primary underline">
              去登录
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
