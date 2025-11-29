package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	configFilename = "glmr-config.yaml"
)

type Project struct {
	ID   int64  `yaml:"id"`
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type Group struct {
	Name     string    `yaml:"name"`
	Projects []Project `yaml:"projects"`
}

type Config struct {
	Gitlab struct {
		URL   string `yaml:"url"`
		Token string `yaml:"token"`
	} `yaml:"gitlab"`

	JIRA struct {
		URL string `yaml:"url"`
	} `yaml:"jira"`

	Editor struct {
		Cmd string `yaml:"cmd"`
	}

	Groups []Group `yaml:"groups"`
}

func loadConfig(path string) (Config, error) {
	cfgFile, err := os.Open(path)
	cfg := Config{}
	if err != nil {
		return cfg, err
	}
	d := yaml.NewDecoder(cfgFile)

	if cfgErr := d.Decode(&cfg); cfgErr != nil {
		return cfg, cfgErr
	}
	return cfg, nil
}
