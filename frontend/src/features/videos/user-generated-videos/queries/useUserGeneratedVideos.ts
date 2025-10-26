import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { CACHE_KEYS } from "@/lib/cacheKeys";
import type { UserGeneratedVideosResponse } from "@/api/models/UserGeneratedVideosResponse";

export function useUserGeneratedVideos() {
  return useQuery({
    queryKey: CACHE_KEYS.USER_GENERATED_VIDEOS,
    queryFn: async (): Promise<UserGeneratedVideosResponse> => {
      return await api.userGeneratedVideos.getUserGeneratedVideos();
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

