import { Button, ComposedModal, ModalHeader, ModalBody, ModalFooter } from "@carbon/react";

interface Props {
  open: boolean;
  onStay: () => void;
  onLeave: () => void;
}

export default function UnsavedChangesModal({ open, onStay, onLeave }: Props) {
  return (
    <ComposedModal open={open} size="sm" onClose={onStay} aria-label="Leave without saving?">
      <ModalHeader title="Leave without saving?" />
      <ModalBody>
        <p className="os-text-md">You have unsaved changes. If you leave now, they'll be lost.</p>
      </ModalBody>
      <ModalFooter>
        <Button kind="secondary" onClick={onStay}>
          Stay
        </Button>
        <Button kind="danger" onClick={onLeave}>
          Leave without saving
        </Button>
      </ModalFooter>
    </ComposedModal>
  );
}
