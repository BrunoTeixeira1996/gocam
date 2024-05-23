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

	if err := cfg.ValidatePaths(); err != nil {
		return err
	}

	if err := cfg.InitJSON(); err != nil {
		return err
	}

	if err := handles.Init(cfg, cfg.Targets, cfg.Conf.ListenPort, cfg.Conf.DumpRecording, cfg.Conf.LogRecording); err != nil {
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
   - add recordings (finished and canceled) to json in order to view that in http (ONGOING)
     - write into recording.Cmd
     - fix recording.Log
     - validate if both are correct on the respective folders and in the web server
   - add slog instead of log and write log to file
   - [ERROR] Process for j4wxmTPAY4 ID did not exit gracefully, force killing...
     - When killing a process I think its good to still have the ffmpeg output
   - test if recording from different cameraId works(this can only be tested when I have 2 cameras)
*/
