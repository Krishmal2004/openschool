import { useThunderID } from "@thunderid/react";
import { useQuery } from "@tanstack/react-query";
import api from "@/shared/api/client";
import { keys } from "@/shared/api/keys";

export interface Me {
  user_id: string;
  email: string;
  username: string;
  given_name: string;
  family_name: string;
  phone_number: string;
  roles: string[];
  must_change_password: boolean;
}

// GET /me provisions the local user row on first sign-in and reports the must-change-password flag.
export function useProvisionUser() {
  const { isSignedIn, isLoading } = useThunderID();
  return useQuery({
    queryKey: keys.me.profile(),
    queryFn: () => api.get<Me>("/me").then((r) => r.data),
    enabled: !isLoading && isSignedIn,
  });
}
