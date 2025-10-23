/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { CheckoutSessionResponse } from '../models/CheckoutSessionResponse';
import type { CreateCheckoutSessionRequest } from '../models/CreateCheckoutSessionRequest';
import type { CreateCustomerPortalRequest } from '../models/CreateCustomerPortalRequest';
import type { CustomerPortalResponse } from '../models/CustomerPortalResponse';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class SubscriptionsService {
    /**
     * Create Stripe checkout session
     * Creates a Stripe checkout session for subscription upgrade
     * @param requestBody
     * @returns CheckoutSessionResponse Checkout session created successfully
     * @throws ApiError
     */
    public static createCheckoutSession(
        requestBody: CreateCheckoutSessionRequest,
    ): CancelablePromise<CheckoutSessionResponse> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/subscription/create-checkout-session',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `Bad request - invalid request data`,
                401: `Unauthorized - invalid or missing token`,
                500: `Internal server error`,
            },
        });
    }
    /**
     * Create customer portal session
     * Creates a Stripe customer portal session for subscription management
     * @param requestBody
     * @returns CustomerPortalResponse Customer portal session created successfully
     * @throws ApiError
     */
    public static createCustomerPortalSession(
        requestBody: CreateCustomerPortalRequest,
    ): CancelablePromise<CustomerPortalResponse> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/subscription/customer-portal',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `Bad request - user not found or no Stripe customer`,
                401: `Unauthorized - invalid or missing token`,
                500: `Internal server error`,
            },
        });
    }
}
