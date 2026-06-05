#!/bin/sh
set -e

INIT_DIR="/init"
DEPLOY_MODE="${DEPLOY_MODE:-prod}"

echo "[entrypoint] deploy mode: ${DEPLOY_MODE}"

cd "${INIT_DIR}"

echo "[entrypoint] running database migrations..."
npx prisma migrate deploy

echo "[entrypoint] seeding base data (admin, plans, categories)..."
npx tsx prisma/seed.ts

if [ "${DEPLOY_MODE}" = "dev" ]; then
  echo "[entrypoint] dev mode: generating random test data..."
  SEED_RANDOM_SKIP=false npx tsx prisma/seed-random.ts
else
  echo "[entrypoint] prod mode: skip random test data"
fi

cd /app
echo "[entrypoint] starting application..."
exec su -s /bin/sh nextjs -c "node server.js"
