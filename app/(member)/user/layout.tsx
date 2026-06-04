import { MemberLayout } from '@/components/member/MemberLayout';

export default function UserLayout({ children }: { children: React.ReactNode }) {
  return <MemberLayout>{children}</MemberLayout>;
}
