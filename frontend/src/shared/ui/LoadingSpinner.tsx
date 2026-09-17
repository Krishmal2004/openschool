import { InlineLoading } from "@carbon/react";

export default function LoadingSpinner() {
  return (
    <div className="os-flex os-items-center os-justify-center os-p-12"
    >
      <InlineLoading description="Loading…" status="active" />
    </div>
  );
}
