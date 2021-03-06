package models

import (
	"github.com/t-Ikonen/bbbookingsystem/internal/forms"
)

//TemplateData holds data send to templates
type TemplateData struct {
	StringMap       map[string]string
	FloatMap        map[string]float32
	IntMap          map[string]int
	Data            map[string]interface{}
	CSRFToken       string
	Flash           string
	Error           string
	Warning         string
	Form            *forms.Form
	IsAuthenticated int
}
