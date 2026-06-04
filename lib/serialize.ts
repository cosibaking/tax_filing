/** 与数据库无关的序列化工具，可供 API 与客户端安全引用 */
export function serializeBigInt<T>(obj: T): T {
  return JSON.parse(
    JSON.stringify(obj, (_key, value) => (typeof value === 'bigint' ? value.toString() : value)),
  ) as T;
}
