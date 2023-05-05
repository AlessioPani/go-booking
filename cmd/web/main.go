package main

import (
	"encoding/gob"
	"flag"
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

	// parse flags provided by cli
	production := flag.Bool("production", true, "True for production, false for development")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database hostname")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbName := flag.String("dbname", "bookings", "Database name")
	dbUser := flag.String("dbuser", "postgres", "Database username")
	dbPassword := flag.String("dbpassword", "", "Database password")
	flag.Parse()

	// create a channel
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// starting mail listener
	log.Println("Starting mail listener...")
	listenForMail()

	// change this to true when in production
	app.InProduction = *production

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
	app.UseCache = *useCache
	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL(fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPassword))
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
