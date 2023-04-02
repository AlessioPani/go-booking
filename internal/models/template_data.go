package models

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{} // Any data => void interface interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
}
