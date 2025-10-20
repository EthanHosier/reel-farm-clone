import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { CACHE_KEYS } from "@/lib/cacheKeys";
import type { UserAccount } from "@/lib/api";

export function useUser() {
  return useQuery({
    queryKey: CACHE_KEYS.USER,
    queryFn: async (): Promise<UserAccount> => {
      return await api.users.getUserAccount();
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
