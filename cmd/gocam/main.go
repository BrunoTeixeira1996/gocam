package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/BrunoTeixeira1996/gocam/internal/config"
	"github.com/BrunoTeixeira1996/gocam/internal/handles"
)

func run() error {
	var configFlag = flag.String("file", "", "use this to provide the config file full path")
	flag.Parse()

	if *configFlag == "" {
		return fmt.Errorf("[ERROR] Please provide a valid config file")
	}

	cfg, err := config.ReadTomlFile(*configFlag)
	if err != nil {
		return err
	}

	if err := handles.Init(cfg.Targets, cfg.Conf.ListenPort, cfg.Conf.DumpRecording, cfg.Conf.LogRecording); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

/*
TODO:
   - add slog instead of log and write log to file
   - create a way to list current recordings
     - Instead of having RecordingChannels, DumpLocation and LogLocation in UI struct I can create another struct that is filled by StartFFMPEG function with data from the current recording
     - This way, UI struct will have a slice of those structs will all the useful information that can be viewed from the web or queried from CLI
   - [ERROR] Process for j4wxmTPAY4 ID did not exit gracefully, force killing...
     - When killing a process I think its good to still have the ffmpeg output
   - check the FIXMEE
   - clean code
*/
