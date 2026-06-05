import type { NextConfig } from 'next';

const basePath = process.env.NEXT_PUBLIC_BASE_PATH ?? '';

const nextConfig: NextConfig = {
  basePath: basePath || undefined,
  // 与 basePath 保持一致，确保 HTML 中静态资源 URL 带 /tax_filing 前缀
  assetPrefix: basePath || undefined,
  output: 'standalone',
  serverExternalPackages: ['pdfkit'],
  experimental: {
    serverActions: {
      bodySizeLimit: '10mb',
    },
  },
};

export default nextConfig;
