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
import { useCreateUserGeneratedVideo } from "./queries/useCreateUserGeneratedVideo";
import { useSubscriptionMutations } from "./queries/useSubscriptionMutations";
import React, { useState, useRef, useEffect } from "react";

export default function Dashboard() {
  // Text wrapping function to match FFmpeg behavior
  const wrapTextToLines = (text: string, maxCharsPerLine: number): string[] => {
    const words = text.split(" ");
    const lines: string[] = [];
    let currentLine = "";

    for (const word of words) {
      // Check if adding this word would exceed the limit
      const testLine = currentLine === "" ? word : currentLine + " " + word;

      if (testLine.length <= maxCharsPerLine) {
        currentLine = testLine;
      } else {
        // If current line is not empty, push it and start a new line
        if (currentLine !== "") {
          lines.push(currentLine);
          currentLine = word;
        } else {
          // If even a single word exceeds the limit, push it anyway
          lines.push(word);
          currentLine = "";
        }
      }
    }

    // Add the last line if it's not empty
    if (currentLine !== "") {
      lines.push(currentLine);
    }

    return lines;
  };

  const { user, session, signOut } = useAuth();
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

  // Video generation mutation
  const createVideoMutation = useCreateUserGeneratedVideo({
    onSuccess: () => {
      alert("Video generated successfully!");
      // Clear the form
      setOverlayText("");
      setSelectedAvatarVideoId(null);
    },
    onError: () => {
      alert("Failed to generate video. Please try again.");
    },
  });

  // Subscription mutations
  const subscriptionMutations = useSubscriptionMutations({
    onCheckoutSuccess: (checkoutUrl) => {
      window.location.href = checkoutUrl;
    },
    onCheckoutError: () => {
      alert("Failed to create checkout session. Please try again.");
    },
    onPortalSuccess: (portalUrl) => {
      window.location.href = portalUrl;
    },
    onPortalError: () => {
      alert("Failed to open customer portal. Please try again.");
    },
  });

  // Video preview state
  const [selectedVideo, setSelectedVideo] = useState<string | null>(null);

  // Video generation state
  const [selectedAvatarVideoId, setSelectedAvatarVideoId] = useState<
    string | null
  >(null);
  const [overlayText, setOverlayText] = useState("");
  const [calculatedFontSize, setCalculatedFontSize] = useState(18);
  const videoContainerRef = useRef<HTMLDivElement>(null);

  // Calculate font size based on video container dimensions
  useEffect(() => {
    const calculateFontSize = () => {
      if (videoContainerRef.current) {
        const containerHeight = videoContainerRef.current.offsetHeight;
        // FFmpeg uses 36px font for 1280px height video
        // So font size = (container height / 1280) * 36
        const fontSize = Math.floor((containerHeight / 1280) * 36);
        setCalculatedFontSize(fontSize); // Clamp between 12-24px
      }
    };

    calculateFontSize();

    // Recalculate on window resize
    const handleResize = () => calculateFontSize();
    window.addEventListener("resize", handleResize);

    return () => window.removeEventListener("resize", handleResize);
  }, [selectedVideo]); // Recalculate when video changes

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

  const handleUpgradeToPro = () => {
    subscriptionMutations.createCheckout.mutate({
      price_id: "price_1SKOuPLa4pEqShgojlivZTLc", // Your Stripe price ID
      success_url: `${window.location.origin}/dashboard?success=true`,
      cancel_url: `${window.location.origin}/dashboard?canceled=true`,
    });
  };

  const handleManageSubscription = () => {
    subscriptionMutations.createPortal.mutate({
      return_url: `${window.location.origin}/dashboard`,
    });
  };

  const handleGenerateVideo = async () => {
    if (!selectedAvatarVideoId || !overlayText.trim()) {
      alert("Please select a video and enter some text");
      return;
    }

    createVideoMutation.mutate({
      ai_avatar_video_id: selectedAvatarVideoId,
      overlay_text: overlayText.trim(),
    });
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
                    disabled={subscriptionMutations.createCheckout.isPending}
                  >
                    {subscriptionMutations.createCheckout.isPending
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
                      disabled={subscriptionMutations.createPortal.isPending}
                    >
                      {subscriptionMutations.createPortal.isPending
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
                Click on a thumbnail to select it for video generation with your
                own text overlay
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
                      onClick={() => {
                        setSelectedAvatarVideoId(video.id);
                        setSelectedVideo(video.video_url);
                      }}
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
                      disabled={
                        createVideoMutation.isPending || !overlayText.trim()
                      }
                      className="flex-1"
                    >
                      {createVideoMutation.isPending
                        ? "Generating..."
                        : "Generate Video"}
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
                  Click on any thumbnail above to preview the video. Text
                  overlay preview appears when you type in the generation form
                  below.
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div
                  ref={videoContainerRef}
                  className="aspect-[9/16] bg-black rounded-lg overflow-hidden max-w-sm mx-auto relative"
                >
                  <video
                    src={selectedVideo}
                    controls
                    className="w-full h-full"
                    autoPlay
                  >
                    Your browser does not support the video tag.
                  </video>
                  {overlayText && (
                    <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
                      <div className="text-center px-4">
                        {wrapTextToLines(overlayText, 35).map((line, index) => (
                          <div
                            key={index}
                            className="text-white leading-tight"
                            style={{
                              fontFamily:
                                '"TikTokDisplay-Medium", Arial, sans-serif',
                              lineHeight: "1.4",
                              fontSize: `${calculatedFontSize}px`,
                              textShadow:
                                "1px 1px 1px black, 1px -1px 1px black, -1px 1px 1px black, -1px -1px 1px black",
                              // WebkitTextStroke: "1px black",
                            }}
                          >
                            {line}
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
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
