'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { cn } from '@/lib/utils';
import type { MemberNavItem } from '@/lib/portal/member-nav';

interface MemberSidebarProps {
  items: MemberNavItem[];
  onNavigate?: () => void;
}

export function MemberSidebar({ items, onNavigate }: MemberSidebarProps) {
  const pathname = usePathname();
  const groups = Array.from(new Set(items.map((i) => i.group).filter(Boolean)));

  const renderItem = (item: MemberNavItem) => {
    const active = pathname === item.href || pathname.startsWith(`${item.href}/`);
    if (item.disabled) {
      return (
        <span
          key={item.href}
          className="block rounded-md px-3 py-2 text-sm text-muted-foreground/50 cursor-not-allowed"
          title="OPC 激活后可用"
        >
          {item.label}
        </span>
      );
    }
    return (
      <Link
        key={item.href}
        href={item.href}
        onClick={onNavigate}
        className={cn(
          'block rounded-md px-3 py-2 text-sm transition-colors min-h-[44px] flex items-center',
          active ? 'bg-primary/10 font-medium text-primary' : 'text-muted-foreground hover:bg-accent hover:text-foreground',
        )}
      >
        {item.label}
      </Link>
    );
  };

  if (groups.length === 0) {
    return <nav className="space-y-1">{items.map(renderItem)}</nav>;
  }

  return (
    <nav className="space-y-6">
      {groups.map((group) => (
        <div key={group}>
          <p className="mb-2 px-3 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
            {group}
          </p>
          <div className="space-y-1">
            {items.filter((i) => i.group === group).map(renderItem)}
          </div>
        </div>
      ))}
      {items
        .filter((i) => !i.group)
        .map(renderItem)}
    </nav>
  );
}
