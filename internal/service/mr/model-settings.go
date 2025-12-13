package mr

import (
	"github.com/samber/lo"
	"github.com/vlanse/glmr/internal/service/gitlab"
)

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
	Groups      []ProjectGroupSettings
	JIRA        JIRA
	ShowStarred bool
}

func (s Settings) GetProjects(starred []gitlab.Project) []Project {
	var projects []Project
	var userAddedIDs []int64
	for _, group := range s.Groups {
		for _, project := range group.Projects {
			projects = append(projects, Project{
				ID:        project.ID,
				Name:      project.Name,
				GroupName: group.Name,
			})
			userAddedIDs = append(userAddedIDs, project.ID)
		}
	}

	for _, p := range starred {
		if lo.Contains(userAddedIDs, p.ID) {
			continue
		}
		projects = append(projects, Project{
			ID:        p.ID,
			Name:      p.Name,
			GroupName: "starred projects",
		})
	}
	return projects
}
