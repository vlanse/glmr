package linker

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
)

type Link struct {
	DisplayName string
	Link        string
}

type Project struct {
	ID    int64
	Attrs map[string]any
}

type Settings struct {
	Templates []Link
	Projects  []Project
}

func formatLink(template string, projectID int64, attrs map[string]any) string {
	for k, v := range attrs {
		template = strings.Replace(template, fmt.Sprintf("{%s}", k), cast.ToString(v), -1)
	}

	template = strings.Replace(template, "{id}", cast.ToString(projectID), -1)
	if strings.Contains(template, "{") && strings.Contains(template, "}") {
		return ""
	}
	return template
}
