/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Hook } from './Hook';
export type GetHooksResponse = {
    /**
     * Array of user's hooks
     */
    hooks: Array<Hook>;
    /**
     * Total number of hooks for the user
     */
    total_count: number;
    /**
     * Video title
     */
    title?: string;
    /**
     * CloudFront URL for video download
     */
    video_url?: string;
    /**
     * CloudFront URL for thumbnail
     */
    thumbnail_url?: string;
    /**
     * When the video was last updated
     */
    updated_at?: string;
};

