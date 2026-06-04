export interface MemberNavState {
  isLoggedIn: boolean;
  hasOrder: boolean;
  opcActive: boolean;
  opcPending: boolean;
}

export interface MemberNavItem {
  label: string;
  href: string;
  disabled?: boolean;
  group?: string;
}

export function getMemberNavItems(state: MemberNavState): MemberNavItem[] {
  const items: MemberNavItem[] = [
    { label: '概览', href: '/user/overview', group: '合规服务' },
    { label: '诊断历史', href: '/user/compliance/diagnosis', group: '合规服务' },
  ];

  if (!state.hasOrder || !state.opcActive) {
    if (!state.hasOrder) {
      items.push({
        label: '方案与签约',
        href: '/user/compliance/plan',
        group: '合规服务',
      });
    }
    if (state.hasOrder && !state.opcActive) {
      items.push({
        label: 'OPC 进度',
        href: '/user/compliance/opc',
        group: '合规服务',
      });
    }
  }

  const ledgerDisabled = state.hasOrder && !state.opcActive;
  const ledgerItems: MemberNavItem[] = [
    { label: '收入台账', href: '/user/compliance/income', disabled: ledgerDisabled },
    { label: '费用台账', href: '/user/compliance/expense', disabled: ledgerDisabled },
    { label: '利润表', href: '/user/compliance/ledger', disabled: ledgerDisabled },
    { label: '申报日历', href: '/user/compliance/tax', disabled: ledgerDisabled },
    { label: '月度对账单', href: '/user/compliance/statement', disabled: ledgerDisabled },
  ];

  if (state.opcActive) {
    ledgerItems.forEach((item) => {
      items.push({ ...item, group: '记账申报' });
    });
  } else if (state.hasOrder) {
    ledgerItems.forEach((item) => {
      items.push({ ...item, group: '记账申报' });
    });
  }

  items.push(
    { label: '个人资料', href: '/user/profile', group: '账户' },
    { label: '通知中心', href: '/user/notification', group: '账户' },
  );

  return items;
}

export function getMemberGuideMessage(state: MemberNavState): {
  variant: 'warning' | 'info' | 'destructive';
  message: string;
  href?: string;
} | null {
  if (!state.isLoggedIn) return null;
  if (!state.hasOrder) {
    return {
      variant: 'warning',
      message: '您尚未签约合规服务，请先选择方案并完成签约。',
      href: '/user/compliance/plan',
    };
  }
  if (state.hasOrder && !state.opcActive) {
    return {
      variant: 'info',
      message: 'OPC 设立进行中，台账功能将在 OPC 激活后开放。',
      href: '/user/compliance/opc',
    };
  }
  return null;
}
