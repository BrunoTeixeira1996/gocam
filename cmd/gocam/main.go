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
   - check the FIXMEE
   - create a way to list current recordings
   - clean code
*/
