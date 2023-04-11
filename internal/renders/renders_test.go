package renders

import (
	"github.com/AlessioPani/go-booking/internal/models"
	"net/http"
	"testing"
)

func TestNewRenderer(t *testing.T) {
	NewRenderer(app)
}

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error("Failed to get a Session")
	}

	session.Put(r.Context(), "flash", "123")
	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Error("Flash value of 123 not found in session")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error("Failed to get session")
	}

	w := myWriter{}

	err = Template(&w, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("Error writing a template to browser")
	}

	err = Template(&w, r, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("Rendered a template that does not exist")
	}
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))

	r = r.WithContext(ctx)

	return r, nil
}
