package models

type Template struct {
	Name    string `json:"templatename,omitempty"`
	Content string `json:"content,omitempty"`
}

type TemplateUpdate struct {
	Content string `json:"content,omitempty"`
}

func (temp *Template) ToMap() (map[string]interface{}, error) {
	tmpMap := make(map[string]interface{})
	A(tmpMap, "templatename", temp.Name)
	A(tmpMap, "content", temp.Content)
	return tmpMap, nil
}

func (temp *TemplateUpdate) ToMap() (map[string]interface{}, error) {
	tmpMap := make(map[string]interface{})
	A(tmpMap, "content", temp.Content)
	return tmpMap, nil
}
