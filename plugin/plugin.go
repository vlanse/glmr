package plugin

import (
	"context"
	"fmt"
	"plugin"
)

const (
	funcNameInit       = "GLMRPluginInitV1"
	funcNamePerProject = "GLMRPluginInvokePerProjectV1"
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

type PerProjectFuncV1 func(context.Context, InputPerProject) (OutputPerProject, error)
type InitFuncV1 func(context.Context) error

type Plugin struct {
	Name string

	methodPerProject *PerProjectFuncV1
	pl               *plugin.Plugin
}

func (p *Plugin) InvokePerProject(ctx context.Context, in InputPerProject) (OutputPerProject, error) {
	return (*p.methodPerProject)(ctx, in)
}

func Load(ctx context.Context, name, path string) (*Plugin, error) {
	p := &Plugin{}
	var err error
	if p.pl, err = plugin.Open(path); err != nil {
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}
	mInit, err := p.pl.Lookup(funcNameInit)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup plugin method %s: %w", funcNameInit, err)
	}
	methodInit, ok := mInit.(*InitFuncV1)
	if !ok {
		return nil, fmt.Errorf("invalid plugin %s method signature", funcNameInit)
	}
	if err = (*methodInit)(ctx); err != nil {
		return nil, fmt.Errorf("failed to init plugin %s: %w", name, err)
	}

	mPerProject, err := p.pl.Lookup(funcNamePerProject)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup plugin method: %w", err)
	}
	if p.methodPerProject, ok = mPerProject.(*PerProjectFuncV1); !ok {
		return nil, fmt.Errorf("invalid plugin method signature")
	}

	if _, err = (*p.methodPerProject)(context.Background(), InputPerProject{ProjectID: -1}); err != nil {
		return nil, fmt.Errorf("plugin internal error: %w", err)
	}

	p.Name = name
	return p, nil
}
