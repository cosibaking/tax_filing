'use client';

import Link from 'next/link';
import { AdminSidebar } from '@/components/admin/AdminSidebar';
import { Button } from '@/components/ui/button';

export function AdminLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen bg-muted/20">
      <header className="border-b bg-background">
        <div className="mx-auto flex h-14 max-w-7xl items-center justify-between px-4">
          <Link href="/admin/compliance/customers" className="font-semibold text-primary">
            顾问工作台
          </Link>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => {
              localStorage.removeItem('admin_token');
              window.location.href = '/admin/login';
            }}
          >
            退出登录
          </Button>
        </div>
      </header>
      <div className="mx-auto flex max-w-7xl gap-6 p-4 lg:p-6">
        <aside className="hidden w-52 shrink-0 rounded-lg border bg-background p-4 md:block">
          <AdminSidebar />
        </aside>
        <main className="min-w-0 flex-1">{children}</main>
      </div>
    </div>
  );
}
