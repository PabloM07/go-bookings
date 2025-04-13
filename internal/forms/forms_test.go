package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

/*

En este archivo de pruebas, iremos creando una prueba por cada casuística de métodos que se encuentren en el
fichero "forms". Según la convención de nombres para los tests de varias casuísticas, deben concatenar el
prefijo "Test" seguido del parámetro a probar (en este caso un formulario web), seguido por un guión bajo '_'
y finalizando con el método a probar.

*/

// Se prueba la instancia de un form sin campos requeridos
func TestForm_valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("form invalid when should have valid")
	}
}

// Se prueba la instancia de un form con campos requeridos incompletos, luego a esos campos se le agregan
// valores para que la validación sea correcta.
func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form invalid when should have valid")
	}
}

// Se prueba la función de validación de campos requeridos
func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	if form.Has("a") {
		t.Error("Form passed when required field is missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "")
	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.Form = postedData

	if form.Has("a") {
		t.Error("Form passed when required field is empty")
	}

	postedData = url.Values{}
	postedData.Add("a", "Field has value")
	form = New(postedData)
	if !form.Has("a") {
		t.Error("Form not passed when has required field")
	}
}

// Se prueba la validación de campos con mínimo de caracteres requeridos.
func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)
	postedData.Add("a", "Test passed successfully")

	form.MinLength("a", 30)
	if form.Valid() {
		t.Error("Form has passed with field without enough characters")
	}

	isError := form.Errors.Get("a")
	if isError == "" {
		t.Error("Should have an error, but did not get one")
	}

	postedData = url.Values{}
	postedData.Add("a", "The system should handle errors gracefully")
	form = New(postedData)

	form.MinLength("a", 30)
	if !form.Valid() {
		t.Error("Form has rejected with field with enough characters")
	}

	isError = form.Errors.Get("a")
	if isError != "" {
		t.Error("Should not have an error, but did get one")
	}
}

// Se prueba la validación de patrones de dirección de correo electrónico.
func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("Form shows valid email for non-existent field")
	}

	postedData = url.Values{}
	postedData.Add("email", "me@here.com")
	form = New(postedData)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("Got an invalid email when we should not have")
	}

	postedData = url.Values{}
	postedData.Add("email", "x")
	form = New(postedData)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("Got a valid email when we should have")
	}
}
