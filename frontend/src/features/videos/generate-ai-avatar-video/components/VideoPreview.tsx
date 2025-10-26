import { useRef, useEffect, useState } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

interface VideoPreviewProps {
  selectedVideo: string | null;
  overlayText: string;
}

export function VideoPreview({
  selectedVideo,
  overlayText,
}: VideoPreviewProps) {
  const [calculatedFontSize, setCalculatedFontSize] = useState(18);
  const videoContainerRef = useRef<HTMLDivElement>(null);

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

  // Calculate font size based on video container dimensions
  useEffect(() => {
    const calculateFontSize = () => {
      if (videoContainerRef.current) {
        const containerHeight = videoContainerRef.current.offsetHeight;
        // FFmpeg uses 36px font for 1280px height video
        // So font size = (container height / 1280) * 36
        const fontSize = Math.floor((containerHeight / 1280) * 36);
        setCalculatedFontSize(fontSize);
      }
    };

    calculateFontSize();

    // Recalculate on window resize
    const handleResize = () => calculateFontSize();
    window.addEventListener("resize", handleResize);

    return () => window.removeEventListener("resize", handleResize);
  }, [selectedVideo]);

  if (!selectedVideo) return null;

  return (
    <Card>
      <CardHeader>
        <CardTitle>Video Preview</CardTitle>
        <CardDescription>
          Click on any thumbnail above to preview the video. Text overlay
          preview appears when you type in the generation form below.
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
            className="w-full h-full video-hide-controls"
            autoPlay
            loop
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
                      fontFamily: '"TikTokDisplay-Medium", Arial, sans-serif',
                      lineHeight: "1.4",
                      fontSize: `${calculatedFontSize}px`,
                      textShadow:
                        "1px 1px 0px black, 1px -1px 0px black, -1px 1px 0px black, -1px -1px 0px black",
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
  );
}
