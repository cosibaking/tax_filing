'use client';

import Link from 'next/link';
import { useEffect, useState } from 'react';
import { ChevronRight } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { AvatarUpload } from '@/components/member/avatar-upload';
import { OpcDetailModal } from '@/components/member/opc-detail-modal';
import { displayMemberName, displayMemberPhone } from '@/lib/member/mask';
import { ApiClientError, memberFetch, notifyMemberProfileUpdated } from '@/lib/api/client';

interface Profile {
  phone?: string;
  name?: string | null;
  avatarUrl?: string | null;
  avatarFileId?: string | null;
}

interface OpcProfileSummary {
  hasOrder?: boolean;
  companyName?: string | null;
  proposedNamePrimary?: string | null;
  creditCode?: string | null;
  opcStatus?: string | null;
  statusLabel?: string;
  bankAccountMasked?: string;
}

function ProfileRow({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div className="flex justify-between gap-4 text-sm">
      <span className="text-muted-foreground">{label}</span>
      <span className="text-right">{value}</span>
    </div>
  );
}

function opcDisplayName(opc: OpcProfileSummary): string {
  return opc.companyName ?? opc.proposedNamePrimary ?? '设立中';
}

function opcStatusBadge(opc: OpcProfileSummary) {
  const isActive = opc.opcStatus === 'active';
  const label = opc.statusLabel ?? (isActive ? '已激活' : '设立中');
  return <Badge variant={isActive ? 'success' : 'warning'}>{label}</Badge>;
}

export default function ProfilePage() {
  const [profile, setProfile] = useState<Profile | null>(null);
  const [opc, setOpc] = useState<OpcProfileSummary | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [editingName, setEditingName] = useState(false);
  const [nameInput, setNameInput] = useState('');
  const [savingName, setSavingName] = useState(false);
  const [savingAvatar, setSavingAvatar] = useState(false);
  const [saveMsg, setSaveMsg] = useState<string | null>(null);
  const [opcModalOpen, setOpcModalOpen] = useState(false);
  const [showFullAccountInfo, setShowFullAccountInfo] = useState(false);

  const loadData = () =>
    Promise.all([
      memberFetch<Profile>('/api/member/profile'),
      memberFetch<OpcProfileSummary>('/api/member/compliance/opc'),
    ]).then(([p, o]) => {
      setProfile(p);
      setOpc(o);
      setNameInput(p.name ?? '');
    });

  useEffect(() => {
    loadData()
      .catch((err) => {
        setError(err instanceof ApiClientError ? err.message : '加载失败');
      })
      .finally(() => setLoading(false));
  }, []);

  const saveProfile = async (patch: { name?: string; avatarFileId?: string | null }) => {
    const updated = await memberFetch<Profile>('/api/member/profile', {
      method: 'PATCH',
      body: JSON.stringify(patch),
    });
    setProfile(updated);
    setNameInput(updated.name ?? '');
    notifyMemberProfileUpdated();
    return updated;
  };

  const handleSaveName = async () => {
    setSavingName(true);
    setSaveMsg(null);
    try {
      await saveProfile({ name: nameInput });
      setEditingName(false);
      setSaveMsg('姓名已保存');
    } catch (err) {
      setSaveMsg(err instanceof ApiClientError ? err.message : '保存失败');
    } finally {
      setSavingName(false);
    }
  };

  const handleAvatarUploaded = async (fileId: string) => {
    setSavingAvatar(true);
    setSaveMsg(null);
    try {
      await saveProfile({ avatarFileId: fileId });
      setSaveMsg('头像已保存');
    } catch (err) {
      setSaveMsg(err instanceof ApiClientError ? err.message : '头像保存失败');
    } finally {
      setSavingAvatar(false);
    }
  };

  if (loading) return <p className="text-muted-foreground">加载中…</p>;
  if (error) {
    return (
      <Alert variant="destructive">
        <AlertDescription>{error}</AlertDescription>
      </Alert>
    );
  }

  const hasOpcBinding = opc?.opcStatus != null;

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">个人资料</h1>

      {saveMsg && (
        <Alert variant={saveMsg.includes('失败') ? 'destructive' : 'info'}>
          <AlertDescription>{saveMsg}</AlertDescription>
        </Alert>
      )}

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0">
          <CardTitle>账号信息</CardTitle>
          <Button
            variant="outline"
            size="sm"
            onClick={() => setShowFullAccountInfo((v) => !v)}
          >
            {showFullAccountInfo ? '隐藏完整信息' : '查看完整信息'}
          </Button>
        </CardHeader>
        <CardContent className="space-y-6">
          <AvatarUpload
            avatarUrl={profile?.avatarUrl}
            onUploaded={handleAvatarUploaded}
            disabled={savingAvatar}
          />

          <div className="space-y-2 text-sm">
            <p>
              <span className="text-muted-foreground">手机号：</span>
              {displayMemberPhone(profile?.phone, showFullAccountInfo)}
            </p>

            <div className="space-y-2">
              <Label htmlFor="member-name">姓名</Label>
              {editingName ? (
                <div className="flex flex-wrap items-center gap-2">
                  <Input
                    id="member-name"
                    value={nameInput}
                    onChange={(e) => setNameInput(e.target.value)}
                    placeholder="请输入真实姓名"
                    className="max-w-xs"
                    maxLength={64}
                  />
                  <Button size="sm" onClick={handleSaveName} disabled={savingName}>
                    {savingName ? '保存中…' : '保存'}
                  </Button>
                  <Button
                    size="sm"
                    variant="ghost"
                    onClick={() => {
                      setEditingName(false);
                      setNameInput(profile?.name ?? '');
                    }}
                  >
                    取消
                  </Button>
                </div>
              ) : (
                <div className="flex items-center gap-2">
                  <span>{displayMemberName(profile?.name, showFullAccountInfo)}</span>
                  <Button size="sm" variant="outline" onClick={() => setEditingName(true)}>
                    {profile?.name ? '修改' : '设置姓名'}
                  </Button>
                </div>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      <Card
        className={hasOpcBinding ? 'cursor-pointer transition-colors hover:bg-muted/30' : undefined}
        onClick={hasOpcBinding ? () => setOpcModalOpen(true) : undefined}
        role={hasOpcBinding ? 'button' : undefined}
        tabIndex={hasOpcBinding ? 0 : undefined}
        onKeyDown={
          hasOpcBinding
            ? (e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                  e.preventDefault();
                  setOpcModalOpen(true);
                }
              }
            : undefined
        }
      >
        <CardHeader>
          <div className="flex items-center justify-between gap-2">
            <div>
              <CardTitle>OPC 主体绑定</CardTitle>
              <CardDescription>
                {hasOpcBinding
                  ? '您已选择绑定收益主体，点击查看完整信息'
                  : '完成方案签约后，系统将为您创建合规经营主体'}
              </CardDescription>
            </div>
            {hasOpcBinding && <ChevronRight className="h-5 w-5 shrink-0 text-muted-foreground" />}
          </div>
        </CardHeader>
        <CardContent>
          {hasOpcBinding ? (
            <div className="space-y-2">
              <ProfileRow label="公司名称" value={<span className="font-medium">{opcDisplayName(opc!)}</span>} />
              <ProfileRow label="统一社会信用代码" value={opc!.creditCode ?? '—'} />
              <ProfileRow label="状态" value={opcStatusBadge(opc!)} />
              {opc!.bankAccountMasked && (
                <ProfileRow label="对公账户" value={opc!.bankAccountMasked} />
              )}
            </div>
          ) : opc?.hasOrder ? (
            <p className="text-sm text-muted-foreground">
              OPC 设立中，请前往{' '}
              <Link
                href="/user/compliance/opc"
                className="font-medium text-primary underline"
                onClick={(e) => e.stopPropagation()}
              >
                OPC 落地进度
              </Link>{' '}
              查看详情。
            </p>
          ) : (
            <p className="text-sm text-muted-foreground">
              尚未绑定 OPC，请先{' '}
              <Link
                href="/user/compliance/plan"
                className="font-medium text-primary underline"
                onClick={(e) => e.stopPropagation()}
              >
                完成方案签约
              </Link>
              。
            </p>
          )}
        </CardContent>
      </Card>

      <OpcDetailModal open={opcModalOpen} onClose={() => setOpcModalOpen(false)} />
    </div>
  );
}
