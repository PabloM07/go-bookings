package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pablom07/go-course/internal/config"
	"github.com/pablom07/go-course/internal/handlers"
)

func routes(_ *config.AppConfig) http.Handler {
	/*	Ruteo usando github.com/bmizerany/pat
		mux := pat.New()
		mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
		mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
	*/

	// Ruteo usando github.com/go-chi/chi/v5
	mux := chi.NewRouter()
	// Seteo funciones middleware HTTP
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/search-availability", handlers.Repo.SearchAvailability)
	mux.Get("/make-reservation", handlers.Repo.MakeReservation)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)
	mux.Get("/reservation-summary", handlers.Repo.ReservSum)

	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/ping", handlers.Repo.Ping)

	// Declaración que sirve para cargar todos los recursos estáticos del folder static
	mux.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	return mux
}
