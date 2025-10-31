import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useHooks } from "@/hooks/useGetHooks";
import { ChevronLeft, ChevronRight } from "lucide-react";

const PLACEHOLDER_TEXT = "I wish someone told me this sooner...";

interface HookSectionProps {
  onTextChange: (text: string) => void;
}

export function HookSection({ onTextChange }: HookSectionProps) {
  const [currentHookIndex, setCurrentHookIndex] = useState(0);
  const [editableText, setEditableText] = useState(PLACEHOLDER_TEXT);

  const { data: hooks, isLoading: hooksLoading } = useHooks(50, 0);

  // Initialize editable text when hooks load
  useEffect(() => {
    if (hooks && hooks.hooks.length > 0) {
      setEditableText(hooks.hooks[currentHookIndex]?.text || "");
    }
  }, [hooks, currentHookIndex]);

  // Handle text change and notify parent
  useEffect(() => {
    onTextChange(editableText);
  }, [editableText, onTextChange]);

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

  const hookCount = hooks?.hooks.length || 0;

  return (
    <div>
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-sm font-medium">1. Hook</h3>
        {hookCount > 0 && (
          <span className="text-xs text-gray-500">
            {currentHookIndex + 1}/{hookCount}
          </span>
        )}
      </div>

      <div className="flex items-center gap-2">
        <Button
          variant="ghost"
          size="icon"
          onClick={handlePreviousHook}
          disabled={hooksLoading || hookCount === 0}
          className="h-8 w-8"
        >
          <ChevronLeft className="h-4 w-4" />
        </Button>
        <Input
          value={editableText}
          onChange={(e) => setEditableText(e.target.value)}
          maxLength={500}
          className="flex-1 bg-defaultbg"
        />
        <Button
          variant="ghost"
          size="icon"
          onClick={handleNextHook}
          disabled={hooksLoading || hookCount === 0}
          className="h-8 w-8"
        >
          <ChevronRight className="h-4 w-4" />
        </Button>
      </div>
    </div>
  );
}
