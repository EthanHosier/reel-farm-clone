import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { CACHE_KEYS } from "@/lib/cacheKeys";
import type { HealthResponse } from "@/lib/api";

export function useHealth() {
  return useQuery({
    queryKey: CACHE_KEYS.HEALTH,
    queryFn: async (): Promise<HealthResponse> => {
      return await api.health.getHealth();
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}
