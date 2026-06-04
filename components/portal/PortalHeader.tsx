'use client';

import Link from 'next/link';
import { useEffect, useState } from 'react';
import { Button } from '@/components/ui/button';
import { getMemberToken } from '@/lib/api/client';

const navLinks = [
  { href: '/', label: '首页' },
  { href: '/diagnosis', label: '免费诊断' },
  { href: '/pricing', label: '服务价格' },
];

export function PortalHeader() {
  const [loggedIn, setLoggedIn] = useState(false);

  useEffect(() => {
    setLoggedIn(!!getMemberToken());
  }, []);

  return (
    <header className="sticky top-0 z-50 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="mx-auto flex h-14 max-w-6xl items-center justify-between gap-4 px-4 sm:h-16 sm:px-6">
        <Link href="/" className="text-lg font-bold text-primary sm:text-xl">
          主播OPC合规
        </Link>
        <nav className="hidden items-center gap-6 text-sm md:flex">
          {navLinks.map((link) => (
            <Link key={link.href} href={link.href} className="text-muted-foreground hover:text-foreground">
              {link.label}
            </Link>
          ))}
        </nav>
        <div className="flex items-center gap-2">
          {loggedIn ? (
            <Button asChild size="sm">
              <Link href="/user/overview">会员中心</Link>
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
      <nav className="flex gap-4 overflow-x-auto border-t px-4 py-2 text-sm md:hidden">
        {navLinks.map((link) => (
          <Link key={link.href} href={link.href} className="whitespace-nowrap text-muted-foreground">
            {link.label}
          </Link>
        ))}
      </nav>
    </header>
  );
}
