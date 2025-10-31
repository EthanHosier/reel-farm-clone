import { useRef, useEffect, useState } from "react";

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
      const testLine = currentLine === "" ? word : currentLine + " " + word;

      if (testLine.length <= maxCharsPerLine) {
        currentLine = testLine;
      } else {
        if (currentLine !== "") {
          lines.push(currentLine);
          currentLine = word;
        } else {
          lines.push(word);
          currentLine = "";
        }
      }
    }

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
        const fontSize = Math.floor((containerHeight / 1280) * 36);
        setCalculatedFontSize(fontSize);
      }
    };

    calculateFontSize();
    const handleResize = () => calculateFontSize();
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, [selectedVideo]);

  if (!selectedVideo) {
    return (
      <div className="aspect-[9/16] bg-gray-200 rounded-lg w-full max-w-sm flex items-center justify-center">
        <p className="text-gray-400">Select an avatar to preview</p>
      </div>
    );
  }

  return (
    <div className="w-full">
      <div
        ref={videoContainerRef}
        className="aspect-[9/16] bg-black rounded-lg overflow-hidden w-full max-w-[340px] mx-auto relative"
      >
        <video
          src={selectedVideo}
          className="w-full h-full object-cover video-hide-controls"
          autoPlay
          loop
          muted
          playsInline
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
    </div>
  );
}
