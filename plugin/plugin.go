package plugin

import (
	"context"
	"fmt"
	"plugin"
)

const (
	methodNamePerProject = "GLMRPluginInvokePerProjectV1"
)

type InputPerProject struct {
	ProjectName  string
	ProjectID    int64
	ProjectAttrs map[string]any
}

type OutputPerProject struct {
	HTML      string
	PlainText string
}

type MethodPerProject func(context.Context, InputPerProject) (OutputPerProject, error)

type Plugin struct {
	Name string

	methodPerProject *MethodPerProject
	pl               *plugin.Plugin
}

func (p *Plugin) InvokePerProject(ctx context.Context, in InputPerProject) (OutputPerProject, error) {
	return (*p.methodPerProject)(ctx, in)
}

func Load(name, path string) (*Plugin, error) {
	p := &Plugin{}
	var err error
	if p.pl, err = plugin.Open(path); err != nil {
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}

	m, err := p.pl.Lookup(methodNamePerProject)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup plugin method: %w", err)
	}

	var ok bool
	if p.methodPerProject, ok = m.(*MethodPerProject); !ok {
		return nil, fmt.Errorf("invalid plugin method signature")
	}

	if _, err = (*p.methodPerProject)(context.Background(), InputPerProject{ProjectID: -1}); err != nil {
		return nil, fmt.Errorf("plugin internal error: %w", err)
	}

	p.Name = name
	return p, nil
}
