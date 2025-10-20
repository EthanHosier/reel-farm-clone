import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import type { HealthResponse } from "@/api";

export function useHealth() {
  return useQuery({
    queryKey: ["health"],
    queryFn: async (): Promise<HealthResponse> => {
      const response = await apiClient.get("/health");
      return response.data;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}
