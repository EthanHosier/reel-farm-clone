import { useQuery } from "@tanstack/react-query";
import { HooksService } from "@/api";

// Query key factory
export const hooksKeys = {
  all: ["hooks"] as const,
  lists: () => [...hooksKeys.all, "list"] as const,
  list: (limit: number, offset: number) =>
    [...hooksKeys.lists(), { limit, offset }] as const,
};

// Get hooks query
export function useHooks(limit: number = 20, offset: number = 0) {
  return useQuery({
    queryKey: hooksKeys.list(limit, offset),
    queryFn: () => HooksService.getHooks(limit, offset),
  });
}
