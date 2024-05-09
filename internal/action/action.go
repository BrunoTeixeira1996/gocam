package action

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

// ffmpeg -i rtsp://brun0teixeira:qwerty654321@192.168.30.44:554/stream1 -c:v copy -c:a aac -strict experimental output.mp4
func StartFFMPEGRecording(recordingDuration string) error {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02-15-04-05")
	outputFile := "output" + formattedTime + ".mp4"

	log.Printf("Starting record duration: %s\n", recordingDuration)

	cmd := exec.Command("ffmpeg", "-i", "rtsp://brun0teixeira:qwerty654321@192.168.30.44:554/stream1", "-c:v", "copy", "-c:a", "aac", "-strict", "experimental", "-t", recordingDuration, "output.mp4")

	// Capture the output
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("[ERROR] FFMPEG while running:%s\n", err)
	}

	log.Printf("FFMPEG Output:%s\n", output.String())
	log.Printf("Finished record\n")

	log.Printf("Changing file name to %s\n========================\n\n", outputFile)

	if err := os.Rename("output.mp4", outputFile); err != nil {
		return fmt.Errorf("[ERROR] renaming file:%s\n", err)
	}

	return nil
}
