package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Conf struct {
	JsonFile      string `toml:"json_file"`
	ListenPort    string `toml:"listen_port"`
	DumpRecording string `toml:"dump_recording"`
	LogRecording  string `toml:"log_recording"`
}

type Target struct {
	CameraId      string `toml:"camera_id"`
	Name          string `toml:"name"`
	Host          string `toml:"host"`
	Port          string `toml:"port"`
	Stream        string `toml:"stream"`
	Protocol      string `toml:"protocol"`
	User          string `toml:"user"`
	Password      string `toml:"password"`
	RecordingPath string `toml:"recording_path"`
}

type Config struct {
	Conf    Conf
	Targets []Target
}

// Reads tom file and parses info to a struct
func ReadTomlFile(fileLocation string) (Config, error) {
	var cfg Config

	input, err := os.ReadFile(fileLocation)
	if err != nil {
		return Config{}, err
	}

	if _, err := toml.Decode(string(input), &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Validates if dump path and log path exist
func (c *Config) ValidatePaths() error {
	_, err := os.Stat(c.Conf.DumpRecording)
	if os.IsNotExist(err) {
		return fmt.Errorf("[ERROR] Dump path does not exist")
	}

	_, err = os.Stat(c.Conf.LogRecording)
	if os.IsNotExist(err) {
		return fmt.Errorf("[ERROR] Log path does not exist")
	}

	return nil
}

// Verifies if json file exist
// if not creates an empty json file
func (c *Config) InitJSON() error {
	_, err := os.Stat(c.Conf.JsonFile)

	// file does not exist so lets create
	if os.IsNotExist(err) {
		log.Println("[INFO] Json file does not exist ... going to create")

		file, err := os.Create(c.Conf.JsonFile)
		if err != nil {
			return fmt.Errorf("[ERROR] while creating the JSON file: %s\n", err)
		}

		defer file.Close()

	} else {
		// validate if its a valid json file
		data, err := os.ReadFile(c.Conf.JsonFile)
		if err != nil {
			return fmt.Errorf("[ERROR] while reading the JSON file: %s\n", err)
		}

		// file is empty
		if len(data) == 0 {
			return nil
		}

		// file contains something
		var j interface{}
		if err := json.Unmarshal(data, &j); err != nil {
			return fmt.Errorf("[ERROR] not a valid JSON file: %s\n", err)

		}
	}
	return nil
}

// Validates if given camera Id exists in toml file
func (c *Config) IsCameraIdValid(cameraId string) bool {
	for _, v := range c.Targets {
		if v.CameraId == cameraId {
			return true
		}
	}
	return false
}
