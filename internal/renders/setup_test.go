package renders

import (
	"encoding/gob"
	"github.com/AlessioPani/go-booking/internal/config"
	"github.com/AlessioPani/go-booking/internal/models"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	// data models I'm going to put to the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	testApp.InProduction = false
	testApp.UseCache = false

	// Set up the infoLog
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	// Set up the errorLog
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())

}

type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	return http.Header{}
}

func (tw *myWriter) WriteHeader(i int) {

}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
