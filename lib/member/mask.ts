/** 手机号脱敏：13079678505 → 130****8505 */
export function maskMemberPhone(phone: string): string {
  if (phone.length < 7) return phone;
  return `${phone.slice(0, 3)}****${phone.slice(-4)}`;
}

/** 姓名脱敏：张三→张*，张三丰→张*丰，欧阳娜娜→欧**娜 */
export function maskMemberName(name: string): string {
  const trimmed = name.trim();
  if (!trimmed) return '未设置';
  if (trimmed.length === 1) return '*';
  if (trimmed.length === 2) return `${trimmed[0]}*`;
  return `${trimmed[0]}${'*'.repeat(trimmed.length - 2)}${trimmed[trimmed.length - 1]}`;
}

export function displayMemberPhone(phone: string | undefined | null, showFull: boolean): string {
  if (!phone) return '—';
  return showFull ? phone : maskMemberPhone(phone);
}

export function displayMemberName(name: string | null | undefined, showFull: boolean): string {
  if (!name?.trim()) return '未设置';
  return showFull ? name : maskMemberName(name);
}
