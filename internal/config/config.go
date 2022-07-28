package config

import (
	"errors"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var (
	ErrParseFile = errors.New("error during file parsing")
)

type LoggerCfg struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type DatabaseCfg struct {
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
	Login    string `yaml:"login"`
	Password string `env:"DATABASE_PASSWORD"`
}

type HttpCfg struct {
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type Config struct {
	Http        HttpCfg     `yaml:"http"`
	Logger      LoggerCfg   `yaml:"logger"`
	Database    DatabaseCfg `yaml:"database"`
	HyperParams struct {
		ShutdownTimeout int `yaml:"shutdown_timeout"`
	} `yaml:"hyper_params"`
}

func readParamsFromConfigFile(pathToFile string, config *Config) error {
	file, err := os.Open(pathToFile)
	if err != nil {
		return err
	}
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(fileContent, config); err != nil {
		return err
	}
	return nil
}

func readParamsFromEnv(config *Config) error {
	if err := godotenv.Load(); err != nil {
		log.Print("failed env var from .env file")
	}

	if err := envconfig.Process("", config); err != nil {
		return ErrParseFile
	}
	return nil
}

func New(pathToConfigFile string) (*Config, error) {
	var config Config

	if err := readParamsFromConfigFile(pathToConfigFile, &config); err != nil {
		return nil, err
	}
	if err := readParamsFromEnv(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
