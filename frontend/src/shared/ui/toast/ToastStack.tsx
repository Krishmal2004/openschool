import { ToastNotification, ActionableNotification } from "@carbon/react";
import { useToast } from "@/shared/ui/toast/useToast";

export default function ToastStack() {
  const { toasts, dismissToast } = useToast();

  if (toasts.length === 0) return null;

  return (
    <div className="os-toast-stack">
      {toasts.map((t) =>
        t.action ? (
          <ActionableNotification
            key={t.id}
            kind={t.kind}
            title={t.title}
            subtitle={t.subtitle}
            role="status"
            lowContrast
            inline={false}
            hasFocus={false}
            actionButtonLabel={t.action.label}
            onActionButtonClick={() => {
              t.action?.onClick();
              dismissToast(t.id);
            }}
            onClose={() => dismissToast(t.id)}
            className="os-toast-stack__item"
          />
        ) : (
          <ToastNotification
            key={t.id}
            kind={t.kind}
            title={t.title}
            subtitle={t.subtitle}
            role="status"
            lowContrast
            onClose={() => dismissToast(t.id)}
            className="os-toast-stack__item"
          />
        ),
      )}
    </div>
  );
}
