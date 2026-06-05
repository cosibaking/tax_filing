'use client';

type ProtectedImageProps = {
  src: string;
  alt: string;
  className?: string;
};

/** 通过签名 URL 展示图片，配合服务端 Referer 防盗链校验 */
export function ProtectedImage({ src, alt, className }: ProtectedImageProps) {
  return (
    // eslint-disable-next-line @next/next/no-img-element
    <img
      src={src}
      alt={alt}
      className={className}
      referrerPolicy="same-origin"
      loading="lazy"
      decoding="async"
    />
  );
}
