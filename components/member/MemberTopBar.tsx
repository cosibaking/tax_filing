'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { Menu } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { clearMemberTokens } from '@/lib/api/client';

export function maskMemberPhone(phone: string): string {
  if (phone.length < 7) return phone;
  return `${phone.slice(0, 3)}****${phone.slice(-4)}`;
}

interface MemberTopBarProps {
  phone?: string | null;
  isLoggedIn: boolean;
  onMenuClick?: () => void;
}

export function MemberTopBar({ phone, isLoggedIn, onMenuClick }: MemberTopBarProps) {
  const router = useRouter();

  const handleLogout = () => {
    clearMemberTokens();
    router.push('/user/login');
  };

  return (
    <header className="sticky top-0 z-50 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="mx-auto flex h-14 max-w-7xl items-center justify-between gap-4 px-4 sm:h-16 sm:px-6">
        <div className="flex min-w-0 items-center gap-2">
          {onMenuClick && (
            <Button
              variant="ghost"
              size="icon"
              className="shrink-0 lg:hidden"
              onClick={onMenuClick}
              aria-label="打开菜单"
            >
              <Menu className="h-5 w-5" />
            </Button>
          )}
          <Link href="/" className="truncate text-lg font-bold text-primary sm:text-xl">
            主播OPC合规
          </Link>
        </div>
        <div className="flex shrink-0 items-center gap-2 sm:gap-3">
          {isLoggedIn && phone ? (
            <span className="text-sm text-muted-foreground">{maskMemberPhone(phone)}</span>
          ) : null}
          {isLoggedIn ? (
            <Button variant="ghost" size="sm" onClick={handleLogout}>
              退出登录
            </Button>
          ) : (
            <>
              <Button asChild variant="ghost" size="sm">
                <Link href="/user/login">登录</Link>
              </Button>
              <Button asChild size="sm">
                <Link href="/user/register">注册</Link>
              </Button>
            </>
          )}
        </div>
      </div>
    </header>
  );
}
