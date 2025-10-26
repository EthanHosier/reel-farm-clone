import { Video } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ROUTES } from "@/types/routes";

export function EmptyState() {
  return (
    <div className="flex flex-col items-center justify-center py-12 text-center min-h-[60vh]">
      <div className="flex h-20 w-20 items-center justify-center rounded-full bg-gray-100 mb-4">
        <Video className="h-10 w-10 text-gray-400" />
      </div>
      <h3 className="text-lg font-semibold mb-2">No Generated Videos Yet</h3>
      <p className="text-gray-500 max-w-sm mb-6">
        You haven't created any videos yet. Get started by selecting an AI
        avatar video and generating your first custom video!
      </p>
      <Button asChild>
        <a href={ROUTES.generateAiAvatarVideo}>Generate Your First Video</a>
      </Button>
    </div>
  );
}
