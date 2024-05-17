package handles

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/BrunoTeixeira1996/gocam/internal/action"
	"github.com/BrunoTeixeira1996/gocam/internal/config"
)

type Record struct {
	Duration string
}

type UI struct {
	tmpl         *template.Template
	Targets      []config.Target
	DumpLocation string
}

func (ui *UI) indexHandler(w http.ResponseWriter, r *http.Request) {
	if err := ui.tmpl.ExecuteTemplate(w, "index.html.tmpl", ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ui *UI) listHandler(w http.ResponseWriter, r *http.Request) {
	if err := ui.tmpl.ExecuteTemplate(w, "targets.html.tmpl", map[string]interface{}{
		"targets": ui.Targets,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Map to store channels for each recording
var recordingChannels = make(map[string]chan struct{})

// receives POST with info to start recording that POST
// contains a duration object with the duration of the recording
// if not present assume its 2h
// after that we execute action.startFFMPEGRecording() function that records the current camera
func (ui *UI) recordHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != "POST" {
		e := "[ERROR] HTTP Method not allowed"
		http.Error(w, e, http.StatusBadRequest)
		log.Println(e)
		return
	}

	responseBody, err := io.ReadAll(r.Body)
	if err != nil {
		e := "[ERROR] While reading response body: " + err.Error()
		http.Error(w, e, http.StatusBadRequest)
		log.Println(e)
		return
	}

	// Unmarshal JSON data
	var record Record
	if err := json.Unmarshal(responseBody, &record); err != nil {
		e := "[INFO] No JSON object provided in POST, assuming 2h record time"
		http.Error(w, e, http.StatusBadRequest)
		log.Println(e)
		record.Duration = "2h"
	}

	recordingDuration, err := time.ParseDuration(record.Duration)
	if err != nil {
		e := "[ERROR] While parsing duration: " + err.Error()
		http.Error(w, e, http.StatusBadRequest)
		log.Println(e)
		return
	}

	camera1Channel := make(chan struct{})
	recordingChannels["camera1"] = camera1Channel

	// convert hour to seconds since ffmpeg uses seconds as time duration
	recordingDurationS := fmt.Sprintf("%.f", recordingDuration.Seconds())
	// if err = action.StartFFMPEGRecording(recordingDuration, recordingDurationS, camera1Channel); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	log.Println(err)
	// 	return
	// }
	// FIXMEE: I want to grab the error here
	go action.StartFFMPEGRecording(recordingDuration, recordingDurationS, camera1Channel, ui.DumpLocation)
}

func cancelRecording(cameraID string) {
	if ch, ok := recordingChannels[cameraID]; ok {
		log.Println("CancelRecording function, channels: ", recordingChannels)
		close(ch)
		delete(recordingChannels, cameraID)
		log.Printf("Cancellation signal sent for camera %s\n", cameraID)
	} else {
		log.Printf("No active recording found for camera %s\n", cameraID)
	}
}

func (ui *UI) cancelHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.Body
	cancelRecording("camera1")
	fmt.Fprintf(w, "Recording for camera 1 cancelled")
}

//go:embed assets/*
var assetsDir embed.FS

func Init(targets []config.Target, listenPort string, dumpPath string) error {
	var err error

	tmpl, err := template.ParseFS(assetsDir, "assets/*.tmpl")
	if err != nil {
		return err
	}

	ui := &UI{
		tmpl:         tmpl,
		Targets:      targets,
		DumpLocation: dumpPath,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", ui.indexHandler)
	mux.HandleFunc("/listcameras/", ui.listHandler)
	mux.HandleFunc("/record/", ui.recordHandler)
	// Example endpoint to cancel recording for camera 1
	mux.HandleFunc("/cancel/", ui.cancelHandler)

	log.Printf("Listening at :%s\n", listenPort)

	err = http.ListenAndServe(":"+listenPort, mux)
	if err != nil && err != http.ErrServerClosed {
		panic("Error trying to start http server: " + err.Error())
	}

	return nil
}
