package main

import (
	"testing"

	"github.com/go-chi/chi"
	"github.com/pablom07/go-course/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// Caso esperado
		break
	default:
		t.Errorf("Type is not *chi.Mux, type is %T", v)
	}
}