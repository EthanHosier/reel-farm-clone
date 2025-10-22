/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type CreateCheckoutSessionRequest = {
    /**
     * Stripe price ID for the subscription
     */
    price_id: string;
    /**
     * URL to redirect to after successful payment
     */
    success_url: string;
    /**
     * URL to redirect to if payment is canceled
     */
    cancel_url: string;
};

