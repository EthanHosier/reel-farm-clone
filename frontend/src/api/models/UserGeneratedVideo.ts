/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type UserGeneratedVideo = {
    /**
     * Unique identifier for the user-generated video
     */
    id: string;
    /**
     * ID of the user who generated the video
     */
    user_id: string;
    /**
     * ID of the original AI avatar video
     */
    ai_avatar_video_id: string;
    /**
     * Text that was overlaid on the video
     */
    overlay_text: string;
    /**
     * CloudFront URL for the generated video
     */
    video_url: string;
    /**
     * CloudFront URL for the video thumbnail
     */
    thumbnail_url: string;
    /**
     * Current processing status
     */
    status: UserGeneratedVideo.status;
    /**
     * When the video was created
     */
    created_at: string;
};
export namespace UserGeneratedVideo {
    /**
     * Current processing status
     */
    export enum status {
        PROCESSING = 'processing',
        COMPLETED = 'completed',
        FAILED = 'failed',
    }
}

