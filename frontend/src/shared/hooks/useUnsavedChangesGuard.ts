import { useEffect, useRef, useState } from "react";

// Guards explicit in-page exits (Cancel/Back buttons, sidebar nav intercepted by the
// caller) and the browser's own close/refresh/back. Does not catch every possible
// navigation - react-router's useBlocker needs a data router, which this app's plain
// BrowserRouter doesn't set up.
export function useUnsavedChangesGuard(hasUnsaved: boolean) {
  const [pending, setPending] = useState(false);
  const pendingActionRef = useRef<(() => void) | null>(null);

  useEffect(() => {
    if (!hasUnsaved) return;
    const handler = (e: BeforeUnloadEvent) => {
      e.preventDefault();
      e.returnValue = "";
    };
    window.addEventListener("beforeunload", handler);
    return () => window.removeEventListener("beforeunload", handler);
  }, [hasUnsaved]);

  const guard = (action: () => void) => {
    if (hasUnsaved) {
      pendingActionRef.current = action;
      setPending(true);
    } else {
      action();
    }
  };

  const confirmLeave = () => {
    setPending(false);
    pendingActionRef.current?.();
    pendingActionRef.current = null;
  };

  const cancelLeave = () => {
    setPending(false);
    pendingActionRef.current = null;
  };

  return { guard, modalOpen: pending, confirmLeave, cancelLeave };
}
