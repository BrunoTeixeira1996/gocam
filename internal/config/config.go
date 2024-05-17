package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Conf struct {
	ListenPort    string `toml:"listen_port"`
	DumpRecording string `toml:"dump_recording"`
	LogRecording  string `toml:"log_recording"`
}

type Target struct {
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
