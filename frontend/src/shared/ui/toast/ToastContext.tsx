import { useCallback, useMemo, useRef, useState, type ReactNode } from "react";
import { ToastContext, type Toast } from "@/shared/ui/toast/toastContext";

const TOAST_TIMEOUT_MS = 5000;

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);
  const nextId = useRef(0);

  const dismissToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const showToast = useCallback(
    (toast: Omit<Toast, "id">) => {
      const id = String(nextId.current++);
      setToasts((prev) => [...prev, { ...toast, id }]);
      if (!toast.action) {
        setTimeout(() => dismissToast(id), TOAST_TIMEOUT_MS);
      }
    },
    [dismissToast],
  );

  const value = useMemo(() => ({ toasts, showToast, dismissToast }), [toasts, showToast, dismissToast]);

  return <ToastContext.Provider value={value}>{children}</ToastContext.Provider>;
}
