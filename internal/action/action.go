package action

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/BrunoTeixeira1996/gocam/internal/utils"
)

// ffmpeg -i rtsp://brun0teixeira:qwerty654321@192.168.30.44:554/stream1 -c:v copy -c:a aac -strict experimental output.mp4
func StartFFMPEGRecording(rD time.Duration, recordingDuration string, cancel chan struct{}, dumpLocation string, logLocation string, recordID string) error {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02-15-04-05")
	outputFile := dumpLocation + "/output" + formattedTime + "-" + recordID + ".mp4"

	log.Printf("[INFO]Starting record duration %s for %s file with %s ID\n", recordingDuration, outputFile, recordID)

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
		log.Printf("[INFO] Recording for %s ID finished\n", recordID)
	case <-cancel:
		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			log.Printf("[ERROR] Error sending interrupt signal: %s\n", err)
		}
		// Wait for the process to exit
		<-time.After(5 * time.Second)
		if cmd.ProcessState == nil || !cmd.ProcessState.Exited() {
			log.Printf("[ERROR] Process for %s ID did not exit gracefully, force killing...\n", recordID)

			err := cmd.Process.Kill()
			if err != nil {
				log.Printf("[ERROR] Error killing process: %s\n", err)
			}
			return err
		}
		log.Printf("[INFO] Recording for %s ID cancelled\n", recordID)
	}

	// Saves a file (recordID.log) in the logs folder to not populate the log stdout
	if err := utils.SaveFFMPEGOutput(logLocation, recordID, output.Bytes()); err != nil {
		return err
	}
	// TODO: finish this log sentence
	log.Printf("[INFO] Finished record for %s ID and created log file in ...\n", recordID)

	return nil
}
