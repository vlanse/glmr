package plugin

import (
	"fmt"
	"plugin"
)

const (
	methodName = "GLMRLoadProjectInfoV1"
)

type Input struct {
	ProjectIID int64
	Name       string
	Params     map[string]string
}

type Output struct {
	HTML      string
	PlainText string
}

type Method func(Input) (Output, error)

type Plugin struct {
	Name   string
	Method *Method

	pl *plugin.Plugin
}

func Load(name, path string) (*Plugin, error) {
	p := &Plugin{}
	var err error
	if p.pl, err = plugin.Open(path); err != nil {
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}

	m, err := p.pl.Lookup(methodName)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup plugin method: %w", err)
	}

	var ok bool
	if p.Method, ok = m.(*Method); !ok {
		return nil, fmt.Errorf("invalid plugin method signature")
	}

	if _, err = (*p.Method)(Input{ProjectIID: -1}); err != nil {
		return nil, fmt.Errorf("plugin internal error: %w", err)
	}

	p.Name = name
	return p, nil
}
