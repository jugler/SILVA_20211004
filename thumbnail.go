package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

//use ffmpeg to generate thumbnails of different sizes
func generateThumbnails(videoPath string, size int) []byte {
	var filter = fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease", size, size)
	var outputName = fmt.Sprintf("thumbnails/%dx%d.jpg", size, size)

	cmd := exec.Command("ffmpeg", "-y", "-ss", "00:00:01.00", "-i", videoPath, "-vf", filter, "-vframes", "1", outputName)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	contents, err := os.ReadFile(outputName)
	return contents
}
