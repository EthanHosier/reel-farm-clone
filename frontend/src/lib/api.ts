import { OpenAPI } from "@/api/core/OpenAPI";
import { HealthService } from "@/api/services/HealthService";
import { UsersService } from "@/api/services/UsersService";

// Configure OpenAPI with auth token injection
const configureOpenAPI = () => {
  OpenAPI.TOKEN = async () => {
    // Get the Supabase auth token from localStorage
    const token = localStorage.getItem("sb-uokbfbxpadivnvjnlhhx-auth-token");
    return token || "";
  };
};

// Initialize the configuration
configureOpenAPI();

// Export configured services
export const api = {
  health: HealthService,
  users: UsersService,
};

// Re-export types for convenience
export type { HealthResponse, UserAccount } from "@/api";
