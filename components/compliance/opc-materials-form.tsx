'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Checkbox } from '@/components/ui/checkbox';
import { OpcFileUpload } from '@/components/compliance/opc-file-upload';
import { BUSINESS_SCOPE_TEMPLATE } from '@/lib/services/compliance/opc/opc-constants';
import {
  parseMaterialsBody,
  validateOpcMaterials,
} from '@/lib/services/compliance/opc/opc-materials-validation';
import type { OpcMaterialsInput } from '@/lib/services/compliance/opc/opc-types';

type Props = {
  onSubmit: (data: OpcMaterialsInput) => Promise<void>;
  submitting: boolean;
  onValidationError?: (message: string) => void;
};

const emptyForm = (): OpcMaterialsInput => ({
  proposedNames: [''],
  registeredCapital: 10,
  capitalTermYears: 20,
  businessTermType: 'long_term',
  businessScope: '',
  registerProvince: '',
  registerCity: '',
  registerDistrict: '',
  registerAddress: '',
  addressProofFileId: '',
  legalPersonName: '',
  idCard: '',
  idCardValidFrom: '',
  idCardValidTo: '',
  householdAddress: '',
  residentialAddress: '',
  phone: '',
  email: '',
  idCardFrontFileId: '',
  idCardBackFileId: '',
  esignAuthorized: false,
  confirmations: { truthful: false, usageConsent: false, opcLimitAck: false },
});

export function OpcMaterialsForm({ onSubmit, submitting, onValidationError }: Props) {
  const [form, setForm] = useState<OpcMaterialsInput>(emptyForm);

  const set = <K extends keyof OpcMaterialsInput>(key: K, value: OpcMaterialsInput[K]) => {
    setForm((f) => ({ ...f, [key]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const names = form.proposedNames.map((n) => n.trim()).filter(Boolean);
    const payload = parseMaterialsBody({ ...form, proposedNames: names });
    const err = validateOpcMaterials(payload);
    if (err) {
      onValidationError?.(err);
      return;
    }
    await onSubmit(payload);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-8">
      <section className="space-y-4">
        <h3 className="font-semibold">拟设公司信息</h3>
        {form.proposedNames.map((name, i) => (
          <div key={i}>
            <Label>备选公司名称 {i + 1}</Label>
            <Input
              className="mt-1"
              value={name}
              onChange={(e) => {
                const next = [...form.proposedNames];
                next[i] = e.target.value;
                set('proposedNames', next);
              }}
              placeholder="如：杭州星播文化传媒有限公司"
            />
          </div>
        ))}
        {form.proposedNames.length < 3 && (
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => set('proposedNames', [...form.proposedNames, ''])}
          >
            添加备选名称
          </Button>
        )}
        <div className="grid gap-4 sm:grid-cols-2">
          <div>
            <Label>注册资本（万元）</Label>
            <Input
              type="number"
              className="mt-1"
              min={1}
              max={1000}
              value={form.registeredCapital}
              onChange={(e) => set('registeredCapital', Number(e.target.value))}
            />
          </div>
          <div>
            <Label>认缴期限</Label>
            <select
              className="mt-1 flex h-10 w-full rounded-md border border-input bg-background px-3 text-sm"
              value={form.capitalTermYears}
              onChange={(e) => set('capitalTermYears', Number(e.target.value))}
            >
              <option value={5}>5 年</option>
              <option value={10}>10 年</option>
              <option value={20}>20 年</option>
              <option value={30}>30 年</option>
            </select>
          </div>
        </div>
        <div>
          <Label>营业期限</Label>
          <div className="mt-2 flex flex-wrap gap-4">
            <label className="flex items-center gap-2 text-sm">
              <input
                type="radio"
                checked={form.businessTermType === 'long_term'}
                onChange={() => set('businessTermType', 'long_term')}
              />
              长期
            </label>
            <label className="flex items-center gap-2 text-sm">
              <input
                type="radio"
                checked={form.businessTermType === 'fixed'}
                onChange={() => set('businessTermType', 'fixed')}
              />
              固定期限
            </label>
          </div>
          {form.businessTermType === 'fixed' && (
            <Input
              type="date"
              className="mt-2"
              value={form.businessTermEnd ?? ''}
              onChange={(e) => set('businessTermEnd', e.target.value)}
            />
          )}
        </div>
        <div>
          <div className="flex items-center justify-between">
            <Label>经营范围</Label>
            <Button
              type="button"
              variant="link"
              className="h-auto p-0 text-xs"
              onClick={() => set('businessScope', BUSINESS_SCOPE_TEMPLATE)}
            >
              使用推荐模板
            </Button>
          </div>
          <textarea
            className="mt-1 flex min-h-[100px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
            value={form.businessScope}
            onChange={(e) => set('businessScope', e.target.value)}
          />
        </div>
        <div className="grid gap-4 sm:grid-cols-3">
          <div>
            <Label>省</Label>
            <Input className="mt-1" value={form.registerProvince} onChange={(e) => set('registerProvince', e.target.value)} />
          </div>
          <div>
            <Label>市</Label>
            <Input className="mt-1" value={form.registerCity} onChange={(e) => set('registerCity', e.target.value)} />
          </div>
          <div>
            <Label>区</Label>
            <Input className="mt-1" value={form.registerDistrict} onChange={(e) => set('registerDistrict', e.target.value)} />
          </div>
        </div>
        <div>
          <Label>详细注册地址</Label>
          <Input className="mt-1" value={form.registerAddress} onChange={(e) => set('registerAddress', e.target.value)} />
        </div>
        <OpcFileUpload
          label="地址证明（租赁合同/产权证/园区证明）"
          value={form.addressProofFileId}
          onChange={(id) => set('addressProofFileId', id)}
        />
      </section>

      <section className="space-y-4">
        <h3 className="font-semibold">法人信息</h3>
        <div>
          <Label>法人姓名</Label>
          <Input className="mt-1" value={form.legalPersonName} onChange={(e) => set('legalPersonName', e.target.value)} />
        </div>
        <div>
          <Label>身份证号</Label>
          <Input className="mt-1" value={form.idCard} onChange={(e) => set('idCard', e.target.value)} />
        </div>
        <div className="grid gap-4 sm:grid-cols-2">
          <div>
            <Label>证件有效期起</Label>
            <Input type="date" className="mt-1" value={form.idCardValidFrom} onChange={(e) => set('idCardValidFrom', e.target.value)} />
          </div>
          <div>
            <Label>证件有效期止</Label>
            <Input type="date" className="mt-1" value={form.idCardValidTo} onChange={(e) => set('idCardValidTo', e.target.value)} />
          </div>
        </div>
        <div>
          <Label>户籍地址</Label>
          <Input className="mt-1" value={form.householdAddress} onChange={(e) => set('householdAddress', e.target.value)} />
        </div>
        <div>
          <Label>现居住地址</Label>
          <Input className="mt-1" value={form.residentialAddress} onChange={(e) => set('residentialAddress', e.target.value)} />
        </div>
        <div className="grid gap-4 sm:grid-cols-2">
          <div>
            <Label>手机号</Label>
            <Input className="mt-1" value={form.phone} onChange={(e) => set('phone', e.target.value)} />
          </div>
          <div>
            <Label>邮箱</Label>
            <Input className="mt-1" type="email" value={form.email} onChange={(e) => set('email', e.target.value)} />
          </div>
        </div>
        <OpcFileUpload label="身份证正面" value={form.idCardFrontFileId} onChange={(id) => set('idCardFrontFileId', id)} />
        <OpcFileUpload label="身份证反面" value={form.idCardBackFileId} onChange={(id) => set('idCardBackFileId', id)} />
      </section>

      <section className="space-y-3">
        <h3 className="font-semibold">确认事项</h3>
        {(
          [
            ['truthful', '本人确认以上信息真实、完整'],
            ['usageConsent', '同意将信息用于工商、税务、银行开户申报'],
            ['opcLimitAck', '知晓自然人 3 年内不得再设立新的 OPC'],
          ] as const
        ).map(([key, text]) => (
          <label key={key} className="flex items-start gap-2 text-sm">
            <Checkbox
              checked={form.confirmations[key]}
              onCheckedChange={(v) =>
                set('confirmations', { ...form.confirmations, [key]: v === true })
              }
            />
            <span>{text}</span>
          </label>
        ))}
        <label className="flex items-start gap-2 text-sm">
          <Checkbox
            checked={form.esignAuthorized}
            onCheckedChange={(v) => set('esignAuthorized', v === true)}
          />
          <span>授权使用签约时的电子签名办理工商登记</span>
        </label>
      </section>

      <Button type="submit" disabled={submitting}>
        {submitting ? '提交中…' : '提交注册资料'}
      </Button>
    </form>
  );
}
