package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// Función de prueba de un middleware HTTP
func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit the page")
		next.ServeHTTP(w, r)
	})
}

// Agrega protección CSRF a todos los requests POST
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// Carga y persiste la sesión en cada request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
