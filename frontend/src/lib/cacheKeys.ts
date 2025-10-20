// Global cache keys for TanStack Query
export const CACHE_KEYS = {
  // Health endpoints
  HEALTH: ["health"] as const,

  // User endpoints
  USER: ["user"] as const,
} as const;

// Type for cache key values
export type CacheKey = (typeof CACHE_KEYS)[keyof typeof CACHE_KEYS];
