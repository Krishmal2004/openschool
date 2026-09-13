import { InlineNotification } from "@carbon/react";
import { getErrorMessage } from "../../lib/errorMessage";

interface Props {
  isError: boolean | undefined;
  error?: unknown;
  title?: string;
  fallback?: string;
  onClose?: () => void;
  style?: React.CSSProperties;
}

// Standardizes the mutation-error banner repeated across nearly every admin page/modal.
export default function MutationErrorNotification({
  isError,
  error,
  title = "Error",
  fallback,
  onClose,
  style,
}: Props) {
  if (!isError) return null;
  return (
    <InlineNotification
      kind="error"
      lowContrast
      hideCloseButton={!onClose}
      title={title}
      subtitle={getErrorMessage(error, fallback)}
      onClose={onClose}
      style={{ marginBottom: "1rem", maxWidth: "100%", ...style }}
    />
  );
}
