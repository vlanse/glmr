package swagger

import (
	"encoding/json"
	"io"
)

type Merger struct {
	Swagger map[string]any
	title   string
}

func NewMerger(title string) *Merger {
	merger := new(Merger)
	merger.Swagger = map[string]any{}
	merger.title = title
	return merger
}

func (m *Merger) AddFile(f io.Reader) error {
	content, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var swaggerMap any

	if err = json.Unmarshal(content, &swaggerMap); err != nil {
		return err
	}

	merge(m.Swagger, swaggerMap.(map[string]any))

	return nil
}

func merge(a, b map[string]any) {
	if a == nil {
		return
	}

	for key, item := range b {
		if i, ok := item.(map[string]any); ok {
			if _, ok := a[key]; ok {
				merge(a[key].(map[string]any), i)
			} else {
				a[key] = i
			}
		} else {
			a[key] = item
		}
	}
}

func (m *Merger) Content() ([]byte, error) {
	var res []byte
	var err error

	info := m.Swagger["info"].(map[string]interface{})
	info["title"] = m.title
	info["version"] = ""

	if res, err = json.Marshal(m.Swagger); err != nil {
		return nil, err
	}
	return res, nil
}
