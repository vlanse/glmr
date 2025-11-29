package editor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

const (
	projectPathPlaceholder = "{project_path}"
)

type Service struct {
	cmd      string
	projects map[int64]Project
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) UpdateSettings(settings Settings) {
	s.cmd = settings.Cmd
	s.projects = lo.SliceToMap(settings.Projects, func(item Project) (int64, Project) {
		return item.ID, item
	})
}

func (s *Service) OpenProject(_ context.Context, id int64) error {
	if len(s.cmd) == 0 {
		return errors.New("editor command line not configured")
	}

	project, ok := s.projects[id]
	if !ok {
		return fmt.Errorf("path to project with ID %d not specified", id)
	}

	parts := strings.Split(s.cmd, " ")
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
	if len(s.cmd) == 0 {
		return false
	}

	_, ok := s.projects[id]
	return ok
}
