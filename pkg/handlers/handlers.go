package handlers

import (
	"net/http"

	"github.com/AlessioPani/go-booking/pkg/config"
	"github.com/AlessioPani/go-booking/pkg/models"
	"github.com/AlessioPani/go-booking/pkg/renders"
)

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// Repo is the repository used by handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the homepage handler.
func (pr *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	pr.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	renders.RenderTemplate(w, "home.page.gotmpl", &models.TemplateData{})
}

// About is the about page handler.
func (pr *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "This come from template data."

	remoteIP := pr.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// send data to template
	renders.RenderTemplate(w, "about.page.gotmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
