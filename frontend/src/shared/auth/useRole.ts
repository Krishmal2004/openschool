import { useState, useEffect, useRef } from "react";
import { useThunderID } from "@thunderid/react";
import { parseJwt, resolveRole, type Role } from "@/shared/auth/jwt";

export type { Role };

export function useRole(): { role: Role | null; loading: boolean } {
  const { getAccessToken, isSignedIn, isLoading } = useThunderID();
  const getAccessTokenRef = useRef(getAccessToken);
  const [role, setRole] = useState<Role | null>(null);
  const [roleResolved, setRoleResolved] = useState(false);

  useEffect(() => {
    getAccessTokenRef.current = getAccessToken;
  });

  useEffect(() => {
    if (isLoading || !isSignedIn) return;
    let cancelled = false;

    getAccessTokenRef.current().then((token) => {
      if (cancelled) return;
      if (token) setRole(resolveRole(parseJwt(token)));
      setRoleResolved(true);
    });

    return () => {
      cancelled = true;
      // Clear on sign-out so a stale role never shows for the next session.
      setRole(null);
      setRoleResolved(false);
    };
  }, [isLoading, isSignedIn]);

  return { role, loading: isLoading || (isSignedIn && !roleResolved) };
}
