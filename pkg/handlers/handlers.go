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

	renders.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler.
func (pr *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "This come from template data."

	remoteIP := pr.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// send data to template
	renders.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders the make a reservation page and displays form
func (pr *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	renders.RenderTemplate(w, "make-reservation.page.tmpl", &models.TemplateData{})
}

// Generals renders the room page
func (pr *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	renders.RenderTemplate(w, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (pr *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	renders.RenderTemplate(w, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (pr *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	renders.RenderTemplate(w, "search-availability.page.tmpl", &models.TemplateData{})
}

// Contact renders the contact page
func (pr *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	renders.RenderTemplate(w, "contact.page.tmpl", &models.TemplateData{})
}
