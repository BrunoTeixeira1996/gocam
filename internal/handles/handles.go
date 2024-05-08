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
	tmpl    *template.Template
	Targets []config.Target
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

	// convert hour to seconds since ffmpeg uses seconds as time duration
	recordingDurationS := fmt.Sprintf("%.f", recordingDuration.Seconds())
	if err := action.StartFFMPEGRecording(recordingDurationS); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}
}

//go:embed assets/*
var assetsDir embed.FS

func Init(targets []config.Target, listenPort string) error {
	var err error

	tmpl, err := template.ParseFS(assetsDir, "assets/*.tmpl")
	if err != nil {
		return err
	}

	ui := &UI{
		tmpl:    tmpl,
		Targets: targets,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", ui.indexHandler)
	mux.HandleFunc("/listcameras/", ui.listHandler)
	mux.HandleFunc("/record/", ui.recordHandler)

	log.Printf("Listening at :%s\n", listenPort)

	err = http.ListenAndServe(":"+listenPort, mux)
	if err != nil && err != http.ErrServerClosed {
		panic("Error trying to start http server: " + err.Error())
	}

	return nil
}
