/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { UserAccount } from '../models/UserAccount';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class UsersService {
    /**
     * Get current user account
     * Retrieves the authenticated user's account information
     * @returns UserAccount User account retrieved successfully
     * @throws ApiError
     */
    public static getUserAccount(): CancelablePromise<UserAccount> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/user',
            errors: {
                401: `Unauthorized - invalid or missing token`,
                404: `User account not found`,
                500: `Internal server error`,
            },
        });
    }
}
