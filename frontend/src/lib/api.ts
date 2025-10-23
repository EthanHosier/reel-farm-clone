import { HooksService } from "@/api";
import { OpenAPI } from "@/api/core/OpenAPI";
import { HealthService } from "@/api/services/HealthService";
import { SubscriptionsService } from "@/api/services/SubscriptionsService";
import { UsersService } from "@/api/services/UsersService";

// Get API URL from environment variable
const API_URL = import.meta.env.VITE_API_URL || "http://localhost:3000";

console.log("API URL:", API_URL);

// Configure OpenAPI with auth token injection
const configureOpenAPI = () => {
  OpenAPI.BASE = API_URL;
  OpenAPI.TOKEN = async () => {
    // Get the Supabase auth token from localStorage
    const token = JSON.parse(
      localStorage.getItem("sb-uokbfbxpadivnvjnlhhx-auth-token") || "{}"
    )?.access_token;
    return token || "";
  };
};

// Initialize the configuration
configureOpenAPI();

// Export configured services
export const api = {
  health: HealthService,
  users: UsersService,
  subscriptions: SubscriptionsService,
  hooks: HooksService,
};

// Re-export types for convenience
export type {
  HealthResponse,
  UserAccount,
  CreateCheckoutSessionRequest,
  CreateCustomerPortalRequest,
  CheckoutSessionResponse,
  CustomerPortalResponse,
  Hook,
  GenerateHooksRequest,
  GenerateHooksResponse,
  GetHooksResponse,
} from "@/api";
