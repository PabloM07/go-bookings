package forms

import (
	"net/http"
	"net/url"
)

/*
	En este documento especificaremos todo lo relacionado con la validación del lado de
	cliente de los formularios que nos envían mediante las peticiones en verbo POST, de
	manera de validar que esas peticiones contengan los datos necesarios para llevar a
	cabo una acción con la información provista. Esto incluye la autorización de ejecución
	de acciones y devolución de un mensaje de confirmación como la interrupción de éstos
	y el retorno de un mensaje de error.
*/

// Form representa una estructura personalizada de formulario, incorpora un objeto url.Values
type Form struct {
	url.Values
	Errors errors
}

// New inicializa una estructura Form
func New(data url.Values) *Form {
	return &Form {
		data,

		// Se declara un mapa de slices de strings vacío, de igual manera que se puede
		// declarar como make(map[string][]string)
		errors(map[string][]string{}),
	}
}

// Has revisa si el campo del formulario está en el cuerpo del request y no está vacío.
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}

// Valid controla el formulario en busca de errores de validación. Devuelve True si el formulario
// tiene datos válidos o false si existen
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}