import { Button, ComposedModal, ModalHeader, ModalBody, ModalFooter } from "@carbon/react";

interface Props {
  open: boolean;
  onStaySignedIn: () => void;
  onSignOut: () => void;
}

// Shown 60 seconds before an idle sign-out (S12) so a still-present user on
// a shared lab or staffroom PC gets a chance to stay signed in, rather than
// losing unsaved work with no warning.
export default function IdleWarningModal({ open, onStaySignedIn, onSignOut }: Props) {
  return (
    <ComposedModal open={open} size="sm" onClose={onStaySignedIn} preventCloseOnClickOutside>
      <ModalHeader title="Still there?" />
      <ModalBody>
        <p className="os-text-md">
          You&apos;ve been inactive for a while. For your security, you&apos;ll be signed out in less than a minute
          unless you choose to stay signed in.
        </p>
      </ModalBody>
      <ModalFooter>
        <Button kind="secondary" onClick={onSignOut}>
          Sign out now
        </Button>
        <Button onClick={onStaySignedIn}>Stay signed in</Button>
      </ModalFooter>
    </ComposedModal>
  );
}
