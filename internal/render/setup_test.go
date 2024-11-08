package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/pablom07/go-course/internal/config"
	"github.com/pablom07/go-course/internal/models"
)

// Declaramos los elementos web que necesitamos para ejecutar las pruebas
var session *scs.SessionManager
var testApp config.AppConfig

// Declaramos una función que inicia todos los elementos necesarios para las pruebas
func TestMain(m *testing.M) {

	gob.Register(models.Reservation{})

	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction

	testApp.Session = session

	app = &testApp

	//Terminamos el proceso con el código que nos devolvió la ejecución de las pruebas
	os.Exit(m.Run())
}