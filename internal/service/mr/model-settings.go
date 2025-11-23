package mr

import "github.com/samber/lo"

type ProjectSettings struct {
	Name string
	ID   int64
}

type ProjectGroupSettings struct {
	Name     string
	Projects []ProjectSettings
}

func (g ProjectGroupSettings) GetAllProjectIDs() []int64 {
	return lo.Map(g.Projects, func(item ProjectSettings, _ int) int64 {
		return item.ID
	})
}

func (g ProjectGroupSettings) ProjectByID(id int64) (ProjectSettings, bool) {
	return lo.Find(g.Projects, func(item ProjectSettings) bool {
		return item.ID == id
	})
}

type JIRA struct {
	URL string
}

type Settings struct {
	Groups []ProjectGroupSettings
	JIRA   JIRA
}

func (s *Settings) GetProjects() []Project {
	var projects []Project
	for _, group := range s.Groups {
		for _, project := range group.Projects {
			projects = append(projects, Project{
				ID:        project.ID,
				Name:      project.Name,
				GroupName: group.Name,
			})
		}
	}
	return projects
}
