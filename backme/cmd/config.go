package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
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

func newConfig() *Config {
	return &Config{OutputDir: "backme"}
}

func getConfig() *Config {
	conf := newConfig()
	err := viper.Unmarshal(conf)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = checkConfig(conf)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	return conf
}

func checkConfig(conf *Config) error {
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
