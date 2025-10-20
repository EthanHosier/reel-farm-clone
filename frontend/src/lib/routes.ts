// Global API routes constants
export const ROUTES = {
  // Health endpoints
  HEALTH: "/health",

  // User endpoints
  USER: "/user",
} as const;

// Type for route values
export type Route = (typeof ROUTES)[keyof typeof ROUTES];
