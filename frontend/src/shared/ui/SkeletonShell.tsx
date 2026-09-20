import { SkeletonText, SkeletonPlaceholder } from "@carbon/react";

export function ContentSkeleton() {
  return (
    <div className="os-p-6">
      <SkeletonPlaceholder className="os-w-full os-h-7 os-mb-4" />
      <SkeletonPlaceholder className="os-w-full os-h-7 os-mb-4" />
      <SkeletonPlaceholder className="os-w-full os-h-7" />
    </div>
  );
}

export default function SkeletonShell() {
  return (
    <>
      <div className="os-skeleton-header" />
      <div className="os-layout">
        <aside className="os-layout__sidebar">
          <div className="os-p-4 os-stack">
            {Array.from({ length: 6 }).map((_, i) => (
              <SkeletonText key={i} width="75%" />
            ))}
          </div>
        </aside>
        <main className="os-layout__content">
          <ContentSkeleton />
        </main>
      </div>
    </>
  );
}
