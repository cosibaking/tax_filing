'use client';

import Link from 'next/link';
import { useEffect, useState } from 'react';
import { Menu, X } from 'lucide-react';
import { MemberSidebar } from '@/components/member/MemberSidebar';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import {
  getMemberGuideMessage,
  getMemberNavItems,
  type MemberNavState,
} from '@/lib/portal/member-nav';
import { ApiClientError, getMemberToken, memberFetch } from '@/lib/api/client';

interface MemberContext {
  hasOrder: boolean;
  opcActive: boolean;
}

export function MemberLayout({ children }: { children: React.ReactNode }) {
  const [mounted, setMounted] = useState(false);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [ctx, setCtx] = useState<MemberContext>({ hasOrder: false, opcActive: false });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setMounted(true);
    const token = getMemberToken();
    setIsLoggedIn(!!token);
    if (!token) {
      setLoading(false);
      return;
    }
    memberFetch<{ hasOrder?: boolean; opcStatus?: string } | null>('/api/member/compliance/opc')
      .then((data) => {
        setCtx({
          hasOrder: !!data?.hasOrder,
          opcActive: data?.opcStatus === 'active',
        });
      })
      .catch((err: unknown) => {
        if (
          err instanceof ApiClientError &&
          (err.status === 401 || err.status === 404)
        ) {
          setCtx({ hasOrder: false, opcActive: false });
        }
      })
      .finally(() => setLoading(false));
  }, []);

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
      <header className="sticky top-0 z-40 flex h-14 items-center gap-4 border-b bg-background px-4 lg:hidden">
        <Button variant="ghost" size="icon" onClick={() => setDrawerOpen(true)} aria-label="打开菜单">
          <Menu className="h-5 w-5" />
        </Button>
        <Link href="/user/overview" className="font-semibold text-primary">
          会员中心
        </Link>
      </header>

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
          {!mounted || loading ? (
            <p className="text-sm text-muted-foreground">加载中…</p>
          ) : (
            children
          )}
        </main>
      </div>
    </div>
  );
}
