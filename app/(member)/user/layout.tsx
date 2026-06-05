import { Suspense } from 'react';
import { MemberLayout } from '@/components/member/MemberLayout';

export default function UserLayout({ children }: { children: React.ReactNode }) {
  return (
    <Suspense fallback={<p className="p-4 text-sm text-muted-foreground">加载中…</p>}>
      <MemberLayout>{children}</MemberLayout>
    </Suspense>
  );
}
