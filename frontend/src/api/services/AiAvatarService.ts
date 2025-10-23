/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { AIAvatarVideosResponse } from '../models/AIAvatarVideosResponse';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class AiAvatarService {
    /**
     * Get all AI avatar videos
     * Retrieves all AI avatar videos with thumbnails for gallery view
     * @returns AIAvatarVideosResponse Videos retrieved successfully
     * @throws ApiError
     */
    public static getAiAvatarVideos(): CancelablePromise<AIAvatarVideosResponse> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/ai-avatar/videos',
            errors: {
                401: `Unauthorized - invalid or missing token`,
                500: `Internal server error`,
            },
        });
    }
}
