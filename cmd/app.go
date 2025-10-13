package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/samber/lo"
	mrV1 "github.com/vlanse/glmr/internal/api/mr/v1"
	"github.com/vlanse/glmr/internal/service/gitlab"
	"github.com/vlanse/glmr/internal/service/mr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	grpcServer *grpc.Server
	mux        *runtime.ServeMux

	cfg Config

	gitlabSvc *gitlab.Service
	mrSvc     *mr.Service
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
	pathPriority := []func() (string, error){
		func() (string, error) {
			return os.Getwd()
		},
		func() (string, error) {
			ex, err := os.Executable()
			if err != nil {
				return "", fmt.Errorf("не удалось получить путь к исполняемому файлу: %w", err)
			}
			return filepath.Dir(ex), nil
		},
		func() (string, error) {
			curUser, err := user.Current()
			if err != nil {
				return "", fmt.Errorf("ошибка получения контекста текущего пользователя ОС: %w", err)
			}
			return curUser.HomeDir, nil
		},
	}

	var err error
	for _, pathGetter := range pathPriority {
		var path string
		if path, err = pathGetter(); err == nil {
			configPath := filepath.Join(path, configFilename)
			if a.cfg, err = loadConfig(configPath); err != nil {
				return fmt.Errorf("ошибка открытия файла настроек: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("ошибка открытия файла конфигурации: %w", err)
}

func (a *App) initServices(_ context.Context) error {
	a.gitlabSvc = gitlab.NewService(a.cfg.Gitlab.URL, a.cfg.Gitlab.Token)

	a.mrSvc = mr.NewService(
		mr.Settings{
			Groups: lo.Map(a.cfg.Groups, func(item Group, _ int) mr.ProjectGroupSettings {
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
		},
		a.gitlabSvc,
	)

	return nil
}

func (a *App) initAPI(ctx context.Context) error {
	a.grpcServer = grpc.NewServer()
	a.mux = runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := mrV1.New(a.mrSvc).Register(ctx, a.grpcServer, a.mux, grpcServerEndpoint, opts); err != nil {
		return fmt.Errorf("инициализация mr v1 API: %w", err)
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
