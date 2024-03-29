package renders

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/AlessioPani/go-booking/internal/config"
	"github.com/AlessioPani/go-booking/internal/models"

	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{
	"humanDate":  HumanDate,
	"formatDate": FormatDate,
	"iterate":    Iterate,
}
var app *config.AppConfig
var pathToTemplates = "./templates"

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// HumanDate converts a time.Time into a string with a YYYY-MM-DD format
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDate converts a time.Time into a string with a f format string
func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

// Iterate returns a slice of ints, starting at 1 and going to count
func Iterate(count int) []int {
	var i int
	var items []int

	for i = 1; i <= count; i++ {
		items = append(items, i)
	}

	return items
}

// AddDefaultData adds default data to templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	td.IsAuthenticated = app.Session.Exists(r.Context(), "user_id")
	return td
}

// Template renders a template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the AppConfig
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		return errors.New("can't get the template from cache")
	}

	td = AddDefaultData(td, r)

	err := t.Execute(w, td)
	if err != nil {
		log.Print(err)
		return errors.New("failed to execute template")
	}

	return nil
}

// CreateTemplateCache populates the template cache
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// get all the files named *.layout.tmpl from ./templates
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
