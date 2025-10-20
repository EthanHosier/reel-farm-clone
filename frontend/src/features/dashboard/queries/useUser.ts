import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { ROUTES } from "@/lib/routes";
import type { UserAccount } from "@/api";

export function useUser() {
  return useQuery({
    queryKey: ["user"],
    queryFn: async (): Promise<UserAccount> => {
      const response = await api.get(ROUTES.USER);
      return response.data;
    },
    staleTime: 2 * 60 * 1000, // 2 minutes
    retry: (failureCount, error: any) => {
      // Don't retry on 401/403 errors
      if (error?.response?.status === 401 || error?.response?.status === 403) {
        return false;
      }
      return failureCount < 3;
    },
  });
}
