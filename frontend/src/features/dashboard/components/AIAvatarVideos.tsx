import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useAIAvatarVideos } from "../queries/useAIAvatarVideos";

interface AIAvatarVideosProps {
  selectedAvatarVideoId: string | null;
  onVideoSelect: (videoId: string, videoUrl: string) => void;
}

export function AIAvatarVideos({
  selectedAvatarVideoId,
  onVideoSelect,
}: AIAvatarVideosProps) {
  const {
    data: aiAvatarVideos,
    isLoading: videosLoading,
    error: videosError,
  } = useAIAvatarVideos();

  return (
    <Card>
      <CardHeader>
        <CardTitle>AI Avatar Videos</CardTitle>
        <CardDescription>
          Click on a thumbnail to select it for video generation with your own
          text overlay
        </CardDescription>
      </CardHeader>
      <CardContent>
        {videosLoading && (
          <p className="text-blue-600">Loading AI avatar videos...</p>
        )}
        {videosError && (
          <p className="text-red-600">Error: {videosError.message}</p>
        )}
        {aiAvatarVideos && (
          <div className="grid grid-cols-4 md:grid-cols-6 lg:grid-cols-8 gap-4">
            {aiAvatarVideos.videos.map((video) => (
              <div
                key={video.id}
                className={`bg-white rounded-lg border shadow-sm hover:shadow-md transition-shadow cursor-pointer ${
                  selectedAvatarVideoId === video.id
                    ? "ring-2 ring-blue-500"
                    : ""
                }`}
                onClick={() => onVideoSelect(video.id, video.video_url)}
              >
                <div className="aspect-square bg-gray-100 rounded-lg overflow-hidden">
                  <img
                    src={video.thumbnail_url}
                    alt={video.title}
                    className="w-full h-full object-cover hover:scale-105 transition-transform duration-200"
                  />
                </div>
              </div>
            ))}
          </div>
        )}
        {aiAvatarVideos && aiAvatarVideos.videos.length === 0 && (
          <div className="text-center py-8 text-gray-500">
            <p>No AI avatar videos available</p>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
