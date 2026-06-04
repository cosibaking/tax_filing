import { jsonOk, handleApiError } from '@/lib/api/envelope';
import { TAX_CHECKLIST_ITEMS } from '@/lib/api/constants';

export async function GET() {
  try {
    return jsonOk({
      items: TAX_CHECKLIST_ITEMS.map((text, index) => ({
        id: index + 1,
        text,
      })),
    });
  } catch (error) {
    return handleApiError(error);
  }
}
