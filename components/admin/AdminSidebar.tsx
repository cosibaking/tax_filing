'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { cn } from '@/lib/utils';

const adminNav = [
  { href: '/admin/compliance/customers', label: '合规客户' },
  { href: '/admin/compliance/opc-tasks', label: 'OPC 任务' },
  { href: '/admin/compliance/filing', label: '申报工作台' },
  { href: '/admin/compliance/statements', label: '对账单管理' },
];

export function AdminSidebar() {
  const pathname = usePathname();

  return (
    <nav className="space-y-1">
      {adminNav.map((item) => {
        const active = pathname.startsWith(item.href);
        return (
          <Link
            key={item.href}
            href={item.href}
            className={cn(
              'block rounded-md px-3 py-2 text-sm min-h-[44px] flex items-center',
              active ? 'bg-primary text-primary-foreground' : 'text-muted-foreground hover:bg-accent',
            )}
          >
            {item.label}
          </Link>
        );
      })}
    </nav>
  );
}
