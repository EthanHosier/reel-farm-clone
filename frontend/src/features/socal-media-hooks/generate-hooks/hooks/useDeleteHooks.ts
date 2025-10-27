import { HooksService } from "@/api";
import { hooksKeys } from "@/hooks/useGetHooks";
import { useMutation, useQueryClient } from "@tanstack/react-query";

// Delete hooks bulk mutation
export function useDeleteHooksBulk() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (request: { hook_ids: string[] }) =>
      HooksService.deleteHooksBulk(request),
    onSuccess: () => {
      // Invalidate and refetch hooks list
      queryClient.invalidateQueries({ queryKey: hooksKeys.lists() });
    },
  });
}
