import { type GenerateHooksRequest, HooksService } from "@/api";
import { hooksKeys } from "@/hooks/useGetHooks";
import { useMutation, useQueryClient } from "@tanstack/react-query";

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
