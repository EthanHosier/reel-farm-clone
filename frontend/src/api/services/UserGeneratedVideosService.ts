/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { CreateUserGeneratedVideoRequest } from '../models/CreateUserGeneratedVideoRequest';
import type { UserGeneratedVideoResponse } from '../models/UserGeneratedVideoResponse';
import type { UserGeneratedVideosResponse } from '../models/UserGeneratedVideosResponse';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class UserGeneratedVideosService {
    /**
     * Get user-generated videos
     * Retrieves all user-generated videos for the authenticated user
     * @returns UserGeneratedVideosResponse User-generated videos retrieved successfully
     * @throws ApiError
     */
    public static getUserGeneratedVideos(): CancelablePromise<UserGeneratedVideosResponse> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/user-generated-videos',
            errors: {
                401: `Unauthorized - invalid or missing token`,
                500: `Internal server error`,
            },
        });
    }
    /**
     * Generate a video with text overlay
     * Creates a new user-generated video by adding text overlay to an existing AI avatar video
     * @param requestBody
     * @returns UserGeneratedVideoResponse Video generated successfully
     * @throws ApiError
     */
    public static createUserGeneratedVideo(
        requestBody: CreateUserGeneratedVideoRequest,
    ): CancelablePromise<UserGeneratedVideoResponse> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/user-generated-videos',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `Bad request - invalid input`,
                401: `Unauthorized - invalid or missing token`,
                404: `AI avatar video not found`,
                500: `Internal server error`,
            },
        });
    }
}
