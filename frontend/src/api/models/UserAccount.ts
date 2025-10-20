/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type UserAccount = {
    /**
     * Unique identifier for the user account
     */
    id: string;
    /**
     * User's subscription plan
     */
    plan: string;
    /**
     * When the current plan started
     */
    plan_started_at: string;
    /**
     * When the current plan ends (null for free plans)
     */
    plan_ends_at?: string | null;
    /**
     * External billing system customer ID
     */
    billing_customer_id?: string | null;
    /**
     * When the account was created
     */
    created_at: string;
    /**
     * When the account was last updated
     */
    updated_at: string;
};

