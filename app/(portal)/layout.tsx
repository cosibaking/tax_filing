import { PortalFooter } from '@/components/portal/PortalFooter';
import { PortalHeader } from '@/components/portal/PortalHeader';

export default function PortalLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex min-h-screen flex-col">
      <PortalHeader />
      <main className="flex-1">{children}</main>
      <PortalFooter />
    </div>
  );
}
