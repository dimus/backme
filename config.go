package backme

import (
	"errors"
	"fmt"
)

type Config struct {
	Logfile   string
	OutputDir string
	InputDirs []InputDir
}

type InputDir struct {
	Path  string
	Files []string
}

func NewConfig() *Config {
	return &Config{OutputDir: "backme"}
}

func CheckConfig(conf *Config) error {
	if len(conf.InputDirs) == 0 {
		return errors.New("InputDirs must have at least one entry.")
	}

	for i, v := range conf.InputDirs {
		if v.Path == "" || len(v.Files) == 0 {
			return fmt.Errorf("InputDirs[%d] must have both path and files set", i)
		}
	}

	return nil
}
