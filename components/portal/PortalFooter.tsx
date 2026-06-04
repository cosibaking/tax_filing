import Link from 'next/link';

export function PortalFooter() {
  return (
    <footer className="border-t bg-muted/40">
      <div className="mx-auto max-w-6xl px-4 py-8 sm:px-6">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <p className="text-sm text-muted-foreground">
            © {new Date().getFullYear()} 主播OPC合规服务 · 合规方案，非逃税方案
          </p>
          <div className="flex flex-wrap gap-4 text-sm">
            <Link href="/legal/privacy" className="text-muted-foreground hover:text-foreground">
              隐私政策
            </Link>
            <Link href="/legal/terms" className="text-muted-foreground hover:text-foreground">
              用户协议
            </Link>
            <Link href="/pricing" className="text-muted-foreground hover:text-foreground">
              服务价格
            </Link>
          </div>
        </div>
      </div>
    </footer>
  );
}
