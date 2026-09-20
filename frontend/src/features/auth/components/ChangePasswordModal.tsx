import { useState } from "react";
import { Button, InlineNotification, ComposedModal, ModalHeader, ModalBody, ModalFooter } from "@carbon/react";
import { useChangePassword } from "@/features/auth/queries/useAuth";
import { validateNewPassword } from "@/shared/auth/password";
import MutationErrorNotification from "@/shared/ui/MutationErrorNotification";
import PasswordFields from "@/shared/ui/PasswordFields";

export default function ChangePasswordModal({ onClose }: { onClose: () => void }) {
  const changePassword = useChangePassword();
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [done, setDone] = useState(false);
  const { valid } = validateNewPassword(password, confirm);

  const handleSubmit = () => {
    if (valid) changePassword.mutate(password, { onSuccess: () => setDone(true) });
  };

  return (
    <ComposedModal open size="sm" onClose={onClose} aria-label="Change password">
      <ModalHeader title="Change password" />
      <ModalBody>
        {done ? (
          <InlineNotification
            kind="success"
            title="Password updated"
            subtitle="Use your new password the next time you sign in."
            lowContrast
            hideCloseButton
            className="os-full-width"
          />
        ) : (
          <>
            <MutationErrorNotification
              isError={changePassword.isError}
              error={changePassword.error}
              title="Could not update password"
              fallback="Please try again."
            />
            <PasswordFields
              idPrefix="change-password"
              password={password}
              confirm={confirm}
              onPasswordChange={setPassword}
              onConfirmChange={setConfirm}
            />
          </>
        )}
      </ModalBody>
      <ModalFooter>
        {done ? (
          <Button kind="primary" onClick={onClose}>Done</Button>
        ) : (
          <>
            <Button kind="secondary" onClick={onClose}>Cancel</Button>
            <Button kind="primary" onClick={handleSubmit} disabled={!valid || changePassword.isPending}>
              {changePassword.isPending ? "Saving…" : "Save"}
            </Button>
          </>
        )}
      </ModalFooter>
    </ComposedModal>
  );
}
