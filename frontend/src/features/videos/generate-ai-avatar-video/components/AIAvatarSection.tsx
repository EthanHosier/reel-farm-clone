import { useState } from "react";
import { Button } from "@/components/ui/button";
import { useAIAvatarVideos } from "@/features/videos/generate-ai-avatar-video/queries/useAIAvatarVideos";
import { ChevronLeft, ChevronRight } from "lucide-react";

interface AIAvatarSectionProps {
  selectedAvatarVideoId: string | null;
  onVideoSelect: (videoId: string, videoUrl: string) => void;
}

export function AIAvatarSection({
  selectedAvatarVideoId,
  onVideoSelect,
}: AIAvatarSectionProps) {
  const [page, setPage] = useState(0);
  const itemsPerPage = 32;

  const {
    data: aiAvatarVideos,
    isLoading: videosLoading,
    error: videosError,
  } = useAIAvatarVideos();

  const totalVideos = aiAvatarVideos?.videos.length || 0;
  const totalPages = Math.ceil(totalVideos / itemsPerPage);
  const startIndex = page * itemsPerPage;
  const endIndex = startIndex + itemsPerPage;
  const currentPageVideos =
    aiAvatarVideos?.videos.slice(startIndex, endIndex) || [];

  const handlePreviousPage = () => {
    if (page > 0) setPage(page - 1);
  };

  const handleNextPage = () => {
    if (page < totalPages - 1) setPage(page + 1);
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-sm font-medium">2. AI avatar</h3>
        {totalPages > 0 && (
          <span className="text-xs text-gray-500">
            {page + 1}/{totalPages}
          </span>
        )}
      </div>

      {videosLoading && (
        <p className="text-sm text-gray-500 py-4">
          Loading AI avatar videos...
        </p>
      )}

      {videosError && (
        <p className="text-sm text-red-600 py-4">
          Error: {videosError.message}
        </p>
      )}

      {aiAvatarVideos && (
        <>
          <div className="grid grid-cols-8 gap-2 mb-2">
            {currentPageVideos.map((video) => (
              <button
                key={video.id}
                onClick={() => onVideoSelect(video.id, video.video_url)}
                className={`aspect-square rounded-lg overflow-hidden border-2 transition ${
                  selectedAvatarVideoId === video.id
                    ? "border-blue-500 ring-2 ring-blue-200"
                    : "border-gray-200 hover:border-gray-300"
                }`}
              >
                <img
                  src={video.thumbnail_url}
                  alt={video.title}
                  className="w-full h-full object-cover"
                />
              </button>
            ))}
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-between text-xs text-gray-500">
              <Button
                variant="ghost"
                size="icon"
                onClick={handlePreviousPage}
                disabled={page === 0}
                className="h-6 w-6"
              >
                <ChevronLeft className="h-3 w-3" />
              </Button>
              <span>
                {page + 1}/{totalPages}
              </span>
              <Button
                variant="ghost"
                size="icon"
                onClick={handleNextPage}
                disabled={page === totalPages - 1}
                className="h-6 w-6"
              >
                <ChevronRight className="h-3 w-3" />
              </Button>
            </div>
          )}
        </>
      )}
    </div>
  );
}
