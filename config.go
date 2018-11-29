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
	KeepAllFiles      bool
	Path              string
	FileRegexPatterns []string
}

func NewConfig() *Config {
	return &Config{OutputDir: "archive"}
}

func CheckConfig(conf *Config) error {
	if len(conf.InputDirs) == 0 {
		return errors.New("InputDirs must have at least one entry.")
	}

	for i, v := range conf.InputDirs {
		if v.Path == "" || len(v.FileRegexPatterns) == 0 {
			return fmt.Errorf("InputDirs[%d] must have both path and files set", i)
		}
	}

	return nil
}
