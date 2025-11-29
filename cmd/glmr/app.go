package main

import (
	"context"
	"fmt"
	"log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/samber/lo"
	editorV1 "github.com/vlanse/glmr/internal/api/editor/v1"
	mrV1 "github.com/vlanse/glmr/internal/api/mr/v1"
	versionV1 "github.com/vlanse/glmr/internal/api/version/v1"
	"github.com/vlanse/glmr/internal/service/editor"
	"github.com/vlanse/glmr/internal/service/gitlab"
	"github.com/vlanse/glmr/internal/service/mr"
	"github.com/vlanse/glmr/internal/util/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	grpcServer *grpc.Server
	mux        *runtime.ServeMux

	cfgProvider *config.Provider[Config]

	gitlabSvc *gitlab.Service
	mrSvc     *mr.Service
	editorSvc *editor.Service
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run(ctx context.Context) error {
	initializers := []func(ctx context.Context) error{
		a.initConfig,
		a.initServices,
		a.initAPI,
		a.startBackgroundWorkers,
	}

	for _, initializer := range initializers {
		if err := initializer(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	var err error

	if a.cfgProvider, err = config.MakeProvider[Config](configFilename); err != nil {
		return err
	}
	a.cfgProvider.ChangeCallback = a.updateConfig

	return nil
}

func (a *App) initServices(_ context.Context) error {
	cfg := a.cfgProvider.GetConfig()

	a.gitlabSvc = gitlab.NewService(cfg.Gitlab.URL, cfg.Gitlab.Token)

	a.mrSvc = mr.NewService(a.gitlabSvc)

	a.editorSvc = editor.NewService()

	a.updateConfig(cfg)

	return nil
}

func (a *App) initAPI(ctx context.Context) error {
	a.grpcServer = grpc.NewServer()
	a.mux = runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := mrV1.New(a.mrSvc, a.editorSvc).Register(ctx, a.grpcServer, a.mux, grpcServerEndpoint, opts); err != nil {
		return fmt.Errorf("init mr v1 API: %w", err)
	}

	if err := versionV1.New().Register(ctx, a.grpcServer, a.mux, grpcServerEndpoint, opts); err != nil {
		return fmt.Errorf("init version v1 API: %w", err)
	}

	if err := editorV1.New(a.editorSvc).Register(ctx, a.grpcServer, a.mux, grpcServerEndpoint, opts); err != nil {
		return fmt.Errorf("init editor v1 API: %w", err)
	}

	return nil
}

func (a *App) startBackgroundWorkers(_ context.Context) error {
	fmt.Printf("Web interface available at http://%s\n", httpServerEndpoint)
	if err := runServer(a.grpcServer, a.mux); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (a *App) updateConfig(cfg Config) {
	mrSettings := mr.Settings{
		JIRA: mr.JIRA{
			URL: cfg.JIRA.URL,
		},
		Groups: lo.Map(cfg.Groups, func(item Group, _ int) mr.ProjectGroupSettings {
			return mr.ProjectGroupSettings{
				Name: item.Name,
				Projects: lo.Map(item.Projects, func(item Project, _ int) mr.ProjectSettings {
					return mr.ProjectSettings{
						Name: item.Name,
						ID:   item.ID,
					}
				}),
			}
		}),
	}
	a.mrSvc.UpdateSettings(mrSettings)

	editorSettings := editor.Settings{
		Cmd: cfg.Editor.Cmd,
		Projects: func() []editor.Project {
			var res []editor.Project
			for _, g := range cfg.Groups {
				for _, p := range g.Projects {
					if len(p.Path) > 0 {
						res = append(res, editor.Project{
							ID:   p.ID,
							Path: p.Path,
						})
					}
				}
			}
			return res
		}(),
	}
	a.editorSvc.UpdateSettings(editorSettings)

	a.gitlabSvc.UpdateSettings(cfg.Gitlab.URL, cfg.Gitlab.Token)
}
