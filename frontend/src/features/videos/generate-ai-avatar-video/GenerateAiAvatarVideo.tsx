import { useState, useEffect } from "react";
import { HookSection } from "@/features/videos/generate-ai-avatar-video/components/HookSection";
import { AIAvatarSection } from "@/features/videos/generate-ai-avatar-video/components/AIAvatarSection";
import { VideoPreview } from "@/features/videos/generate-ai-avatar-video/components/VideoPreview";
import { useAIAvatarVideos } from "@/features/videos/generate-ai-avatar-video/queries/useAIAvatarVideos";
import { Button } from "@/components/ui/button";
import { Title } from "@/components/dashboard/Title";

export const GenerateAiAvatarVideo = () => {
  const { data: aiAvatarVideos } = useAIAvatarVideos();

  // Video preview state
  const [selectedVideo, setSelectedVideo] = useState<string | null>(null);

  // Video generation state
  const [selectedAvatarVideoId, setSelectedAvatarVideoId] = useState<
    string | null
  >(null);
  const [overlayText, setOverlayText] = useState("");

  // Set default selected video when videos load
  useEffect(() => {
    if (aiAvatarVideos && aiAvatarVideos.videos.length > 0 && !selectedVideo) {
      setSelectedVideo(aiAvatarVideos.videos[0].video_url);
      setSelectedAvatarVideoId(aiAvatarVideos.videos[0].id);
    }
  }, [aiAvatarVideos, selectedVideo]);

  const handleVideoSelect = (videoId: string, videoUrl: string) => {
    setSelectedAvatarVideoId(videoId);
    setSelectedVideo(videoUrl);
  };

  const handleTextChange = (newText: string) => {
    setOverlayText(newText);
  };

  return (
    <div>
      <Title
        title="Generate AI Avatar Video"
        description="Generate videos with your own text overlay."
      />
      <div className="grid grid-cols-2 gap-4 h-[calc(100vh-8rem)]">
        <div className="overflow-y-auto space-y-6">
          <HookSection onTextChange={handleTextChange} />
          <AIAvatarSection
            selectedAvatarVideoId={selectedAvatarVideoId}
            onVideoSelect={handleVideoSelect}
          />
        </div>

        <div className="flex flex-col items-center justify-center bg-gray-50 rounded-lg p-6">
          <VideoPreview
            selectedVideo={selectedVideo}
            overlayText={overlayText}
          />
          <Button
            size="lg"
            className="mt-6 w-full max-w-xs"
            disabled={!selectedAvatarVideoId || !overlayText}
          >
            Subscription required to use
          </Button>
        </div>
      </div>
    </div>
  );
};
