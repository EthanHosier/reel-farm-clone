import { useState, useEffect } from "react";
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
import { useHooks } from "@/hooks/useGetHooks";
import { useCreateUserGeneratedVideo } from "@/features/videos/generate-ai-avatar-video/queries/useCreateUserGeneratedVideo";

interface VideoGenerationFormProps {
  selectedAvatarVideoId: string | null;
  onCancel: () => void;
  onTextChange: (text: string) => void;
}

export function VideoGenerationForm({
  selectedAvatarVideoId,
  onCancel,
  onTextChange,
}: VideoGenerationFormProps) {
  const [currentHookIndex, setCurrentHookIndex] = useState(0);
  const [editableText, setEditableText] = useState("");

  const { data: hooks, isLoading: hooksLoading } = useHooks(50, 0);

  const createVideoMutation = useCreateUserGeneratedVideo({
    onSuccess: () => {
      alert("Video generated successfully!");
      onCancel();
    },
    onError: (error) => {
      alert(`Error generating video: ${error.message}`);
    },
  });

  // Initialize editable text when hooks load
  useEffect(() => {
    if (hooks && hooks.hooks.length > 0) {
      setEditableText(hooks.hooks[currentHookIndex]?.text || "");
    }
  }, [hooks, currentHookIndex]);

  // Hook navigation functions
  const handlePreviousHook = () => {
    if (hooks && hooks.hooks.length > 0) {
      const newIndex =
        currentHookIndex > 0 ? currentHookIndex - 1 : hooks.hooks.length - 1;
      setCurrentHookIndex(newIndex);
    }
  };

  const handleNextHook = () => {
    if (hooks && hooks.hooks.length > 0) {
      const newIndex =
        currentHookIndex < hooks.hooks.length - 1 ? currentHookIndex + 1 : 0;
      setCurrentHookIndex(newIndex);
    }
  };

  const handleTextChange = (newText: string) => {
    setEditableText(newText);
    onTextChange(newText);
  };

  const handleGenerateVideo = () => {
    if (!selectedAvatarVideoId || !editableText.trim()) return;

    createVideoMutation.mutate({
      ai_avatar_video_id: selectedAvatarVideoId,
      overlay_text: editableText.trim(),
    });
  };

  if (!selectedAvatarVideoId) return null;

  return (
    <Card>
      <CardHeader>
        <CardTitle>Generate Video with Text Overlay</CardTitle>
        <CardDescription>
          Navigate through your hooks and edit the text before generating
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {/* Hook Navigation */}
          {hooks && hooks.hooks.length > 0 && (
            <div className="space-y-2">
              <Label>
                Select Hook ({currentHookIndex + 1} of {hooks.hooks.length})
              </Label>
              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handlePreviousHook}
                  disabled={hooksLoading}
                >
                  ← Previous
                </Button>
                <div className="flex-1 text-center text-sm text-gray-600">
                  {hooks.hooks[currentHookIndex]?.text}
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleNextHook}
                  disabled={hooksLoading}
                >
                  Next →
                </Button>
              </div>
            </div>
          )}

          {/* Text Editor */}
          <div>
            <Label htmlFor="overlay-text">Edit text for video overlay</Label>
            <Input
              id="overlay-text"
              placeholder="Enter your text here..."
              value={editableText}
              onChange={(e) => handleTextChange(e.target.value)}
              maxLength={500}
            />
            <p className="text-xs text-gray-500 mt-1">
              {editableText.length}/500 characters
            </p>
          </div>
          <div className="flex gap-2">
            <Button
              onClick={handleGenerateVideo}
              disabled={createVideoMutation.isPending || !editableText.trim()}
              className="flex-1"
            >
              {createVideoMutation.isPending
                ? "Generating..."
                : "Generate Video"}
            </Button>
            <Button variant="outline" onClick={onCancel}>
              Cancel
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
