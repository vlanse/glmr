package editor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/samber/lo"
)

const (
	projectPathPlaceholder = "{project_path}"
)

type Service struct {
	mx       sync.RWMutex
	settings Settings
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) UpdateSettings(settings Settings) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.settings = settings
}

func (s *Service) getSettings() Settings {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.settings
}

func (s *Service) OpenProject(_ context.Context, id int64) error {
	settings := s.getSettings()

	if len(settings.Cmd) == 0 {
		return errors.New("editor command line not configured")
	}

	project, ok := lo.Find(settings.Projects, func(item Project) bool {
		return item.ID == id
	})
	if !ok {
		return fmt.Errorf("path to project with ID %d not specified", id)
	}

	parts := strings.Split(settings.Cmd, " ")
	var cmd string
	var args []string
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if i == 0 {
			cmd = part
			continue
		}
		if part == projectPathPlaceholder {
			path := project.Path
			if strings.HasPrefix(path, "~") {
				hd, _ := os.UserHomeDir()
				path = strings.Replace(path, "~", hd, 1)
			}

			part, _ = filepath.Abs(path)
		}
		args = append(args, part)
	}

	c := exec.Command(cmd, args...)
	if err := c.Start(); err != nil {
		return fmt.Errorf("start editor on project: %w", err)
	}

	if err := c.Wait(); err != nil {
		return fmt.Errorf("start editor on project: %w", err)
	}

	return nil
}

func (s *Service) IsProjectConfigured(id int64) bool {
	settings := s.getSettings()

	if len(settings.Cmd) == 0 {
		return false
	}

	_, ok := lo.Find(settings.Projects, func(item Project) bool {
		return item.ID == id
	})
	return ok
}
