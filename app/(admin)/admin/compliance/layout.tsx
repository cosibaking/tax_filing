import { Suspense } from 'react';
import { AdminLayout } from '@/components/admin/AdminLayout';

export default function AdminComplianceLayout({ children }: { children: React.ReactNode }) {
  return (
    <Suspense fallback={<p className="p-4 text-sm text-muted-foreground">加载中…</p>}>
      <AdminLayout>{children}</AdminLayout>
    </Suspense>
  );
}
