import { useState } from "react";
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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  useHooks,
  useGenerateHooks,
  useDeleteHook,
} from "@/features/videos/generate-ai-avatar-video/queries/useHooks";
import { Trash2, Loader2 } from "lucide-react";
import type { Hook } from "@/api";

export function HooksManager() {
  const [prompt, setPrompt] = useState("");
  const [numHooks, setNumHooks] = useState(3);
  const [limit] = useState(20);
  const [offset] = useState(0);

  // Queries and mutations
  const {
    data: hooksData,
    isLoading: hooksLoading,
    error: hooksError,
  } = useHooks(limit, offset);
  const generateHooks = useGenerateHooks();
  const deleteHook = useDeleteHook();

  const handleGenerateHooks = async () => {
    if (!prompt.trim()) {
      alert("Please enter a prompt");
      return;
    }

    try {
      await generateHooks.mutateAsync({
        prompt: prompt.trim(),
        num_hooks: numHooks,
      });
      setPrompt(""); // Clear prompt after successful generation
    } catch (error) {
      console.error("Error generating hooks:", error);
      alert("Failed to generate hooks. Please try again.");
    }
  };

  const handleDeleteHook = async (hookId: string) => {
    if (!confirm("Are you sure you want to delete this hook?")) {
      return;
    }

    try {
      await deleteHook.mutateAsync(hookId);
    } catch (error) {
      console.error("Error deleting hook:", error);
      alert("Failed to delete hook. Please try again.");
    }
  };

  return (
    <div className="space-y-6">
      {/* Generate Hooks Card */}
      <Card>
        <CardHeader>
          <CardTitle>Generate Hooks</CardTitle>
          <CardDescription>
            Create TikTok hooks for your slideshow content
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="prompt">Prompt</Label>
            <Input
              id="prompt"
              placeholder="e.g., Plants dying in my house"
              value={prompt}
              onChange={(e) => setPrompt(e.target.value)}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="numHooks">Number of Hooks</Label>
            <Input
              id="numHooks"
              type="number"
              min="1"
              max="10"
              value={numHooks}
              onChange={(e) => setNumHooks(parseInt(e.target.value) || 3)}
            />
          </div>
          <Button
            onClick={handleGenerateHooks}
            disabled={generateHooks.isPending || !prompt.trim()}
            className="w-full"
          >
            {generateHooks.isPending ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Generating Hooks...
              </>
            ) : (
              "Generate Hooks"
            )}
          </Button>
        </CardContent>
      </Card>

      {/* Hooks List Card */}
      <Card>
        <CardHeader>
          <CardTitle>Your Hooks</CardTitle>
          <CardDescription>Manage your generated hooks</CardDescription>
        </CardHeader>
        <CardContent>
          {hooksLoading && (
            <div className="flex items-center justify-center py-8">
              <Loader2 className="h-6 w-6 animate-spin mr-2" />
              Loading hooks...
            </div>
          )}

          {hooksError && (
            <div className="text-red-600 text-center py-8">
              Error loading hooks: {hooksError.message}
            </div>
          )}

          {hooksData && (
            <div className="space-y-4">
              <div className="text-sm text-gray-600">
                Showing {hooksData.hooks.length} of {hooksData.total_count}{" "}
                hooks
              </div>

              {hooksData.hooks.length === 0 ? (
                <div className="text-center py-8 text-gray-500">
                  No hooks yet. Generate some hooks to get started!
                </div>
              ) : (
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Hook Text</TableHead>
                      <TableHead className="w-[100px]">Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {hooksData.hooks.map((hook: Hook) => (
                      <TableRow key={hook.id}>
                        <TableCell className="font-medium">
                          {hook.text}
                        </TableCell>
                        <TableCell>
                          <Button
                            variant="destructive"
                            size="sm"
                            onClick={() => handleDeleteHook(hook.id)}
                            disabled={deleteHook.isPending}
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              )}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
