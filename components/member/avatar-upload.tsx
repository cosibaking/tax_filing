'use client';

import { useEffect, useRef, useState } from 'react';
import { Button } from '@/components/ui/button';
import { getMemberToken } from '@/lib/api/client';
import { ProtectedImage } from '@/components/media/protected-image';

type Props = {
  avatarUrl?: string | null;
  onUploaded: (fileId: string) => void;
  disabled?: boolean;
};

export function AvatarUpload({ avatarUrl, onUploaded, disabled }: Props) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [localPreview, setLocalPreview] = useState<string | null>(null);

  useEffect(() => {
    if (avatarUrl) setLocalPreview(null);
  }, [avatarUrl]);

  const displayUrl = localPreview ?? avatarUrl ?? null;

  const handleFile = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    if (!file.type.startsWith('image/')) {
      setError('请选择图片文件');
      return;
    }
    if (file.size > 5 * 1024 * 1024) {
      setError('图片大小不能超过 5MB');
      return;
    }

    const token = getMemberToken();
    if (!token) {
      setError('请先登录');
      return;
    }

    setUploading(true);
    setError(null);
    try {
      const form = new FormData();
      form.append('file', file);
      const res = await fetch('/api/upload', {
        method: 'POST',
        headers: { Authorization: `Bearer ${token}` },
        body: form,
      });
      const json = await res.json();
      if (json.code !== 0) throw new Error(json.message || '上传失败');
      setLocalPreview(URL.createObjectURL(file));
      onUploaded(String(json.data.id));
    } catch (err) {
      setError(err instanceof Error ? err.message : '上传失败');
    } finally {
      setUploading(false);
      e.target.value = '';
    }
  };

  return (
    <div className="flex items-center gap-4">
      <div className="flex h-16 w-16 shrink-0 items-center justify-center overflow-hidden rounded-full border bg-muted">
        {displayUrl ? (
          <ProtectedImage src={displayUrl} alt="头像" className="h-full w-full object-cover" />
        ) : (
          <span className="text-xs text-muted-foreground">无头像</span>
        )}
      </div>
      <div className="space-y-1">
        <Button
          type="button"
          variant="outline"
          size="sm"
          disabled={disabled || uploading}
          onClick={() => inputRef.current?.click()}
        >
          {uploading ? '上传中…' : '更换头像'}
        </Button>
        <input
          ref={inputRef}
          type="file"
          className="hidden"
          accept="image/jpeg,image/png,image/webp,image/gif"
          onChange={handleFile}
        />
        <p className="text-xs text-muted-foreground">支持 JPG/PNG/WebP，最大 5MB</p>
        {error && <p className="text-xs text-destructive">{error}</p>}
      </div>
    </div>
  );
}
