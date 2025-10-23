import { useAuth } from "@/contexts/AuthContext";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useHealth } from "./queries/useHealth";
import { useUser } from "./queries/useUser";
import { HooksManager } from "./components/HooksManager";
import { useAIAvatarVideos } from "./queries/useAIAvatarVideos";
import { useUserGeneratedVideos } from "./queries/useUserGeneratedVideos";
import { api } from "@/lib/api";
import React, { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { CACHE_KEYS } from "@/lib/cacheKeys";

export default function Dashboard() {
  const { user, session, signOut } = useAuth();
  const queryClient = useQueryClient();
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
  const {
    data: aiAvatarVideos,
    isLoading: videosLoading,
    error: videosError,
  } = useAIAvatarVideos();
  const {
    data: userGeneratedVideos,
    isLoading: userVideosLoading,
    error: userVideosError,
  } = useUserGeneratedVideos();

  // Subscription state
  const [isCreatingCheckout, setIsCreatingCheckout] = useState(false);
  const [isCreatingPortal, setIsCreatingPortal] = useState(false);

  // Video preview state
  const [selectedVideo, setSelectedVideo] = useState<string | null>(null);

  // Video generation state
  const [selectedAvatarVideoId, setSelectedAvatarVideoId] = useState<
    string | null
  >(null);
  const [overlayText, setOverlayText] = useState("");
  const [isGeneratingVideo, setIsGeneratingVideo] = useState(false);

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

  const handleUpgradeToPro = async () => {
    setIsCreatingCheckout(true);
    try {
      const response = await api.subscriptions.createCheckoutSession({
        price_id: "price_1SKOuPLa4pEqShgojlivZTLc", // Your Stripe price ID
        success_url: `${window.location.origin}/dashboard?success=true`,
        cancel_url: `${window.location.origin}/dashboard?canceled=true`,
      });

      // Redirect to Stripe Checkout
      window.location.href = response.checkout_url;
    } catch (error) {
      console.error("Error creating checkout session:", error);
      alert("Failed to create checkout session. Please try again.");
    } finally {
      setIsCreatingCheckout(false);
    }
  };

  const handleManageSubscription = async () => {
    setIsCreatingPortal(true);
    try {
      const response = await api.subscriptions.createCustomerPortalSession({
        return_url: `${window.location.origin}/dashboard`,
      });

      // Redirect to Stripe Customer Portal
      window.location.href = response.portal_url;
    } catch (error) {
      console.error("Error creating customer portal session:", error);
      alert("Failed to open customer portal. Please try again.");
    } finally {
      setIsCreatingPortal(false);
    }
  };

  const handleGenerateVideo = async () => {
    if (!selectedAvatarVideoId || !overlayText.trim()) {
      alert("Please select a video and enter some text");
      return;
    }

    setIsGeneratingVideo(true);
    try {
      await api.userGeneratedVideos.createUserGeneratedVideo({
        ai_avatar_video_id: selectedAvatarVideoId,
        overlay_text: overlayText.trim(),
      });

      // Show success message
      alert("Video generated successfully!");

      // Clear the form
      setOverlayText("");
      setSelectedAvatarVideoId(null);

      // Refresh the user-generated videos list
      await queryClient.invalidateQueries({
        queryKey: CACHE_KEYS.USER_GENERATED_VIDEOS,
      });
    } catch (error) {
      console.error("Error generating video:", error);
      alert("Failed to generate video. Please try again.");
    } finally {
      setIsGeneratingVideo(false);
    }
  };

  const accessToken = session?.access_token;

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-4xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold">Dashboard</h1>
          <Button onClick={handleSignOut} variant="outline">
            Sign Out
          </Button>
        </div>

        <div className="space-y-6">
          {/* User Info Card */}
          <Card>
            <CardHeader>
              <CardTitle>Welcome to Reel Farm!</CardTitle>
              <CardDescription>
                You're successfully authenticated
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                <p>
                  <strong>Email:</strong> {user?.email}
                </p>
                <p>
                  <strong>User ID:</strong> {user?.id}
                </p>
                <p>
                  <strong>Last Sign In:</strong>{" "}
                  {user?.last_sign_in_at
                    ? new Date(user.last_sign_in_at).toLocaleString()
                    : "N/A"}
                </p>
                <p>
                  <strong>Access Token:</strong>
                </p>
                <div className="bg-gray-100 p-2 rounded border font-mono text-sm break-all">
                  {accessToken || "N/A"}
                </div>
              </div>
            </CardContent>
          </Card>

          {/* API Health Card */}
          <Card>
            <CardHeader>
              <CardTitle>API Health</CardTitle>
            </CardHeader>
            <CardContent>
              {healthLoading && (
                <p className="text-blue-600">Checking API health...</p>
              )}
              {healthError && (
                <p className="text-red-600">Error: {healthError.message}</p>
              )}
              {health && (
                <div className="bg-green-50 p-3 rounded border">
                  <p>
                    <strong>Status:</strong> {health.status}
                  </p>
                  <p>
                    <strong>Message:</strong> {health.message}
                  </p>
                  <p>
                    <strong>Port:</strong> {health.port}
                  </p>
                </div>
              )}
            </CardContent>
          </Card>

          {/* User Account Card */}
          <Card>
            <CardHeader>
              <CardTitle>User Account (from API)</CardTitle>
            </CardHeader>
            <CardContent>
              {userLoading && (
                <p className="text-blue-600">Loading user account...</p>
              )}
              {userError && (
                <p className="text-red-600">Error: {userError.message}</p>
              )}
              {userAccount && (
                <div className="bg-green-50 p-3 rounded border">
                  <p>
                    <strong>Plan:</strong> {userAccount.plan}
                  </p>
                  <p>
                    <strong>Credits:</strong> {userAccount.credits || 0}
                  </p>
                  <p>
                    <strong>Plan Started:</strong>{" "}
                    {new Date(userAccount.plan_started_at).toLocaleDateString()}
                  </p>
                  <p>
                    <strong>Created:</strong>{" "}
                    {new Date(userAccount.created_at).toLocaleDateString()}
                  </p>
                  {userAccount.plan_ends_at && (
                    <p>
                      <strong>Plan Ends:</strong>{" "}
                      {new Date(userAccount.plan_ends_at).toLocaleDateString()}
                    </p>
                  )}
                </div>
              )}
            </CardContent>
          </Card>

          {/* Subscription Card */}
          <Card>
            <CardHeader>
              <CardTitle>Subscription</CardTitle>
              <CardDescription>
                Upgrade to Pro for more credits and features
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="bg-blue-50 p-4 rounded-lg border">
                  <h3 className="font-semibold text-lg mb-2">Reel Farm Pro</h3>
                  <ul className="space-y-1 text-sm">
                    <li>• 500 credits per month</li>
                    <li>• Credits never expire</li>
                    <li>• Priority support</li>
                    <li>• Advanced features</li>
                  </ul>
                  <div className="mt-3">
                    <span className="text-2xl font-bold">£0.00</span>
                    <span className="text-gray-600">/month</span>
                  </div>
                </div>

                {userAccount?.plan === "free" ? (
                  <Button
                    className="w-full"
                    size="lg"
                    onClick={handleUpgradeToPro}
                    disabled={isCreatingCheckout}
                  >
                    {isCreatingCheckout
                      ? "Creating Checkout..."
                      : "Upgrade to Pro"}
                  </Button>
                ) : (
                  <div className="text-center">
                    <p className="text-green-600 font-medium">
                      ✓ You're on the Pro plan!
                    </p>
                    <Button
                      variant="outline"
                      className="mt-2"
                      onClick={handleManageSubscription}
                      disabled={isCreatingPortal}
                    >
                      {isCreatingPortal
                        ? "Opening Portal..."
                        : "Manage Subscription"}
                    </Button>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>

          {/* Hooks Management */}
          <HooksManager />
          {/* AI Avatar Videos Card */}
          <Card>
            <CardHeader>
              <CardTitle>AI Avatar Videos</CardTitle>
              <CardDescription>
                Click on a thumbnail to watch the video, or select one to add
                your own text overlay
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
                      onClick={() => setSelectedVideo(video.video_url)}
                    >
                      <div className="aspect-square bg-gray-100 rounded-lg overflow-hidden">
                        <img
                          src={video.thumbnail_url}
                          alt={video.title}
                          className="w-full h-full object-cover hover:scale-105 transition-transform duration-200"
                        />
                      </div>
                      <div className="p-2">
                        <p
                          className="text-xs text-gray-600 truncate"
                          title={video.title}
                        >
                          {video.title}
                        </p>
                        <Button
                          size="sm"
                          variant="outline"
                          className="w-full mt-1 text-xs"
                          onClick={(e) => {
                            e.stopPropagation();
                            setSelectedAvatarVideoId(video.id);
                          }}
                        >
                          {selectedAvatarVideoId === video.id
                            ? "Selected"
                            : "Select"}
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
              {aiAvatarVideos && aiAvatarVideos.videos.length === 0 && (
                <div className="text-center py-8 text-gray-500">
                  <p>No AI avatar videos available yet.</p>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Video Generation Form */}
          {selectedAvatarVideoId && (
            <Card>
              <CardHeader>
                <CardTitle>Generate Video with Text Overlay</CardTitle>
                <CardDescription>
                  Add your own text to the selected AI avatar video
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div>
                    <Label htmlFor="overlay-text">
                      Text to overlay on video
                    </Label>
                    <Input
                      id="overlay-text"
                      placeholder="Enter your text here..."
                      value={overlayText}
                      onChange={(e) => setOverlayText(e.target.value)}
                      maxLength={500}
                    />
                    <p className="text-xs text-gray-500 mt-1">
                      {overlayText.length}/500 characters
                    </p>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      onClick={handleGenerateVideo}
                      disabled={isGeneratingVideo || !overlayText.trim()}
                      className="flex-1"
                    >
                      {isGeneratingVideo ? "Generating..." : "Generate Video"}
                    </Button>
                    <Button
                      variant="outline"
                      onClick={() => {
                        setSelectedAvatarVideoId(null);
                        setOverlayText("");
                      }}
                    >
                      Cancel
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Video Preview Section */}
          {selectedVideo && (
            <Card>
              <CardHeader>
                <CardTitle>Video Preview</CardTitle>
                <CardDescription>
                  Click on any thumbnail above to preview the video
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="aspect-[9/16] bg-black rounded-lg overflow-hidden max-w-sm mx-auto">
                  <video
                    src={selectedVideo}
                    controls
                    className="w-full h-full"
                    autoPlay
                  >
                    Your browser does not support the video tag.
                  </video>
                </div>
              </CardContent>
            </Card>
          )}

          {/* User Generated Videos Card */}
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
              {userGeneratedVideos && (
                <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
                  {userGeneratedVideos.videos.map((video) => (
                    <div
                      key={video.id}
                      className="bg-white rounded-lg border shadow-sm hover:shadow-md transition-shadow cursor-pointer"
                      onClick={() => setSelectedVideo(video.video_url)}
                    >
                      <div className="aspect-square bg-gray-100 rounded-lg overflow-hidden">
                        <img
                          src={video.thumbnail_url}
                          alt={`Video with text: ${video.overlay_text}`}
                          className="w-full h-full object-cover hover:scale-105 transition-transform duration-200"
                        />
                      </div>
                      <div className="p-2">
                        <p
                          className="text-xs text-gray-600 truncate"
                          title={video.overlay_text}
                        >
                          "{video.overlay_text}"
                        </p>
                        <div className="flex items-center justify-between mt-1">
                          <span
                            className={`text-xs px-2 py-1 rounded-full ${
                              video.status === "completed"
                                ? "bg-green-100 text-green-800"
                                : video.status === "processing"
                                ? "bg-yellow-100 text-yellow-800"
                                : "bg-red-100 text-red-800"
                            }`}
                          >
                            {video.status}
                          </span>
                          <span className="text-xs text-gray-500">
                            {new Date(video.created_at).toLocaleDateString()}
                          </span>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
              {userGeneratedVideos &&
                userGeneratedVideos.videos.length === 0 && (
                  <div className="text-center py-8 text-gray-500">
                    <p>You haven't created any videos yet.</p>
                    <p className="text-sm mt-2">
                      Select an AI avatar video above and add your own text!
                    </p>
                  </div>
                )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
