import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useUserGeneratedVideos } from "../queries/useUserGeneratedVideos";

interface UserGeneratedVideosProps {
  onVideoSelect: (videoUrl: string) => void;
}

export function UserGeneratedVideos({
  onVideoSelect,
}: UserGeneratedVideosProps) {
  const {
    data: userGeneratedVideos,
    isLoading: userVideosLoading,
    error: userVideosError,
  } = useUserGeneratedVideos();

  return (
    <Card>
      <CardHeader>
        <CardTitle>My Generated Videos</CardTitle>
        <CardDescription>
          Videos you've created with custom text overlays
        </CardDescription>
      </CardHeader>
      <CardContent>
        {userVideosLoading && (
          <p className="text-blue-600">Loading your videos...</p>
        )}
        {userVideosError && (
          <p className="text-red-600">Error: {userVideosError.message}</p>
        )}
        {userGeneratedVideos && userGeneratedVideos.videos.length > 0 && (
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {userGeneratedVideos.videos.map((video) => (
              <div
                key={video.id}
                className="bg-white rounded-lg border shadow-sm hover:shadow-md transition-shadow cursor-pointer"
                onClick={() => onVideoSelect(video.video_url)}
              >
                <div className="aspect-[9/16] bg-gray-100 rounded-lg overflow-hidden">
                  <img
                    src={video.thumbnail_url}
                    alt={video.overlay_text}
                    className="w-full h-full object-cover hover:scale-105 transition-transform duration-200"
                  />
                </div>
              </div>
            ))}
          </div>
        )}
        {userGeneratedVideos && userGeneratedVideos.videos.length === 0 && (
          <div className="text-center py-8 text-gray-500">
            <p>No generated videos yet</p>
            <p className="text-sm mt-1">
              Select an AI avatar video and generate your first custom video!
            </p>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
