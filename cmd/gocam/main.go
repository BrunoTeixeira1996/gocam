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
