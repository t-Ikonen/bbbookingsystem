package models

//TemplateData holds data send to templates
type TemplateData struct {
	StringMap map[string]string
	FloatMap  map[string]float32
	IntMap    map[string]int
	Data      map[string]interface{}
	CSFRToken string
	Flash     string
	Error     string
	Warning   string
}
