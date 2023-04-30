package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AlessioPani/go-booking/internal/config"
	"github.com/AlessioPani/go-booking/internal/driver"
	"github.com/AlessioPani/go-booking/internal/handlers"
	"github.com/AlessioPani/go-booking/internal/helpers"
	"github.com/AlessioPani/go-booking/internal/models"
	"github.com/AlessioPani/go-booking/internal/renders"

	"github.com/alexedwards/scs/v2"
)

// Port number of our website
const portNumber = ":8080"

// app is the config struct of our webapp
var app config.AppConfig

// Session
var session *scs.SessionManager

var infoLog *log.Logger
var errorLog *log.Logger

// main is the entry point.
func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)
	listenForMail()

	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	fmt.Println("Starting an application on port", portNumber[1:])
	err = serve.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// data models I'm going to put to the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// create a channel
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// starting mail listener
	log.Println("Starting mail listener...")
	listenForMail()

	// change this to true when in production
	app.InProduction = false

	// Set up the infoLog
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	// Set up the errorLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	tc, err := renders.CreateTemplateCache()
	if err != nil {
		log.Println(err)
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false
	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=postgres")
	if err != nil {
		log.Fatal("cannot connect to database")
	}
	log.Println("Connected to database")

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	renders.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
