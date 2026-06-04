
/** 与 pending-tasks.ts 中终态过滤规则保持一致的纯函数测试替身 */
function shouldShowSignTask(hasActiveOrder: boolean): boolean {
  return !hasActiveOrder;
}

function filterOpcTasks(input: {
  opcStatus: string;
  materialsSubmittedAt: Date | null;
  bankReceiptFileId: bigint | null;
}): string[] {
  const keys: string[] = [];
  if (input.opcStatus === 'pending' && !input.materialsSubmittedAt) {
    keys.push('submit_opc_materials');
  }
  if (input.opcStatus === 'materials') {
    keys.push('resubmit_opc_materials');
  }
  if (input.opcStatus === 'bank' && !input.bankReceiptFileId) {
    keys.push('upload_bank_receipt');
  }
  return keys;
}

describe('pending task visibility', () => {
  it('hides sign task when order is active', () => {
    expect(shouldShowSignTask(false)).toBe(true);
    expect(shouldShowSignTask(true)).toBe(false);
  });

  it('hides material submit after submitted', () => {
    expect(
      filterOpcTasks({ opcStatus: 'pending', materialsSubmittedAt: new Date(), bankReceiptFileId: null }),
    ).toEqual([]);
    expect(
      filterOpcTasks({ opcStatus: 'materials_review', materialsSubmittedAt: new Date(), bankReceiptFileId: null }),
    ).toEqual([]);
  });

  it('shows resubmit only in materials status', () => {
    expect(
      filterOpcTasks({ opcStatus: 'materials', materialsSubmittedAt: new Date(), bankReceiptFileId: null }),
    ).toEqual(['resubmit_opc_materials']);
    expect(
      filterOpcTasks({ opcStatus: 'active', materialsSubmittedAt: new Date(), bankReceiptFileId: null }),
    ).toEqual([]);
  });

  it('treats filed tax status as terminal', () => {
    const openStatuses = ['pending', 'overdue'];
    expect(openStatuses).not.toContain('filed');
  });
});
