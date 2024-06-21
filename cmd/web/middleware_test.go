package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler

	h := NoSurf(&myH)

	// Corte de control para determinar que implementación lleva una estructura.
	switch v := h.(type) {
	case http.Handler:
		// No hacer nada, dado que es lo que esperamos
	default:
		// EL placeholder %T es un marcador para inserción de valores de tipos de datos.
		t.Errorf("type is not http.Handler, but is %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler

	h := SessionLoad(&myH)

	// Corte de control para determinar que implementación lleva una estructura.
	switch v := h.(type) {
	case http.Handler:
		// No hacer nada, dado que es lo que esperamos
	default:
		// EL placeholder %T es un marcador para inserción de valores de tipos de datos.
		t.Errorf("type is not http.Handler, but is %T", v)
	}

}
