package editor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"unicode"

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

	parts := splitCmdLine(settings.Cmd)
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

	// wait editor process without blocking the request handler:
	// GUI editors stay alive for as long as the IDE is open
	go func() {
		if err := c.Wait(); err != nil {
			log.Printf("editor %q exited with error: %v", cmd, err)
		}
	}()

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

func splitCmdLine(line string) []string {
	var (
		parts  []string
		cur    strings.Builder
		quote  rune
		quoted bool
	)

	flush := func() {
		if cur.Len() > 0 || quoted {
			parts = append(parts, cur.String())
			cur.Reset()
			quoted = false
		}
	}

	for _, r := range line {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				cur.WriteRune(r)
			}
		case r == '"' || r == '\'':
			quote = r
			quoted = true
		case unicode.IsSpace(r):
			flush()
		default:
			cur.WriteRune(r)
		}
	}
	flush()

	return parts
}
