FROM node:20-alpine AS base
WORKDIR /app
RUN apk add --no-cache libc6-compat
RUN npm config set registry https://registry.npmjs.org/

FROM base AS deps
COPY package.json package-lock.json* ./
RUN npm ci

FROM base AS builder
ARG NEXT_PUBLIC_BASE_PATH=/tax_filing
ENV NEXT_PUBLIC_BASE_PATH=$NEXT_PUBLIC_BASE_PATH
COPY --from=deps /app/node_modules ./node_modules
COPY . .
# Next.js 可选 public 目录；无静态资源时仍需存在，否则 runner 阶段 COPY 失败
RUN mkdir -p public
# 避免本地 .next 缓存污染镜像构建（导致 basePath 未生效）
RUN rm -rf .next
RUN npx prisma generate
RUN npm run build

# 启动时迁移/种子脚本依赖已生成的 Prisma Client（deps 阶段仅有原始 node_modules）
FROM base AS init
COPY --from=deps /app/node_modules ./node_modules
COPY --from=builder /app/prisma ./prisma
COPY --from=builder /app/package.json ./package.json
RUN npx prisma generate

FROM base AS runner
ENV NODE_ENV=production
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

ARG NEXT_PUBLIC_BASE_PATH=/tax_filing
ENV NEXT_PUBLIC_BASE_PATH=$NEXT_PUBLIC_BASE_PATH
ENV PORT=3000
ENV HOSTNAME=0.0.0.0

COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

# 数据库初始化工具（迁移 + 种子脚本，含 prisma generate 产物）
COPY --from=init /app/node_modules /init/node_modules
COPY --from=builder /app/prisma /init/prisma
COPY --from=builder /app/package.json /init/package.json
COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

RUN mkdir -p /app/uploads && chown -R nextjs:nodejs /app/uploads

EXPOSE 3000
ENTRYPOINT ["/entrypoint.sh"]
