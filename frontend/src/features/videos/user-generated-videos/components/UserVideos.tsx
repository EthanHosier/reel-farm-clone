import { useState } from "react";
import type { UserGeneratedVideo } from "@/api";
import { EmptyState } from "./EmptyState";
import { Dialog, DialogClose, DialogContent } from "@/components/ui/dialog";
import { XIcon } from "lucide-react";

interface UserVideosProps {
  userVideos: UserGeneratedVideo[];
}

export function UserVideos({ userVideos }: UserVideosProps) {
  const [selectedVideo, setSelectedVideo] = useState<UserGeneratedVideo | null>(
    null
  );

  if (userVideos.length === 0) {
    return <EmptyState />;
  }

  return (
    <>
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
        {userVideos.map((video: UserGeneratedVideo) => (
          <div
            key={video.id}
            className="bg-white rounded-lg border shadow-sm hover:shadow-md transition-shadow cursor-pointer"
            onClick={() => setSelectedVideo(video)}
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

      <Dialog
        open={!!selectedVideo}
        onOpenChange={(open) => !open && setSelectedVideo(null)}
      >
        <DialogContent
          showCloseButton={false}
          className="max-w-2xl p-0 overflow-hidden border-none"
        >
          {selectedVideo && (
            <div className="w-full h-full bg-black flex items-center justify-center">
              <video
                src={selectedVideo.video_url}
                controls
                autoPlay
                className="w-full h-full object-contain"
              />
              <DialogClose className="absolute top-4 right-4 bg-white/80 rounded-full p-2 cursor-pointer">
                <XIcon className="size-4" />
              </DialogClose>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </>
  );
}
