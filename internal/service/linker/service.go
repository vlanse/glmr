package linker

import "github.com/samber/lo"

type Service struct {
	settings Settings
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetProjectLinks(projectID int64) []Link {
	project, found := lo.Find(s.settings.Projects, func(item Project) bool {
		return item.ID == projectID
	})
	if !found {
		return nil
	}

	res := make([]Link, 0, len(s.settings.Templates))
	for _, t := range s.settings.Templates {
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
	s.settings = settings
}
