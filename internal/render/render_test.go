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

// Test unitario que prueba la función RenderTemplate
func TestRenderTemplate(t *testing.T) {

	// Debemos sobreescribir el valor de la ruta de los templates debido a que las pruebas
	// corren en un subdirectorio de la ejecución normal.
	tmplPath = "./../../templates"

	// Instanciamos el caché de templates para renderizar en las pruebas requeridas.
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	// Asignamos el caché a la AppConfig que correremos para levantarlo.
	app.TemplateCache = tc

	// Obtenemos una instancia de prueba del tipo request para pasar como parámetro.
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	// Creamos una estructura que implemente la interfaz para simular un ResponseWriter
	// para pasarlo como parámetro al método que queremos probar
	var ww myWriter

	// Corremos el primer test pidiendo que renderice la home page y que no devuelva errores
	err = RenderTemplate(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("Error writing template to browser")
	}

	// Efectuamos una prueba con un nombre de template inexistente para probar que devuelva
	// error de template no encontrado.
	err = RenderTemplate(&ww, r, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("Rendered template that does not exists")
	}

}

// Prueba unitaria que testea la función que setea la config del paquete templates
func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

// Prueba unitaria que testea la implementación de la creación de la caché de templates
func TestCreateTemplateCache(t *testing.T) {

	// Asignamos el path de los templates para la ubicación de la ejecución de la prueba unitaria
	tmplPath = "./../../templates"

	// Corremos la función ignorando el resultado feliz, pero capturando errores en caso de fallos
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

// Función que prueba y genera los parámetros que debe contener el contexto de ejecución http
func getSession() (*http.Request, error) {
	// Creamos el request para probar
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	// Obtenemos el contexto del request creado
	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}
