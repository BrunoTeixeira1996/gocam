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

// Struct responsable for holding a respective recording
type Recording struct {
	Id                 string        // random string to identify the recording
	Start              string        // start date
	WantDuration       string        // recording duration from POST
	WantDurationParsed time.Duration // recording duration parsed
	WantDurationS      string        // recording duration in seconds for ffmpeg
	Cmd                string        // cmd used for ffmpeg exec
	Channel            chan struct{} // channel used in the respective recording in order to cancel the goroutine
	DumpOutput         string        // .mp4 dump
	LogOutput          string        // .log dump
}

// Initializes recording struct
func (r *Recording) Init(dumpOutput, logOutput string) {
	r.Id = utils.GenerateRandomString(10)
	r.DumpOutput = dumpOutput
	r.LogOutput = logOutput
	r.Channel = make(chan struct{})
}

// Parses duration and converts to seconds since ffmpeg needs that value
func (r *Recording) ParseDurationToSeconds() error {
	var err error

	r.WantDurationParsed, err = time.ParseDuration(r.WantDuration)
	if err != nil {
		return fmt.Errorf("While parsing duration: %s", err)
	}

	r.WantDurationS = fmt.Sprintf("%.f", r.WantDurationParsed.Seconds())

	return nil
}

// ffmpeg -i rtsp://brun0teixeira:qwerty654321@192.168.30.44:554/stream1 -c:v copy -c:a aac -strict experimental output.mp4
func StartFFMPEGRecording(recording *Recording, recordings *[]Recording) error {
	currentTime := time.Now()
	recording.Start = currentTime.Format("2006-01-02-15-04-05")
	recording.DumpOutput = recording.DumpOutput + recording.Start + "-" + recording.Id + ".mp4"

	log.Printf("[INFO] Starting record duration %s for %s file with %s ID\n", recording.WantDurationS, recording.DumpOutput, recording.Id)

	cmd := exec.Command("ffmpeg", "-i", "rtsp://brun0teixeira:qwerty654321@192.168.30.44:554/stream1", "-c:v", "copy", "-c:a", "aac", "-strict", "experimental", "-t", recording.WantDurationS, recording.DumpOutput)

	// Capture the output
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("[ERROR] FFMPEG while starting:%s\n", err)
	}

	// Add recording to slice of recordings
	*recordings = append(*recordings, *recording)

	select {
	case <-time.After(recording.WantDurationParsed):
		log.Printf("[INFO] Recording for %s ID finished\n", recording.Id)
	case <-recording.Channel:
		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			log.Printf("[ERROR] Error sending interrupt signal: %s\n", err)
		}
		// Wait for the process to exit
		<-time.After(5 * time.Second)
		if cmd.ProcessState == nil || !cmd.ProcessState.Exited() {
			log.Printf("[ERROR] Process for %s ID did not exit gracefully, force killing...\n", recording.Id)

			err := cmd.Process.Kill()
			if err != nil {
				log.Printf("[ERROR] Error killing process: %s\n", err)
			}
			return err
		}
		log.Printf("[INFO] Recording for %s ID cancelled\n", recording.Id)
	}

	// Saves a file (recording.Id.log) in the logs folder to not populate the log stdout
	if err := utils.SaveFFMPEGOutput(recording.LogOutput, recording.Id, output.Bytes()); err != nil {
		return err
	}

	log.Printf("[INFO] Finished recording for %s ID log file at %s\n", recording.Id, recording.LogOutput)

	// TODO: Remove finished recording from recordings slice using mutex

	return nil
}

// Map to store channels for each recording
var RecordingChannels = make(map[string]chan struct{})

func CancelRecording(recordID string) error {
	log.Printf("[INFO] Cancellation signal sent for camera %s\n", recordID)
	if ch, ok := RecordingChannels[recordID]; ok {
		close(ch)
		delete(RecordingChannels, recordID)
		return nil
	}
	return fmt.Errorf("[ERROR] No active recording found for camera %s\n", recordID)
}
