'use client';

import Link from 'next/link';
import { useCallback, useEffect, useState } from 'react';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { X } from 'lucide-react';
import { MemberSidebar } from '@/components/member/MemberSidebar';
import { MemberTopBar } from '@/components/member/MemberTopBar';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import {
  getMemberGuideMessage,
  getMemberNavItems,
  type MemberNavState,
} from '@/lib/portal/member-nav';
import {
  ApiClientError,
  getMemberToken,
  MEMBER_AUTH_CHANGED_EVENT,
  MEMBER_PROFILE_UPDATED_EVENT,
  memberFetch,
} from '@/lib/api/client';
import {
  buildLoginRedirectUrl,
  isMemberPublicPath,
  MEMBER_LOGIN_PATH,
} from '@/lib/auth/login-redirect';

interface MemberContext {
  hasOrder: boolean;
  opcActive: boolean;
}

export function MemberLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const router = useRouter();
  const [mounted, setMounted] = useState(false);
  const [authReady, setAuthReady] = useState(false);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [phone, setPhone] = useState<string | null>(null);
  const [name, setName] = useState<string | null>(null);
  const [avatarUrl, setAvatarUrl] = useState<string | null>(null);
  const [ctx, setCtx] = useState<MemberContext>({ hasOrder: false, opcActive: false });
  const [loading, setLoading] = useState(true);

  const loadMemberSession = useCallback(() => {
    const token = getMemberToken();
    setIsLoggedIn(!!token);
    if (!token) {
      setPhone(null);
      setName(null);
      setAvatarUrl(null);
      setCtx({ hasOrder: false, opcActive: false });
      setLoading(false);
      return;
    }
    setLoading(true);
    Promise.all([
      memberFetch<{ hasOrder?: boolean; opcStatus?: string } | null>('/api/member/compliance/opc'),
      memberFetch<{ phone?: string; name?: string | null; avatarUrl?: string | null }>(
        '/api/member/profile',
      ),
    ])
      .then(([opcData, profile]) => {
        setCtx({
          hasOrder: !!opcData?.hasOrder,
          opcActive: opcData?.opcStatus === 'active',
        });
        setPhone(profile.phone ?? null);
        setName(profile.name ?? null);
        setAvatarUrl(profile.avatarUrl ?? null);
      })
      .catch((err: unknown) => {
        if (
          err instanceof ApiClientError &&
          (err.status === 401 || err.status === 404)
        ) {
          setCtx({ hasOrder: false, opcActive: false });
          setPhone(null);
          setName(null);
          setAvatarUrl(null);
        }
      })
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    setMounted(true);
  }, []);

  useEffect(() => {
    if (!mounted) return;
    if (isMemberPublicPath(pathname)) {
      setAuthReady(true);
      return;
    }
    if (!getMemberToken()) {
      const query = searchParams.toString();
      const returnTo = query ? `${pathname}?${query}` : pathname;
      router.replace(buildLoginRedirectUrl(MEMBER_LOGIN_PATH, returnTo));
      return;
    }
    setAuthReady(true);
  }, [mounted, pathname, searchParams, router]);

  useEffect(() => {
    loadMemberSession();
  }, [pathname, loadMemberSession]);

  useEffect(() => {
    const onAuthChange = () => loadMemberSession();
    const onProfileUpdate = () => loadMemberSession();
    window.addEventListener(MEMBER_AUTH_CHANGED_EVENT, onAuthChange);
    window.addEventListener(MEMBER_PROFILE_UPDATED_EVENT, onProfileUpdate);
    return () => {
      window.removeEventListener(MEMBER_AUTH_CHANGED_EVENT, onAuthChange);
      window.removeEventListener(MEMBER_PROFILE_UPDATED_EVENT, onProfileUpdate);
    };
  }, [loadMemberSession]);

  const navState: MemberNavState = {
    isLoggedIn,
    hasOrder: ctx.hasOrder,
    opcActive: ctx.opcActive,
    opcPending: ctx.hasOrder && !ctx.opcActive,
  };
  const navItems = getMemberNavItems(navState);
  const guide = mounted ? getMemberGuideMessage(navState) : null;

  return (
    <div className="min-h-screen bg-muted/20">
      <MemberTopBar
        phone={phone}
        name={name}
        avatarUrl={avatarUrl}
        isLoggedIn={isLoggedIn}
        onMenuClick={() => setDrawerOpen(true)}
      />

      {drawerOpen && (
        <div className="fixed inset-0 z-50 lg:hidden">
          <div className="absolute inset-0 bg-black/40" onClick={() => setDrawerOpen(false)} />
          <aside className="absolute left-0 top-0 h-full w-72 bg-background p-4 shadow-lg">
            <div className="mb-4 flex items-center justify-between">
              <span className="font-semibold">菜单</span>
              <Button variant="ghost" size="icon" onClick={() => setDrawerOpen(false)}>
                <X className="h-5 w-5" />
              </Button>
            </div>
            <MemberSidebar items={navItems} onNavigate={() => setDrawerOpen(false)} />
          </aside>
        </div>
      )}

      <div className="mx-auto flex max-w-7xl gap-0 lg:gap-8 lg:p-6">
        <aside className="hidden w-56 shrink-0 lg:block">
          <div className="sticky top-6 rounded-lg border bg-background p-4">
            <p className="mb-4 px-3 text-lg font-semibold">合规中心</p>
            <MemberSidebar items={navItems} />
          </div>
        </aside>
        <main className="min-w-0 flex-1 p-4 lg:p-0">
          {guide && (
            <Alert variant={guide.variant === 'warning' ? 'warning' : 'info'} className="mb-4">
              <AlertDescription>
                {guide.message}{' '}
                {guide.href && (
                  <Link href={guide.href} className="font-medium underline">
                    立即前往
                  </Link>
                )}
              </AlertDescription>
            </Alert>
          )}
          {!mounted || !authReady || loading ? (
            <p className="text-sm text-muted-foreground">加载中…</p>
          ) : (
            children
          )}
        </main>
      </div>
    </div>
  );
}
