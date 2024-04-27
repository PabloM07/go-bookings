package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/pablom07/go-course/internal/config"
	"github.com/pablom07/go-course/internal/handlers"
	"github.com/pablom07/go-course/internal/models"
	"github.com/pablom07/go-course/internal/render"
)

const portNumber string = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	/*
		Declaración de tipos de datos que se añadirán a la sesión. Se usa para el manejo
		de formularios (en forma de structs), y/o el intercambio de información suelta
		(valores primitivos) para su exposición del lado del cliente 
	*/
	gob.Register(models.Reservation{})

	// Cambiar este valor cuando se encuentre en prod
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// Obtener el template caché desde la configuración de la aplicación
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache:", err)
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	// Web Routing using 'net/http'
	// http.HandleFunc("/", repo.Home)
	// http.HandleFunc("/about", repo.About)

	fmt.Println(fmt.Sprintf("Starting app on port %s", portNumber))
	// _ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = nil
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
