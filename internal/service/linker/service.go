package linker

import (
	"sync"

	"github.com/samber/lo"
)

type Service struct {
	mx       sync.RWMutex
	settings Settings
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetProjectLinks(projectID int64) []Link {
	settings := s.getSettings()

	project, found := lo.Find(settings.Projects, func(item Project) bool {
		return item.ID == projectID
	})
	if !found {
		return nil
	}

	res := make([]Link, 0, len(settings.Templates))
	for _, t := range settings.Templates {
		formatted := formatLink(t.Link, projectID, project.Attrs)
		if len(formatted) == 0 {
			continue
		}
		res = append(res, Link{
			DisplayName: t.DisplayName,
			Link:        formatted,
		})
	}
	return res
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
