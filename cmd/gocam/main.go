package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/BrunoTeixeira1996/gocam/internal/config"
	"github.com/BrunoTeixeira1996/gocam/internal/handles"
)

func run() error {
	var configFlag = flag.String("f", "", "use this to provide the config file full path")
	var listenPortFlag = flag.String("l", "", "use this to provide the listening port")
	flag.Parse()

	if *configFlag == "" || *listenPortFlag == "" {
		return fmt.Errorf("[ERROR] Please provide valid flags")
	}

	cfg, err := config.ReadTomlFile(*configFlag)
	if err != nil {
		return err
	}

	if err := handles.Init(cfg.Targets, *listenPortFlag); err != nil {
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
   - create dump flag to dump the recording
   - create way of stoping record without breaking the .mp4
     - its working but now i need
       - clean useless prints
       - add ID and return that ID on the record response
       - add that ID to the map of channels so I can list and delete specific goroutines
       - cancel handle is a POST that gets an object like {"recordID":123} and that ID is the ID belonging to the goroutine 123
       - check the FIXMEE
   - create a way to list current recordings
   - clean code
*/
