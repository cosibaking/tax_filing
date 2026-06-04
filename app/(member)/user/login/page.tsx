'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { ApiClientError, publicFetch, setMemberTokens } from '@/lib/api/client';

export default function LoginPage() {
  const router = useRouter();
  const [phone, setPhone] = useState('');
  const [password, setPassword] = useState('');
  const [otp, setOtp] = useState('');
  const [loading, setLoading] = useState(false);
  const [otpSending, setOtpSending] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const login = async (mode: 'password' | 'otp') => {
    setError(null);
    setLoading(true);
    try {
      const data = await publicFetch<{
        accessToken: string;
        refreshToken: string;
      }>('/api/member/auth/login', {
        method: 'POST',
        body: JSON.stringify({
          phone,
          password: mode === 'password' ? password : undefined,
          otp: mode === 'otp' ? otp : undefined,
        }),
      });
      setMemberTokens(data.accessToken, data.refreshToken);
      router.push('/user/overview');
    } catch (err) {
      setError(err instanceof ApiClientError ? err.message : '登录失败');
    } finally {
      setLoading(false);
    }
  };

  const sendOtp = async () => {
    if (!phone) {
      setError('请先输入手机号');
      return;
    }
    setOtpSending(true);
    setError(null);
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

  return (
    <div className="mx-auto max-w-md py-8">
      <Card>
        <CardHeader>
          <CardTitle>会员登录</CardTitle>
          <CardDescription>使用手机号登录合规中心</CardDescription>
        </CardHeader>
        <CardContent>
          {error && (
            <Alert variant="destructive" className="mb-4">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}
          <div className="mb-4">
            <Label htmlFor="phone">手机号</Label>
            <Input
              id="phone"
              type="tel"
              className="mt-1"
              placeholder="11 位手机号"
              value={phone}
              onChange={(e) => setPhone(e.target.value)}
            />
          </div>
          <Tabs defaultValue="password">
            <TabsList>
              <TabsTrigger value="password">密码登录</TabsTrigger>
              <TabsTrigger value="otp">验证码登录</TabsTrigger>
            </TabsList>
            <TabsContent value="password" className="space-y-4">
              <div>
                <Label htmlFor="password">密码</Label>
                <Input
                  id="password"
                  type="password"
                  className="mt-1"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                />
              </div>
              <Button className="w-full" disabled={loading} onClick={() => login('password')}>
                {loading ? '登录中…' : '登录'}
              </Button>
            </TabsContent>
            <TabsContent value="otp" className="space-y-4">
              <div className="flex gap-2">
                <div className="flex-1">
                  <Label htmlFor="otp">验证码</Label>
                  <Input
                    id="otp"
                    className="mt-1"
                    value={otp}
                    onChange={(e) => setOtp(e.target.value)}
                  />
                </div>
                <Button
                  variant="outline"
                  className="mt-auto shrink-0"
                  disabled={otpSending}
                  onClick={sendOtp}
                >
                  {otpSending ? '发送中' : '获取验证码'}
                </Button>
              </div>
              <Button className="w-full" disabled={loading} onClick={() => login('otp')}>
                {loading ? '登录中…' : '登录'}
              </Button>
            </TabsContent>
          </Tabs>
          <p className="mt-4 text-center text-sm text-muted-foreground">
            还没有账号？
            <Link href="/user/register" className="text-primary underline">
              立即注册
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
