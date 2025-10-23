package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// Get the videos directory path
	videosDir := "videos"
	if len(os.Args) > 1 {
		videosDir = os.Args[1]
	}

	// Check if videos directory exists
	if _, err := os.Stat(videosDir); os.IsNotExist(err) {
		fmt.Printf("❌ Videos directory '%s' does not exist\n", videosDir)
		os.Exit(1)
	}

	// Find the first video file
	var firstVideo string
	err := filepath.Walk(videosDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if it's a video file
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".mp4" || ext == ".mov" || ext == ".avi" || ext == ".mkv" {
			firstVideo = path
			return filepath.SkipDir // Stop walking after finding first video
		}
		return nil
	})

	if err != nil {
		fmt.Printf("❌ Error walking videos directory: %v\n", err)
		os.Exit(1)
	}

	if firstVideo == "" {
		fmt.Printf("❌ No video files found in '%s'\n", videosDir)
		os.Exit(1)
	}

	fmt.Printf("🎬 Found first video: %s\n", firstVideo)

	// Get text to overlay (default or from command line)
	text := "Sample Text"
	if len(os.Args) > 2 {
		text = os.Args[2]
	}

	// Create output filename
	baseName := strings.TrimSuffix(filepath.Base(firstVideo), filepath.Ext(firstVideo))
	outputFile := fmt.Sprintf("%s_with_text.mp4", baseName)

	// Wrap text if it's too long (approximately 20 characters per line for 48px font)
	wrappedLines := wrapTextToLines(text, 20)

	fmt.Printf("📝 Adding text (%d lines):\n", len(wrappedLines))
	for i, line := range wrappedLines {
		fmt.Printf("  Line %d: '%s'\n", i+1, line)
	}
	fmt.Printf("💾 Output file: %s\n", outputFile)

	// Create a temporary text file with the wrapped text
	tempTextFile := "temp_text.txt"
	joinedText := strings.Join(wrappedLines, "\n")
	err = os.WriteFile(tempTextFile, []byte(joinedText), 0644)
	if err != nil {
		fmt.Printf("❌ Failed to create temporary text file: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(tempTextFile) // Clean up the temp file

	// Use textfile parameter instead of inline text
	videoFilter := fmt.Sprintf("drawtext=textfile=%s:fontfile=TikTokDisplay-Medium.ttf:fontsize=48:fontcolor=white:x=(w-text_w)/2:y=(h-text_h)/2:borderw=5:bordercolor=black:text_align=center", tempTextFile)

	cmd := exec.Command("ffmpeg",
		"-i", firstVideo,
		"-vf", videoFilter,
		"-c:a", "copy", // Copy audio without re-encoding
		"-y", // Overwrite output file if it exists
		outputFile,
	)

	fmt.Printf("🚀 Running FFmpeg command...\n")
	fmt.Printf("Command: %s\n", strings.Join(cmd.Args, " "))

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ FFmpeg error: %v\n", err)
		fmt.Printf("Output: %s\n", string(output))
		os.Exit(1)
	}

	fmt.Printf("✅ Successfully created %s\n", outputFile)
	fmt.Printf("📁 Output location: %s\n", filepath.Join(".", outputFile))
}

// wrapTextToLines wraps text to fit within a specified number of characters per line
func wrapTextToLines(text string, maxCharsPerLine int) []string {
	words := strings.Fields(text)
	var lines []string
	var currentLine string

	for _, word := range words {
		if len(currentLine)+len(word)+1 <= maxCharsPerLine {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
				currentLine = word
			} else {
				// Word is longer than maxCharsPerLine, add it anyway
				lines = append(lines, word)
			}
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}
