import Redis from 'ioredis';

const globalForRedis = globalThis as unknown as { redis: Redis | null };

function createRedisClient(): Redis | null {
  const url = process.env.REDIS_URL;
  if (!url) return null;
  return new Redis(url, { maxRetriesPerRequest: 2, lazyConnect: true });
}

export function getRedis(): Redis {
  if (!globalForRedis.redis) {
    const client = createRedisClient();
    if (!client) {
      throw new Error('REDIS_URL is not configured');
    }
    globalForRedis.redis = client;
  }
  return globalForRedis.redis;
}

export function isRedisConfigured(): boolean {
  return Boolean(process.env.REDIS_URL);
}

if (process.env.NODE_ENV !== 'production') {
  globalForRedis.redis = globalForRedis.redis ?? null;
}
