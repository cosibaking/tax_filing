'use client';

import { useRef, useState } from 'react';
import { Button } from '@/components/ui/button';
import { getMemberToken } from '@/lib/api/client';

type Props = {
  label: string;
  value?: string;
  onChange: (fileId: string) => void;
};

export function OpcFileUpload({ label, value, onChange }: Props) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleFile = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    const token = getMemberToken();
    if (!token) {
      setError('请先登录后再上传附件');
      return;
    }
    setUploading(true);
    setError(null);
    try {
      const form = new FormData();
      form.append('file', file);
      const res = await fetch('/api/upload', {
        method: 'POST',
        headers: token ? { Authorization: `Bearer ${token}` } : {},
        body: form,
      });
      const json = await res.json();
      if (json.code !== 0) throw new Error(json.message || '上传失败');
      onChange(String(json.data.id));
    } catch (err) {
      setError(err instanceof Error ? err.message : '上传失败');
    } finally {
      setUploading(false);
      e.target.value = '';
    }
  };

  return (
    <div className="space-y-1">
      <span className="text-sm font-medium">{label}</span>
      <div className="flex items-center gap-2">
        <Button
          type="button"
          variant="outline"
          size="sm"
          disabled={uploading}
          onClick={() => inputRef.current?.click()}
        >
          {uploading ? '上传中…' : value ? '重新上传' : '选择文件'}
        </Button>
        <input
          ref={inputRef}
          type="file"
          className="hidden"
          accept="image/*,.pdf"
          onChange={handleFile}
        />
        {value && (
          <span className="text-xs text-green-600 dark:text-green-400">已上传 #{value}</span>
        )}
      </div>
      {error && <p className="text-xs text-destructive">{error}</p>}
    </div>
  );
}
