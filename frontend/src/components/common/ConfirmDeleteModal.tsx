import {
  Button,
  ComposedModal,
  ModalHeader,
  ModalBody,
  ModalFooter,
} from "@carbon/react";

interface Props {
  open: boolean;
  title: string;
  description: React.ReactNode;
  isPending?: boolean;
  // Disables confirming for a reason unrelated to the mutation itself (e.g. a dependent check still loading), without implying one is in flight.
  disabled?: boolean;
  confirmLabel?: string;
  pendingLabel?: string;
  onClose: () => void;
  onConfirm: () => void;
}

export default function ConfirmDeleteModal({
  open,
  title,
  description,
  isPending,
  disabled,
  confirmLabel = "Delete",
  pendingLabel = "Deleting…",
  onClose,
  onConfirm,
}: Props) {
  return (
    <ComposedModal open={open} size="sm" onClose={onClose}>
      <ModalHeader title={title} />
      <ModalBody>
        <p style={{ fontSize: "0.875rem" }}>{description}</p>
      </ModalBody>
      <ModalFooter>
        <Button kind="secondary" onClick={onClose}>
          Cancel
        </Button>
        <Button kind="danger" onClick={onConfirm} disabled={isPending || disabled}>
          {isPending ? pendingLabel : confirmLabel}
        </Button>
      </ModalFooter>
    </ComposedModal>
  );
}
