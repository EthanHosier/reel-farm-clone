import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { HooksService, type GenerateHooksRequest } from "@/api";

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

// Generate hooks mutation
export function useGenerateHooks() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (request: GenerateHooksRequest) =>
      HooksService.generateHooks(request),
    onSuccess: () => {
      // Invalidate and refetch hooks list
      queryClient.invalidateQueries({ queryKey: hooksKeys.lists() });
    },
  });
}

// Delete hook mutation
export function useDeleteHook() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (hookId: string) => HooksService.deleteHook(hookId),
    onSuccess: () => {
      // Invalidate and refetch hooks list
      queryClient.invalidateQueries({ queryKey: hooksKeys.lists() });
    },
  });
}
