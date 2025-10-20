import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { HealthResponse } from "@/lib/api";

export function useHealth() {
  return useQuery({
    queryKey: ["health"],
    queryFn: async (): Promise<HealthResponse> => {
      return await api.health.getHealth();
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}
