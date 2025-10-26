import type { UserGeneratedVideo } from "@/api";
import { EmptyState } from "./EmptyState";

interface UserVideosProps {
  userVideos: UserGeneratedVideo[];
}

export function UserVideos({ userVideos }: UserVideosProps) {
  if (userVideos.length === 0 || true) {
    return <EmptyState />;
  }

  return (
    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
      {userVideos.map((video: UserGeneratedVideo) => (
        <div
          key={video.id}
          className="bg-white rounded-lg border shadow-sm hover:shadow-md transition-shadow cursor-pointer"
          onClick={() => {}}
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
  );
}
