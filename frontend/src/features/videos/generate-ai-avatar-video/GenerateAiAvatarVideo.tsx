import { useState, useEffect } from "react";
import { AIAvatarVideos } from "@/features/videos/generate-ai-avatar-video/components/AIAvatarVideos";
import { VideoGenerationForm } from "@/features/videos/generate-ai-avatar-video/components/VideoGenerationForm";
import { VideoPreview } from "@/features/videos/generate-ai-avatar-video/components/VideoPreview";
import { useAIAvatarVideos } from "@/features/videos/generate-ai-avatar-video/queries/useAIAvatarVideos";

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
    }
  }, [aiAvatarVideos, selectedVideo]);

  const handleVideoSelect = (videoId: string, videoUrl: string) => {
    setSelectedAvatarVideoId(videoId);
    setSelectedVideo(videoUrl);
  };

  const handleCancelGeneration = () => {
    setSelectedAvatarVideoId(null);
    setOverlayText("");
  };

  const handleTextChange = (newText: string) => {
    setOverlayText(newText);
  };

  return (
    <div className="space-y-6">
      <div>
        <AIAvatarVideos
          selectedAvatarVideoId={selectedAvatarVideoId}
          onVideoSelect={handleVideoSelect}
        />
      </div>

      <div>
        <VideoGenerationForm
          selectedAvatarVideoId={selectedAvatarVideoId}
          onCancel={handleCancelGeneration}
          onTextChange={handleTextChange}
        />
      </div>

      <div>
        <VideoPreview selectedVideo={selectedVideo} overlayText={overlayText} />
      </div>
    </div>
  );
};
