import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { CACHE_KEYS } from "@/lib/cacheKeys";
import type { AIAvatarVideosResponse } from "@/api/models/AIAvatarVideosResponse";

export function useAIAvatarVideos() {
  return useQuery({
    queryKey: CACHE_KEYS.AI_AVATAR_VIDEOS,
    queryFn: async (): Promise<AIAvatarVideosResponse> => {
      return await api.aiAvatar.getAiAvatarVideos();
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}
