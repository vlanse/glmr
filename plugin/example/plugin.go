package main

import (
	"context"
	"fmt"

	"github.com/vlanse/glmr/plugin"
)

func pluginInit(_ context.Context) error {
	return nil
}

func pluginMain(_ context.Context, input plugin.InputPerProject) (plugin.OutputPerProject, error) {
	if input.ProjectID == -1 {
		return plugin.OutputPerProject{
			PlainText: "self-test",
		}, nil
	}

	return plugin.OutputPerProject{
		PlainText: "plugin result",
		HTML:      fmt.Sprintf("<p>%+v</p>", input),
	}, nil
}

var GLMRPluginInvokePerProjectV1 plugin.PerProjectFuncV1 = pluginMain

var GLMRPluginInitV1 plugin.InitFuncV1 = pluginInit

func main() {
	ctx := context.Background()
	_ = GLMRPluginInitV1(ctx)
	_, _ = GLMRPluginInvokePerProjectV1(ctx, plugin.InputPerProject{})
}
