package main

const (
	configFilename = "glmr-config.yaml"
)

type Attrs map[string]any

type Project struct {
	Attrs Attrs  `yaml:",inline"`
	ID    int64  `yaml:"id"`
	Name  string `yaml:"name"`
	Path  string `yaml:"path"`
}

type Group struct {
	Name     string    `yaml:"name"`
	Projects []Project `yaml:"projects"`
}

type ProjectLink struct {
	DisplayName string `yaml:"display_name"`
	Template    string `yaml:"template"`
}

type Plugin struct {
	Name    string
	Path    string
	Enabled bool
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

	ShowStarred bool `yaml:"show_starred"`

	Groups []Group `yaml:"groups"`

	ProjectLinks []ProjectLink `yaml:"project_links"`

	Plugins []Plugin `yaml:"plugins"`
}
