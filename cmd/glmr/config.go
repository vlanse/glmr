package main

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
