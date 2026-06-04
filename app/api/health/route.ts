import { jsonOk, handleApiError } from '@/lib/api/envelope';

export async function GET() {
  try {
    return jsonOk({
      status: 'ok',
      timestamp: new Date().toISOString(),
    });
  } catch (error) {
    return handleApiError(error);
  }
}
