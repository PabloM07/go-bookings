package forms

import (
	"fmt"
	"net/url"
	"strings"

	goval "github.com/asaskevich/govalidator"
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
	return &Form{
		data,

		// Se declara un mapa de slices de strings vacío, de igual manera que se puede
		// declarar como make(map[string][]string)
		errors(map[string][]string{}),
	}
}

// Has revisa si el campo del formulario está en el cuerpo del request y no está vacío.
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}

/*
	Como veremos en esta función, en los parámetros requeridos, usamos unos puntos suspensivos
	por detrás de el tipo de datos que esperamos recibir. Eso se llama una función variádica, y
	significa que la cantidad de datos que puedo recibir de este tipo es variable específicamente
	en ese orden o lugar del input de parámetros, es decir, que puedo pasarle, en este caso, la
	cantidad de strings que querramos, cuyos valores dentro de la función se tratarán como un
	slice de strings.
*/

// Required recopila los campos requeridos del formulario para validar su contenido
func (f *Form) Required(fields ...string) {

	// Recorremos la colección (slice) de nombres de campos
	for _, field := range fields {

		// Obtenemos del formulario el valor del campo que coincide con el nombre que le pasamos
		val := f.Get(field)

		// Quitamos los espacios en blanco del principio y del final de la cadena de string
		if strings.TrimSpace(val) == "" {

			// Si el valor devuelto por la función está vacío, agregamos un error a la lista de
			// errores a mostrar en el cliente
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Valid controla el formulario en busca de errores de validación. Devuelve True si el formulario
// tiene datos válidos o false si existen
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// MinLength valida que el campo dado contenga un dato string que supere una cantidad dada de caracteres
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail verifica si el valor ingresado como cadena string es una dirección de correo electrónico
func (f *Form) IsEmail(field string) bool {

	// Utilizamos un validador externo de la librería 'govalidator' para la verificación del formato
	if !goval.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
	return false
}
