import { useAuth } from "@/contexts/AuthContext";
import { Button } from "@/components/ui/button";
import { useHealth } from "./queries/useHealth";
import { useUser } from "./queries/useUser";
import { HooksManager } from "./components/HooksManager";
import { useAIAvatarVideos } from "./queries/useAIAvatarVideos";
import { AIAvatarVideos } from "./components/AIAvatarVideos";
import { VideoPreview } from "./components/VideoPreview";
import { VideoGenerationForm } from "./components/VideoGenerationForm";
import { UserGeneratedVideos } from "./components/UserGeneratedVideos";
import { AccountInfo } from "./components/AccountInfo";
import { HealthStatus } from "./components/HealthStatus";
import React, { useState } from "react";

export default function Dashboard() {
  const { user, signOut } = useAuth();
  const {
    data: health,
    isLoading: healthLoading,
    error: healthError,
  } = useHealth();
  const {
    data: userAccount,
    isLoading: userLoading,
    error: userError,
  } = useUser();
  const { data: aiAvatarVideos } = useAIAvatarVideos();

  // Video preview state
  const [selectedVideo, setSelectedVideo] = useState<string | null>(null);

  // Video generation state
  const [selectedAvatarVideoId, setSelectedAvatarVideoId] = useState<
    string | null
  >(null);
  const [overlayText, setOverlayText] = useState("");

  // Set default selected video when videos load
  React.useEffect(() => {
    if (aiAvatarVideos && aiAvatarVideos.videos.length > 0 && !selectedVideo) {
      setSelectedVideo(aiAvatarVideos.videos[0].video_url);
    }
  }, [aiAvatarVideos, selectedVideo]);

  const handleSignOut = async () => {
    try {
      await signOut();
    } catch (error) {
      console.error("Error signing out:", error);
    }
  };

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

  const handleUserVideoSelect = (videoUrl: string) => {
    setSelectedVideo(videoUrl);
  };

  if (healthLoading || userLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading dashboard...</p>
        </div>
      </div>
    );
  }

  if (healthError || userError) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-red-600 mb-4">Error</h1>
          <p className="text-gray-600 mb-4">
            {healthError?.message || userError?.message}
          </p>
          <Button onClick={() => window.location.reload()}>Retry</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
            <p className="text-gray-600">Welcome back, {user?.email}!</p>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={handleSignOut}>
              Sign Out
            </Button>
          </div>
        </div>

        <HealthStatus health={health} />

        {userAccount && (
          <div className="mb-6">
            <AccountInfo userAccount={userAccount} />
          </div>
        )}

        <div className="mb-6">
          <HooksManager />
        </div>

        <div className="mb-6">
          <AIAvatarVideos
            selectedAvatarVideoId={selectedAvatarVideoId}
            onVideoSelect={handleVideoSelect}
          />
        </div>

        <div className="mb-6">
          <VideoGenerationForm
            selectedAvatarVideoId={selectedAvatarVideoId}
            onCancel={handleCancelGeneration}
            onTextChange={handleTextChange}
          />
        </div>

        <div className="mb-6">
          <VideoPreview
            selectedVideo={selectedVideo}
            overlayText={overlayText}
          />
        </div>

        <div className="mb-6">
          <UserGeneratedVideos onVideoSelect={handleUserVideoSelect} />
        </div>
      </div>
    </div>
  );
}
