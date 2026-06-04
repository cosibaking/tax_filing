import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: '主播OPC合规服务',
  description: '个人主播合规诊断、OPC 设立与记账申报一站式服务',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="zh-CN">
      <body>{children}</body>
    </html>
  );
}
