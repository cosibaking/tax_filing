import { Suspense } from 'react';
import { DiagnosisResultContent } from './result-content';

export default function DiagnosisResultPage() {
  return (
    <Suspense fallback={<div className="py-16 text-center text-muted-foreground">加载中…</div>}>
      <DiagnosisResultContent />
    </Suspense>
  );
}
