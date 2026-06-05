'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { Menu, User } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { ProtectedImage } from '@/components/media/protected-image';
import { maskMemberName, maskMemberPhone } from '@/lib/member/mask';
import { clearMemberTokens } from '@/lib/api/client';

export { maskMemberName, maskMemberPhone } from '@/lib/member/mask';

interface MemberTopBarProps {
  phone?: string | null;
  name?: string | null;
  avatarUrl?: string | null;
  isLoggedIn: boolean;
  onMenuClick?: () => void;
}

export function MemberTopBar({ phone, name, avatarUrl, isLoggedIn, onMenuClick }: MemberTopBarProps) {
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
          {isLoggedIn && (
            <Link
              href="/user/profile"
              className="flex items-center gap-2.5 rounded-md px-2 py-1 transition-colors hover:bg-muted"
              aria-label="个人资料"
            >
              <div className="flex h-9 w-9 shrink-0 items-center justify-center overflow-hidden rounded-full border bg-muted">
                {avatarUrl ? (
                  <ProtectedImage src={avatarUrl} alt="" className="h-full w-full object-cover" />
                ) : (
                  <User className="h-4 w-4 text-muted-foreground" aria-hidden />
                )}
              </div>
              <div className="flex flex-col">
                <span className="text-sm font-medium leading-tight">{maskMemberName(name ?? '')}</span>
                {phone && (
                  <span className="text-xs leading-tight text-muted-foreground">
                    {maskMemberPhone(phone)}
                  </span>
                )}
              </div>
            </Link>
          )}
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
