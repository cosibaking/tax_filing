export type OpcMaterialsInput = {
  proposedNames: string[];
  registeredCapital: number;
  capitalTermYears: number;
  businessTermType: 'long_term' | 'fixed';
  businessTermEnd?: string;
  businessScope: string;
  registerProvince: string;
  registerCity: string;
  registerDistrict: string;
  registerAddress: string;
  addressProofFileId: string;
  legalPersonName: string;
  idCard: string;
  idCardValidFrom: string;
  idCardValidTo: string;
  householdAddress: string;
  residentialAddress: string;
  phone: string;
  email: string;
  idCardFrontFileId: string;
  idCardBackFileId: string;
  ethnicity?: string;
  esignAuthorized: boolean;
  confirmations: {
    truthful: boolean;
    usageConsent: boolean;
    opcLimitAck: boolean;
  };
};

export type OpcAdminAction =
  | 'approve_materials'
  | 'reject_materials'
  | 'issue_license'
  | 'complete_tax'
  | 'complete_bank';
