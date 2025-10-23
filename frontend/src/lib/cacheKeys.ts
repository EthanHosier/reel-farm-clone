// Global cache keys for TanStack Query
export const CACHE_KEYS = {
  // Health endpoints
  HEALTH: ["health"] as const,

  // User endpoints
  USER: ["user"] as const,

  // AI Avatar endpoints
  AI_AVATAR_VIDEOS: ["ai-avatar-videos"] as const,
} as const;

// Type for cache key values
export type CacheKey = (typeof CACHE_KEYS)[keyof typeof CACHE_KEYS];
