package models

import "github.com/pablom07/go-course/internal/forms"

// TemplateData alberga los datos enviados desde los controladores hacia los templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]any
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}
