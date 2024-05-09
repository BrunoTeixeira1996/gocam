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
func StartFFMPEGRecording(rD time.Duration, recordingDuration string, cancel chan struct{}) error {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02-15-04-05")
	outputFile := "output" + formattedTime + ".mp4"

	log.Printf("Starting record duration %s for %s file\n", recordingDuration, outputFile)

	cmd := exec.Command("ffmpeg", "-i", "rtsp://brun0teixeira:qwerty654321@192.168.30.44:554/stream1", "-c:v", "copy", "-c:a", "aac", "-strict", "experimental", "-t", recordingDuration, outputFile)

	// Capture the output
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("[ERROR] FFMPEG while starting:%s\n", err)
	}

	select {
	case <-time.After(rD):
		log.Printf("Recording for camera %s finished\n", "camera1")
	case <-cancel:
		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			log.Printf("Error sending interrupt signal: %s\n", err)
		}
		// Wait for the process to exit
		<-time.After(5 * time.Second)
		if cmd.ProcessState == nil || !cmd.ProcessState.Exited() {
			log.Printf("[ERROR] Process for camera %s did not exit gracefully, force killing...\n", "camera1")

			err := cmd.Process.Kill()
			if err != nil {
				log.Printf("[ERROR] Error killing process: %s\n", err)
			}
			return err
		}
		log.Printf("[INFO] Recording for camera %s cancelled\n", "camera1")
	}

	log.Printf("[INFO] FFMPEG Output:%s\n", output.String())
	log.Printf("[INFO] Finished record\n")

	return nil
}
