import { useMutation } from "@tanstack/react-query";
import { api } from "@/lib/api";
import type { CreateCheckoutSessionRequest } from "@/api/models/CreateCheckoutSessionRequest";
import type { CreateCustomerPortalRequest } from "@/api/models/CreateCustomerPortalRequest";

interface UseSubscriptionMutationsOptions {
  onCheckoutSuccess?: (checkoutUrl: string) => void;
  onCheckoutError?: (error: Error) => void;
  onPortalSuccess?: (portalUrl: string) => void;
  onPortalError?: (error: Error) => void;
}

export function useSubscriptionMutations(
  options?: UseSubscriptionMutationsOptions
) {
  const createCheckoutMutation = useMutation({
    mutationFn: async (request: CreateCheckoutSessionRequest) => {
      return await api.subscriptions.createCheckoutSession(request);
    },
    onSuccess: (response) => {
      options?.onCheckoutSuccess?.(response.checkout_url);
    },
    onError: (error: Error) => {
      console.error("Error creating checkout session:", error);
      options?.onCheckoutError?.(error);
    },
  });

  const createPortalMutation = useMutation({
    mutationFn: async (request: CreateCustomerPortalRequest) => {
      return await api.subscriptions.createCustomerPortalSession(request);
    },
    onSuccess: (response) => {
      options?.onPortalSuccess?.(response.portal_url);
    },
    onError: (error: Error) => {
      console.error("Error creating customer portal session:", error);
      options?.onPortalError?.(error);
    },
  });

  return {
    createCheckout: createCheckoutMutation,
    createPortal: createPortalMutation,
  };
}
