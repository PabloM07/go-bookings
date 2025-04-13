package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/pablom07/go-course/internal/config"
)

// En los paquetes helpers se agregan métodos para ser usados como utilidades en
// todo el scope del proyecto

var app *config.AppConfig

// Método iniciaizador de la configuración de la app
func NewHelpers(a *config.AppConfig) {
	app = a
}

// Función para procesar errores de cliente
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

// Función para procesar errores internos de la app
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
