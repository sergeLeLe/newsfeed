package main

import (
	"errors"
	"log"
	"newsfeed/internal/application"
	"newsfeed/internal/config"
	"os"
	"path/filepath"
)

func getConfigPath() (string, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		return "", errors.New("the path to the config is not specified")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", errors.New("could not get the full path to the configuration file")
	}

	return absPath, nil
}

func main() {
	pathToConfig, err := getConfigPath()
	if err != nil {
		log.Fatalf("Error %v", err.Error())
	}
	cfg, err := config.New(pathToConfig)
	if err != nil {
		log.Fatalf("failed to load configuration on the path: %v, err: %v", pathToConfig, err.Error())
	}

	application.Run(cfg)
}