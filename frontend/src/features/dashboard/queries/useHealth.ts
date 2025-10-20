import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { ROUTES } from "@/lib/routes";
import type { HealthResponse } from "@/api";

export function useHealth() {
  return useQuery({
    queryKey: ["health"],
    queryFn: async (): Promise<HealthResponse> => {
      const response = await api.get(ROUTES.HEALTH);
      return response.data;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}
