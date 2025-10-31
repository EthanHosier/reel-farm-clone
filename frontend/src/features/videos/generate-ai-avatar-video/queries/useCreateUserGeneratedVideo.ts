import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { CACHE_KEYS } from "@/lib/cacheKeys";
import type { CreateUserGeneratedVideoRequest } from "@/api/models/CreateUserGeneratedVideoRequest";

interface UseCreateUserGeneratedVideoOptions {
  onSuccess?: () => void;
  onError?: (error: Error) => void;
}

export function useCreateUserGeneratedVideo(
  options?: UseCreateUserGeneratedVideoOptions
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (request: CreateUserGeneratedVideoRequest) => {
      return await api.userGeneratedVideos.createUserGeneratedVideo(request);
    },
    onSuccess: () => {
      // Invalidate and refetch user-generated videos
      queryClient.invalidateQueries({
        queryKey: CACHE_KEYS.USER_GENERATED_VIDEOS,
      });

      // Call custom success handler if provided
      options?.onSuccess?.();
    },
    onError: (error: Error) => {
      console.error("Error generating video:", error);

      // Call custom error handler if provided
      options?.onError?.(error);
    },
  });
}
