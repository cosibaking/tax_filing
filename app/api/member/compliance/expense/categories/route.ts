import { jsonOk, handleApiError } from '@/lib/api/envelope';
import { getExpenseCategories } from '@/lib/services/compliance/ledger/income-expense-service';
import { serializeBigInt } from '@/lib/db';

export async function GET() {
  try {
    const categories = await getExpenseCategories();
    return jsonOk(serializeBigInt(categories));
  } catch (error) {
    return handleApiError(error);
  }
}
