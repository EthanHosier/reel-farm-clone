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

	fmt.Printf("📝 Adding text: '%s'\n", text)
	fmt.Printf("💾 Output file: %s\n", outputFile)

	// FFmpeg command to add text overlay with border outline
	// -i: input video
	// -vf: video filter with drawtext
	// drawtext parameters:
	//   text='Sample Text': the text to display
	//   fontfile=TikTokDisplay-Medium.ttf: custom font file
	//   fontsize=48: font size
	//   fontcolor=white: text color
	//   x=(w-text_w)/2: center horizontally
	//   y=(h-text_h)/2: center vertically
	//   borderw=3: border width around text
	//   bordercolor=black: black border color
	cmd := exec.Command("ffmpeg",
		"-i", firstVideo,
		"-vf", fmt.Sprintf("drawtext=text='%s':fontfile=TikTokDisplay-Medium.ttf:fontsize=48:fontcolor=white:x=(w-text_w)/2:y=(h-text_h)/2:borderw=5:bordercolor=black", text),
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
