package action

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/BrunoTeixeira1996/gocam/internal/config"
	"github.com/BrunoTeixeira1996/gocam/internal/utils"
)

var (
	mutex             sync.Mutex                       // mutex to perform actions on memory objects
	RecordingChannels = make(map[string]chan struct{}) // Map to store channels for each recording
)

// Struct that hold some stuff from config file
type RecordingConfig struct {
	Name     string
	Host     string
	Port     string
	Stream   string
	Protocol string
	User     string
	Password string
}

// Struct responsable for holding a respective recording
type Recording struct {
	CameraId           string          `json:"CameraId"`  // camera id from POST
	Id                 string          `json:"Id"`        // random string to identify the recording
	Start              string          `json:"-"`         // start date
	StartDate          string          `json:"StartDate"` // start date with proper output
	WantDuration       string          `json:"-"`         // recording duration from POST
	WantDurationParsed time.Duration   `json:"-"`         // recording duration parsed
	WantDurationS      string          `json:"Duration"`  // recording duration in seconds for ffmpeg
	UntilDate          string          `json:"UntilDate"` // value of start date + duration with proper output
	Cmd                string          `json:"Cmd"`       // cmd used for ffmpeg exec
	CancelChannel      chan struct{}   `json:"-"`         // channel used in the respective recording in order to cancel the goroutine
	DumpOutput         string          `json:"Dump"`      // .mp4 dump
	LogOutput          string          `json:"Log"`       // .log dump
	Status             string          `json:"Status"`    // finished / canceled
	Config             RecordingConfig `json:"-"`         // configs from recording
}

// Initializes recording struct
func (r *Recording) Init(cameraId, dumpOutput, logOutput string, config config.Config) {
	r.CameraId = cameraId
	r.Id = utils.GenerateRandomString(10)
	r.DumpOutput = dumpOutput
	r.LogOutput = logOutput
	r.CancelChannel = make(chan struct{})

	// grab the correct target that is being used in the recording
	var targetIndex int
	for i, v := range config.Targets {
		if v.CameraId == cameraId {
			targetIndex = i
			break
		}
	}

	// copies target config to the respective recording
	r.Config.Name = config.Targets[targetIndex].Name
	r.Config.Host = config.Targets[targetIndex].Host
	r.Config.Port = config.Targets[targetIndex].Port
	r.Config.Stream = config.Targets[targetIndex].Stream
	r.Config.Protocol = config.Targets[targetIndex].Protocol
	r.Config.User = config.Targets[targetIndex].User
	r.Config.Password = config.Targets[targetIndex].Password
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

// Writes into the JSON file the newly recording
func (r *Recording) UpdateJSON(status string, config config.Config) error {
	var recordings []Recording

	fC, err := os.ReadFile(config.Conf.JsonFile)
	if err != nil {
		return fmt.Errorf("[ERROR] while reading the json file:%s\n", err)
	}

	// file its not empty so we can unmarshal to read what it contains
	if len(fC) > 0 {
		if err := json.Unmarshal(fC, &recordings); err != nil {
			return fmt.Errorf("[ERROR] while unmarshal the json file:%s\n", err)
		}
	}

	r.Status = status

	recordings = append(recordings, *r)

	fC, err = json.MarshalIndent(recordings, "", "  ")
	if err != nil {
		return fmt.Errorf("[ERROR] while marshal the json file:%s\n", err)
	}

	// write JSON data to the file
	if err := os.WriteFile(config.Conf.JsonFile, fC, 0644); err != nil {
		return fmt.Errorf("[ERROR] while writing the json file:%s\n", err)
	}

	return nil
}

// Calculates end time since start
func (r *Recording) CaculateUntilDate(currentTime time.Time) error {
	parsedDuration, err := time.ParseDuration(r.WantDuration)
	if err != nil {
		return fmt.Errorf("[ERROR] parsing duration: %s\n", err)
	}

	// Add the duration to the start time
	endTime := currentTime.Add(parsedDuration)

	// Format the end time
	r.UntilDate = endTime.Format("January 02, 2006 15:04:05")

	return nil
}

/*
Setups boilerplate stuff for recording
Defines:
  - Start time
  - Until time
  - Dump output
  - Log output
*/
func (r *Recording) Setup() {
	currentTime := time.Now()

	layoutFormat := "2006-01-02_15:04:05"

	// define start time
	r.Start = currentTime.Format(layoutFormat)

	// format the time into the desired output into r.StartDate to show in template
	t, _ := time.Parse(layoutFormat, r.Start)
	r.StartDate = t.Format("January 02, 2006 15:04:05")

	// calculate end date
	r.CaculateUntilDate(currentTime)

	// define output locations
	r.DumpOutput = r.DumpOutput + r.Start + "-" + r.Id + ".mp4"
	r.LogOutput = r.LogOutput + r.Id + ".log"
}

// Validates if rtsp connection is valid
func isRTSPValid(recording *Recording) bool {
	cmd := exec.Command("ffprobe", "rtsp://"+recording.Config.User+":"+recording.Config.Password+"@"+recording.Config.Host+":"+recording.Config.Port+recording.Config.Stream)

	if _, err := cmd.Output(); err != nil {
		return false
	}
	return true
}

// ffmpeg -i rtsp://brun0teixeira:qwerty654321@192.168.30.44:554/stream1 -c:v copy -c:a aac -strict experimental output.mp4
// Starts ffmpeg process
func StartFFMPEGRecording(recording *Recording, recordings *[]Recording, config config.Config, w http.ResponseWriter) error {

	// Setup boilerplate stuff for recording
	recording.Setup()

	// First, check if rtsp connection is valid
	if !isRTSPValid(recording) {
		return fmt.Errorf("[ERROR] FFMPEG RTSP connection is not valid")
	}

	log.Printf("[INFO] Starting record duration %s for %s file with %s ID on cameraID %s\n", recording.WantDurationS, recording.DumpOutput, recording.Id, recording.CameraId)

	cmd := exec.Command("ffmpeg", "-i", recording.Config.Protocol+"://"+recording.Config.User+":"+recording.Config.Password+"@"+recording.Config.Host+":"+recording.Config.Port+recording.Config.Stream, "-c:v", "copy", "-c:a", "aac", "-strict", "experimental", "-t", recording.WantDurationS, recording.DumpOutput)

	recording.Cmd = cmd.String()

	// Capture the output
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("[ERROR] FFMPEG while starting:%s\n", err)
	}

	// Add recording to slice of recordings
	*recordings = append(*recordings, *recording)

	// Send immediate HTTP response to avoid curl waiting
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Recording started for ID: %s", recording.Id)))

	go func() {
		defer func() {
			// Remove finished recording from recordings slice using mutex
			mutex.Lock()
			removeFinishedRecording(recordings, recording.Id)
			mutex.Unlock()
		}()

		select {
		case <-time.After(recording.WantDurationParsed):
			log.Printf("[INFO] Recording for %s ID finished\n", recording.Id)
			if err := recording.UpdateJSON("Completed", config); err != nil {
				log.Printf("[ERROR] %s\n", err)
			}
		case <-recording.CancelChannel:
			if err := recording.UpdateJSON("Canceled", config); err != nil {
				log.Printf("[ERROR] %s\n", err)
			}
			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				log.Printf("[ERROR] Error sending interrupt signal: %s\n", err)
			}
			// Wait for the process to exit
			done := make(chan error, 1)
			go func() {
				done <- cmd.Wait()
			}()
			if err := <-done; err != nil {
				log.Printf("[ERROR] Process for %s ID finished with error: %s\n", recording.Id, err)
			}
			log.Printf("[INFO] Recording for %s ID cancelled\n", recording.Id)
		}

		// Save the output log
		if err := utils.SaveFFMPEGOutput(recording.LogOutput, recording.Id, output.Bytes()); err != nil {
			log.Printf("[ERROR] %s\n", err)
		}
		log.Printf("[INFO] Finished recording for %s ID log file at %s\n", recording.Id, recording.LogOutput)
	}()

	return nil
}

// Removes finished recording from current recording slice
func removeFinishedRecording(r *[]Recording, recordingToRemove string) {
	var indexToRemove int

	for i, v := range *r {
		if v.Id == recordingToRemove {
			indexToRemove = i
			break
		}
	}

	slice := *r
	if indexToRemove < 0 || indexToRemove >= len(slice) {
		return // Index out of bounds, do nothing
	}

	*r = append(slice[:indexToRemove], slice[indexToRemove+1:]...)
}

func CancelRecording(recordID string) error {
	log.Printf("[INFO] Cancellation signal sent for camera %s\n", recordID)
	if ch, ok := RecordingChannels[recordID]; ok {
		close(ch)
		delete(RecordingChannels, recordID)
		return nil
	}
	return fmt.Errorf("[ERROR] No active recording found for camera %s\n", recordID)
}
