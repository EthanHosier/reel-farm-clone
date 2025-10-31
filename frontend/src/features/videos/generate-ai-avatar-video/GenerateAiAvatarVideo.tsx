import { useState, useEffect } from "react";
import { HookSection } from "@/features/videos/generate-ai-avatar-video/components/HookSection";
import { AIAvatarSection } from "@/features/videos/generate-ai-avatar-video/components/AIAvatarSection";
import { VideoPreview } from "@/features/videos/generate-ai-avatar-video/components/VideoPreview";
import { useAIAvatarVideos } from "@/features/videos/generate-ai-avatar-video/queries/useAIAvatarVideos";
import { Button } from "@/components/ui/button";
import { Title } from "@/components/dashboard/Title";
import { useCreateUserGeneratedVideo } from "@/features/videos/generate-ai-avatar-video/queries/useCreateUserGeneratedVideo";
import { toast } from "sonner";
import { Loader2 } from "lucide-react";

export const GenerateAiAvatarVideo = () => {
  const { data: aiAvatarVideos } = useAIAvatarVideos();
  const { mutateAsync: createVideo, isPending: isCreatingVideo } =
    useCreateUserGeneratedVideo();

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

  const handleGenerateVideo = async () => {
    if (!selectedAvatarVideoId) {
      toast.error("Please select an AI avatar video");
      return;
    }

    if (!overlayText) {
      toast.error("Please enter text to overlay on the video");
      return;
    }

    try {
      await createVideo({
        ai_avatar_video_id: selectedAvatarVideoId,
        overlay_text: overlayText,
      });
    } catch (error) {
      toast.error("Failed to generate video");
      return;
    }

    toast.success("Video generated successfully");
  };

  return (
    <div>
      <Title
        title="Generate AI Avatar Video"
        description="Generate videos with your own text overlay."
      />
      <div className="flex flex-col lg:grid lg:grid-cols-2 gap-4 lg:h-[calc(100vh-8rem)]">
        <div className="lg:order-1 order-2 overflow-y-auto space-y-6">
          <HookSection onTextChange={handleTextChange} />
          <AIAvatarSection
            selectedAvatarVideoId={selectedAvatarVideoId}
            onVideoSelect={handleVideoSelect}
          />
        </div>

        <div className="lg:order-2 order-1 flex flex-col items-center justify-center rounded-lg p-6">
          <VideoPreview
            selectedVideo={selectedVideo}
            overlayText={overlayText}
          />
          <Button
            size="lg"
            className="mt-6 w-full max-w-xs"
            disabled={!selectedAvatarVideoId || !overlayText || isCreatingVideo}
            onClick={handleGenerateVideo}
          >
            {isCreatingVideo ? (
              <Loader2 className="size-4 animate-spin" />
            ) : (
              "Generate UGC"
            )}
          </Button>
        </div>
      </div>
    </div>
  );
};
