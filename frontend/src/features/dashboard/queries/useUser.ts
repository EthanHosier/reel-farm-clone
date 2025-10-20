import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import type { UserAccount } from "@/api";

export function useUser() {
  return useQuery({
    queryKey: ["user"],
    queryFn: async (): Promise<UserAccount> => {
      const response = await apiClient.get("/user");
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
