package main

import (
	"net/http"
	"os"
	"testing"
)

/*
Aquí declaramos los elementos que precisaríamos antes de la prueba de TestMain,
*/
func TestMain(m *testing.M) {

	// Comando usado para referenciar un método a correr antes de que muera el proceso de prueba.
	os.Exit(m.Run())
}

type myHandler struct{}

func (m *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
