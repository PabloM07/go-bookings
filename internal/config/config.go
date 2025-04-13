package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// AppConfig mantiene la configuración de la aplicación
type AppConfig struct {
	// Variable que podemos usar para activar y desactivar el uso de cache
	UseCache bool
	// Mapa de caché de templates
	TemplateCache map[string]*template.Template
	// Logueador de procesos de la Webapp
	InfoLog *log.Logger
	// Logueador de errores de la Webapp
	ErrorLog *log.Logger
	// Flag que determina si el programa corre en producción
	InProduction bool
	// Contiene los datos de la sesión HTTP
	Session *scs.SessionManager
}
