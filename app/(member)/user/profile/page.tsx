'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { ApiClientError, getMemberToken, memberFetch } from '@/lib/api/client';

interface Profile {
  phone?: string;
  name?: string | null;
  createdAt?: string;
}

interface OpcInfo {
  companyName?: string;
  creditCode?: string;
  opcStatus?: string;
  bankAccount?: string;
}

export default function ProfilePage() {
  const [profile, setProfile] = useState<Profile | null>(null);
  const [opc, setOpc] = useState<OpcInfo | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!getMemberToken()) {
      setLoading(false);
      setError('请先登录');
      return;
    }
    Promise.all([
      memberFetch<Profile>('/api/member/profile'),
      memberFetch<OpcInfo>('/api/member/compliance/opc').catch(() => null),
    ])
      .then(([p, o]) => {
        setProfile(p);
        setOpc(o);
      })
      .catch((err) => {
        setError(err instanceof ApiClientError ? err.message : '加载失败');
      })
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p className="text-muted-foreground">加载中…</p>;
  if (error) {
    return (
      <Alert variant="destructive">
        <AlertDescription>{error}</AlertDescription>
      </Alert>
    );
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">个人资料</h1>
      <Card>
        <CardHeader>
          <CardTitle>账号信息</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2 text-sm">
          <p>
            <span className="text-muted-foreground">手机号：</span>
            {profile?.phone ?? '—'}
          </p>
          <p>
            <span className="text-muted-foreground">姓名：</span>
            {profile?.name ?? '未设置'}
          </p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>OPC 主体绑定</CardTitle>
          <CardDescription>您的合规服务经营主体</CardDescription>
        </CardHeader>
        <CardContent>
          {opc?.companyName ? (
            <dl className="space-y-2 text-sm">
              <div className="flex justify-between gap-4">
                <dt className="text-muted-foreground">公司名称</dt>
                <dd className="font-medium">{opc.companyName}</dd>
              </div>
              <div className="flex justify-between gap-4">
                <dt className="text-muted-foreground">统一社会信用代码</dt>
                <dd>{opc.creditCode ?? '—'}</dd>
              </div>
              <div className="flex justify-between gap-4">
                <dt className="text-muted-foreground">状态</dt>
                <dd>
                  <Badge variant={opc.opcStatus === 'active' ? 'success' : 'warning'}>
                    {opc.opcStatus === 'active' ? '已激活' : opc.opcStatus ?? '设立中'}
                  </Badge>
                </dd>
              </div>
              {opc.bankAccount && (
                <div className="flex justify-between gap-4">
                  <dt className="text-muted-foreground">对公账户</dt>
                  <dd>{opc.bankAccount}</dd>
                </div>
              )}
            </dl>
          ) : (
            <p className="text-sm text-muted-foreground">尚未绑定 OPC，请先完成方案签约。</p>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
