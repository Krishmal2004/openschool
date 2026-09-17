import { useThunderID } from "@thunderid/react";
import { Navigate, useLocation } from "react-router";
import { useIdleLogout } from "@/shared/auth/useIdleLogout";

export default function ProtectedRoute({
  children,
}: {
  children: React.ReactNode;
}) {
  const { isSignedIn, isLoading } = useThunderID();
  const location = useLocation();
  useIdleLogout();

  if (isLoading) return <div className="os-full-height" />;

  if (!isSignedIn) {
    return <Navigate to="/signin" replace state={{ from: location }} />;
  }

  return <>{children}</>;
}
