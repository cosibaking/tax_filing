/** Next.js basePath，Docker 部署时为 /tax_filing */
export const BASE_PATH = process.env.NEXT_PUBLIC_BASE_PATH ?? '';

export function withBasePath(path: string): string {
  if (!path.startsWith('/') || !BASE_PATH) return path;
  if (path === BASE_PATH || path.startsWith(`${BASE_PATH}/`)) return path;
  return `${BASE_PATH}${path}`;
}

export function stripBasePath(pathname: string): string {
  if (!BASE_PATH || !pathname.startsWith(BASE_PATH)) return pathname;
  const stripped = pathname.slice(BASE_PATH.length);
  return stripped || '/';
}
