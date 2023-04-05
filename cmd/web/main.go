package main

import (
	"encoding/gob"
	"fmt"
	"github.com/AlessioPani/go-booking/internal/config"
	"github.com/AlessioPani/go-booking/internal/handlers"
	"github.com/AlessioPani/go-booking/internal/models"
	"github.com/AlessioPani/go-booking/internal/renders"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

// Port number of our website
const portNumber = ":8080"

// app is the config struct of our webapp
var app config.AppConfig

// Session
var session *scs.SessionManager

// main is the entry point.
func main() {

	err := run()
	if err != nil {
		log.Fatal(err)
	}

	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	fmt.Println("Starting an application on port", portNumber[1:])
	err = serve.ListenAndServe()
	log.Fatal(err)
}

func run() error {
	// data models I'm going to put to the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	tc, err := renders.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}
	app.TemplateCache = tc
	app.UseCache = false
	app.Session = session

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	renders.NewTemplates(&app)

	return nil
}
