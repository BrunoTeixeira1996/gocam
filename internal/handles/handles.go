package handles

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"text/template"

	"github.com/BrunoTeixeira1996/gocam/internal/action"
	"github.com/BrunoTeixeira1996/gocam/internal/config"
)

type Record struct {
	Duration string
}

type UI struct {
	Config            config.Config            // config Object
	tmpl              *template.Template       // pointer to template
	Targets           []config.Target          // list of all cameras
	RecordingChannels map[string]chan struct{} // list of all ongoing records
	Recordings        []action.Recording
	DumpOutput        string // path for the .mp4 dump
	LogOutput         string // path for the .log dump
}

// Removes from ui.Recordings the current canceled/terminated recording
func (ui *UI) RemoveCanceledRecording(Id string) {
	var indexToRemove int
	for i, v := range ui.Recordings {
		if v.Id == Id {
			indexToRemove = i
			break
		}
	}
	ui.Recordings[indexToRemove] = ui.Recordings[len(ui.Recordings)-1]
	ui.Recordings = ui.Recordings[:len(ui.Recordings)-1]
}

// Receives GET and displays current cameras from config file
func (ui *UI) listHandler(w http.ResponseWriter, r *http.Request) {
	if err := ui.tmpl.ExecuteTemplate(w, "targets.html.tmpl", map[string]interface{}{
		"targets": ui.Targets,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Receives GET and displays ongoing recordings
func (ui *UI) listRecordingsHandler(w http.ResponseWriter, r *http.Request) {
	if err := ui.tmpl.ExecuteTemplate(w, "recordings.html.tmpl", map[string]interface{}{
		"recordings": ui.Recordings,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Receives POST with info to start recording that POST
// contains a duration object with the duration of the recording
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
	var t struct {
		CameraId string
		Duration string
	}
	if err := json.Unmarshal(responseBody, &t); err != nil {
		e := "[ERROR] No JSON object provided in POST"
		http.Error(w, e, http.StatusBadRequest)
		log.Println(e)
		return
	}

	// Validate if cameraID exist in the config toml file
	if !ui.Config.IsCameraIdValid(t.CameraId) {
		e := fmt.Sprintf("[ERROR] camera %s is not valid", t.CameraId)
		http.Error(w, e, http.StatusBadRequest)
		log.Println(e)
		return
	}

	var recording action.Recording
	recording.Init(t.CameraId, ui.DumpOutput, ui.LogOutput, ui.Config)

	recording.WantDuration = t.Duration

	if err := recording.ParseDurationToSeconds(); err != nil {
		e := "[ERROR] " + err.Error()
		http.Error(w, e, http.StatusBadRequest)
		log.Println(e)
		return
	}

	action.RecordingChannels[recording.Id] = recording.Channel

	ui.RecordingChannels = action.RecordingChannels

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = action.StartFFMPEGRecording(&recording, &ui.Recordings, ui.Config)
	}()
	wg.Wait()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}
}

// Receives POST with an ID that contains the recording identifier
// after that we execute action.CancelRecording(ID) function that stops a given recording
func (ui *UI) cancelHandler(w http.ResponseWriter, r *http.Request) {
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
	data := struct {
		ID string
	}{}

	if err := json.Unmarshal(responseBody, &data); err != nil {
		e := "[INFO] No JSON object provided in POST"
		http.Error(w, e, http.StatusBadRequest)
		log.Println(e)
	}

	// Cancels specific recording
	if err := action.CancelRecording(data.ID); err != nil {
		http.Error(w, err.Error(), http.StatusAccepted)
		log.Println(err)
	}

	// Removes canceled recordings from UI
	ui.RemoveCanceledRecording(data.ID)
}

//go:embed assets/*
var assetsDir embed.FS

func Init(config config.Config, targets []config.Target, listenPort string, dumpRecording string, logRecording string) error {
	var err error

	tmpl, err := template.ParseFS(assetsDir, "assets/*.tmpl")
	if err != nil {
		return err
	}

	ui := &UI{
		Config:     config,
		tmpl:       tmpl,
		DumpOutput: dumpRecording,
		Targets:    targets,
		LogOutput:  logRecording,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/listcameras", http.StatusMovedPermanently)
	})
	mux.HandleFunc("/listcameras/", ui.listHandler)
	mux.HandleFunc("/listrecordings/", ui.listRecordingsHandler)
	mux.HandleFunc("/record/", ui.recordHandler)
	mux.HandleFunc("/cancel/", ui.cancelHandler)

	log.Printf("Listen at port %s - Dump recording at %s - Log recording at %s\n", listenPort, dumpRecording, logRecording)

	err = http.ListenAndServe(":"+listenPort, mux)
	if err != nil && err != http.ErrServerClosed {
		panic("Error trying to start http server: " + err.Error())
	}

	return nil
}
