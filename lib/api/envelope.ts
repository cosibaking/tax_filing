import { NextResponse } from 'next/server';

export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T | null;
}

export function ok<T>(data: T, message = 'success'): ApiResponse<T> {
  return { code: 0, message, data };
}

export function fail(code: number, message: string): ApiResponse<null> {
  return { code, message, data: null };
}

export function jsonOk<T>(data: T, message = 'success', status = 200) {
  return NextResponse.json(ok(data, message), { status });
}

export function jsonFail(code: number, message: string, status = 400) {
  return NextResponse.json(fail(code, message), { status });
}

export function handleApiError(error: unknown) {
  if (error instanceof Error) {
    if (error.message === 'UNAUTHORIZED') return jsonFail(1001, '未授权', 401);
    if (error.message === 'FORBIDDEN') return jsonFail(1003, '无权限', 403);
    if (error.message === 'NOT_FOUND') return jsonFail(1004, '资源不存在', 404);
    return jsonFail(1000, error.message, 400);
  }
  return jsonFail(1000, '服务器错误', 500);
}
