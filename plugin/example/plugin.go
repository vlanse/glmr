package main

import (
	"context"
	"fmt"

	"github.com/vlanse/glmr/plugin"
)

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

var GLMRPluginInvokePerProjectV1 plugin.MethodPerProject = pluginMain

func main() {
	_, _ = GLMRPluginInvokePerProjectV1(context.Background(), plugin.InputPerProject{})
}
