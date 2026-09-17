import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router";
import { ThunderIDProvider } from "@thunderid/react";
import { QueryClientProvider } from "@tanstack/react-query";
import "@/shared/styles/index.scss";
import App from "@/app/App";
import { queryClient } from "@/app/queryClient";
import ErrorBoundary from "@/shared/ui/ErrorBoundary";

const origin = `${window.location.origin}/`;

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ErrorBoundary>
      <ThunderIDProvider
        clientId={import.meta.env.VITE_THUNDERID_CLIENT_ID}
        baseUrl={import.meta.env.VITE_THUNDERID_BASE_URL}
        scopes={import.meta.env.VITE_THUNDERID_SCOPES}
        afterSignInUrl={import.meta.env.VITE_THUNDERID_AFTER_SIGN_IN_URL || origin}
        afterSignOutUrl={import.meta.env.VITE_THUNDERID_AFTER_SIGN_OUT_URL || origin}
      >
        <QueryClientProvider client={queryClient}>
          <BrowserRouter>
            <App />
          </BrowserRouter>
        </QueryClientProvider>
      </ThunderIDProvider>
    </ErrorBoundary>
  </StrictMode>,
);
