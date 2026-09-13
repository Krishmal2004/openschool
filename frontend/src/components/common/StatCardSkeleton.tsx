import { SkeletonText } from "@carbon/react";

// Placeholder for one `.os-stat-card` while its value is loading.
export default function StatCardSkeleton() {
  return (
    <div className="os-stat-card">
      <div style={{ marginBottom: "0.5rem" }}>
        <SkeletonText width="60%" />
      </div>
      <SkeletonText width="30%" heading />
    </div>
  );
}
