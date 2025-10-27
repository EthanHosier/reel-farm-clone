import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
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
import { useAuth } from "@/contexts/AuthContext";

const NUM_HOOKS = 5;
const SUGGESTIONS = [
  "Fun facts about cooking",
  "Money-saving tips",
  "Life hacks everyone should know",
  "Relationship advice",
  "Health and wellness tips",
];

export function GenerateHooks() {
  const { user } = useAuth();
  const [prompt, setPrompt] = useState("");
  const [limit] = useState(20);
  const [offset] = useState(0);
  const [selectedSuggestion, setSelectedSuggestion] = useState<string | null>(
    null
  );

  // Queries and mutations
  const {
    data: hooksData,
    isLoading: hooksLoading,
    error: hooksError,
  } = useHooks(limit, offset);
  const { mutateAsync: generateHooks, isPending: generateHooksPending } =
    useGenerateHooks();
  const deleteHook = useDeleteHook();

  const handleGenerateHooks = async () => {
    if (!prompt.trim()) {
      alert("Please enter a prompt");
      return;
    }

    try {
      await generateHooks({
        prompt: prompt.trim(),
        num_hooks: NUM_HOOKS,
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

  const firstName = user?.email?.split("@")[0] || "there";

  return (
    <div className="space-y-8">
      {/* Main Prompt Section */}
      <div className="flex flex-col items-center justify-center min-h-[60vh] space-y-8">
        {/* Greeting */}
        <h2 className="text-3xl font-light">
          Ready to create viral hooks,{" "}
          {firstName.charAt(0).toUpperCase() + firstName.slice(1)}?
        </h2>

        {/* Central Input */}
        <div className="w-full max-w-2xl relative">
          <div className="relative flex items-center bg-white rounded-full border-2 border-gray-200 shadow-lg hover:shadow-xl transition-shadow px-2 justify-between ">
            <Input
              id="prompt"
              disabled={generateHooksPending}
              placeholder="Ask anything"
              value={prompt}
              onChange={(e) => setPrompt(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter" && prompt.trim()) {
                  handleGenerateHooks();
                }
              }}
              className="mx-2 w-full pr-20 py-6 text-lg border-0 focus-visible:ring-0 focus-visible:ring-offset-0"
            />
            <Button
              className="py-2 rounded-full"
              onClick={handleGenerateHooks}
              disabled={generateHooksPending}
            >
              {generateHooksPending ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                "Generate"
              )}
            </Button>
          </div>

          {/* Suggestions */}
          <div className="mt-6">
            <div className="flex flex-wrap gap-2 justify-center">
              {SUGGESTIONS.map((suggestion) => (
                <Button
                  size="sm"
                  variant="outline"
                  key={suggestion}
                  onClick={() => {
                    setPrompt(suggestion);
                    setSelectedSuggestion(suggestion);
                  }}
                  className={`px-4 py-2 rounded-full text-sm transition ${
                    selectedSuggestion === suggestion
                      ? "bg-black text-white"
                      : " text-gray-700 hover:bg-gray-100"
                  }`}
                >
                  {suggestion}
                </Button>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Hooks List Section */}
      <div className="space-y-4">
        <div>
          <h2 className="text-2xl font-bold mb-2">Generated Hooks</h2>
        </div>
        <div>
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
                            variant="ghost"
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
        </div>
      </div>
    </div>
  );
}
