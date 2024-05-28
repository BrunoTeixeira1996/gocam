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
   - validate if its possible to login with current user:password for rtsp
   - add slog instead of log and write log to file
   - [ERROR] Process for j4wxmTPAY4 ID did not exit gracefully, force killing...
     - When killing a process I think its good to still have the ffmpeg output
   - Write a README
   - test if recording from different cameraId works(this can only be tested when I have 2 cameras)
*/
