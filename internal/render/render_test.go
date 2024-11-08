package render

import (
	"net/http"
	"testing"

	"github.com/pablom07/go-course/internal/models"
)

// Función que prueba el método render.AddDefaultData()
func TestAddDefaultData(t *testing.T) {
	
	// Se obtienen los modelos de templates renderizados
	var td models.TemplateData

	// Obtenemos la sesión que formamos en el método getSession()
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")


	// Almacenamos en la variable result la devolución del método a probar y controlamos
	// que no devuelva error
	res := AddDefaultData(&td, r) 
	if res.Flash != "123" {
		t.Error("Failed ")
	}
}

// Función que prueba y genera los parámetros que debe contener el contexto de ejecución http
func getSession() (*http.Request, error) {
	// Creamos el request para probar
	r, err := http.NewRequest("GET", "/some-url", nil);
	if err != nil {
		return nil, err
	}
	
	// Obtenemos el contexto del request creado
	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}