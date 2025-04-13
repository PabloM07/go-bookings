package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/pablom07/go-course/internal/config"
	"github.com/pablom07/go-course/internal/forms"
	"github.com/pablom07/go-course/internal/helpers"
	"github.com/pablom07/go-course/internal/models"
	"github.com/pablom07/go-course/internal/render"
)

// Repository es el tipo repositorio de configuraciones
type Repository struct {
	App *config.AppConfig
}

// Repo es el reposotorio usado por los controladores
var Repo *Repository

// NewRepo crea un nuevo repositorio de configuraciones
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers Setea el repositorio para los controladores
func NewHandlers(r *Repository) {
	Repo = r
}

// Home es la renderización de la página de inicio de nuestro sitio
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// Aqui se implementa un poco de lógica de negocios

	// [DEPRECADO] Usaremos la sesión para imprimir la IP desde donde se llama a la webapp
	// remoteIP := r.RemoteAddr

	// [DEPRECADO] Colocaremos el valor asignándoselo a una clave dentro de la sesión levantada por el cliente
	// m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	// Luego debemos mostrar los datos finales al template
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About es la renderización de la página About de nuestro sitio
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	// Escribimos dentro de un mapa los parámetros que debemos pasarle a la Webpage para que los muestre.
	stringMap := make(map[string]string)
	stringMap["test"] = "Sarlanga"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Contact es la renderización de la página de contacto
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// SearchAvailability es la renderización de la página de búsqueda de vacantes de habitaciones
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability es la renderización de la página de búsqueda de vacantes de habitaciones con método post
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {

	// Obtenemos los valores de los inputs del form del template html a través de los valores de los tags 'name'
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	// Escribimos en el response, un slice de bytes que contiene un string a mostrar.
	w.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s", start, end)))
}

/*
Por convención en Golang, la declaración de las estructuras suelen hacerse antecediendo al primer método de uso
inmediato dentro del archivo del código fuente, e incluso, agrupar los métodos que vayan a usar dicho tipo de
datos, facilitando la lectura de las instanciaciones de la estructura y sus usos correspondientes dentro del
comportamiento de nuestra app.
*/

// Struct de prueba para parsear en un formato string JSON para devolver al request POST.
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON es la función RESTful que devuelve una petición JSON
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	//Instanciamos un mock struct a modo de una respuesta a devolver
	resp := jsonResponse{
		OK:      false,
		Message: "Available!",
	}

	// Formateamos la instancia anterior en un string con formato JSON
	out, err := json.MarshalIndent(resp, "", "	")
	if err != nil {
		log.Println(err)
	}

	log.Println(string(out))

	// Agregamos un encabezado al response para determinarle al navegador que la respuesta recibida es
	// en formato JSON
	w.Header().Set("Content-type", "application/json")

	// Escribimos el cuerpo del response en formato Bytes para devolverselo al cliente
	w.Write([]byte(out))

}

// PostReservation maneja la recepción de una petición de reserva de la habitación.
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	// Parseamos la instancia de formulario que le pasamos en el método MakeReservation() y detectamos
	// si existen errores.
	err := r.ParseForm()
	err = errors.New("Este es un mensaje de error")
	if err != nil {
		helpers.ServerError(w, err) // Centralizamos el manejo de errores con los helpers
		return
	}

	// Se instancia una variable del modelo de formulario que definimos en el paquete 'forms' y se
	// le carga todos los datos traidos del formulario del request POST.
	reserv := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
	}

	// Una vez cargados los datos que trajimos del formulario, procederemos a verificar la fiabilidad
	// de éstos.

	form := forms.New(r.PostForm)

	// Se verifica con el método Has() si el campo solicitado del formulario trae datos.
	//form.Has("first_name", r)

	// Verificamos con un método masivo, la validez de los campos requeridos del formulario.
	form.Required("first_name", "last_name", "email")

	// Incorporamos un ejemplo de validación de mínimos caracteres soportados
	form.MinLength("first_name", 3)

	// Implementamos una validación de formato email
	form.IsEmail("email")

	if !form.Valid() {
		// En caso que el formulario contenga errores, se procede a reabrir la página con los datos
		// precargados anteriormente para solicitar su corrección.
		data := make(map[string]any)
		data["reservation"] = reserv

		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	/*
		Hasta aquí tenemos todo el mecanismo por el cual verificamos todos los datos del formulario
		para efectuar la reserva. Lo que debemos lograr ahora, es redirigir a la página del nuevo
		template que creamos para mostrar los datos precargados antes de confirmar la reserva, o bién,
		cuando ésta ya fué confirmada.
	*/

	// Obtenemos la sesión activa para agregarle los datos a mostrar en la página de precarga de la
	// reserva. Como valores, se pasa el contexto del request, la clave en string, y el valor del tipo
	// 'any' o 'interface{}'
	m.App.Session.Put(r.Context(), "reservation", reserv)

	/*
		Luego de subir el formulario de la reserva con sus datos pertinentes, debemos evitar que se repita
		la recepción inmediata de los mismos dentro de la misma sesión para evitar cargar el formulario,
		por lo que debemos indicarle a la sesión que nos redirija a la página que muestra los datos de la
		nueva reserva
	*/
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// MakeReservation es la renderización de la página de reservas de habitaciones
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	// Se crea una instancia de la struct que nos trae el esquema del formulario necesario para generar
	// la reserva
	var emptyRes models.Reservation

	// Se crea un mapa de claves-valores que poseen un alias a un dato de cualquier tipo que interactúe
	// con el cliente
	data := make(map[string]any)

	// Se le agrega el formulario vacío declarado anteriormente como abstracción del formulario a mostrar
	data["reservation"] = emptyRes
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		/*
			Como nuestra página de reservación tiene una acción que se ejecuta con la carga y el envio de
			ciertos datos, debemos pasarle al template la abstracción de dicho formulario para el check de
			datos del lado del servidor (en caso de que el cliente tenga JS desactivado de su navegador),
			además de performar la acción en concreto de la reserva.
		*/
		Form: forms.New(nil),

		// Se adjunta el arreglo de datos que debemos pasarle al controlador del cliente
		Data: data,
	})
}

// Generals es la renderización de la página General's Suite de nuestro sitio
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors es la renderización de la página Major's Room de nuestro sitio
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// ReserveSum renderiza la página del resumen de reserva con los datos traidos del formulario
func (m *Repository) ReservSum(w http.ResponseWriter, r *http.Request) {
	/*
		Antes de renderizar la página, debemos obtener los datos del formulario cargado que obtuvimos de
		parte de la redirección de página, luego de haber confirmado la carga de datos desde la sección
		anterior. Para esto, utilizamos el método Get() de la sesión para pedirle un valor almacenado.
		Como este método devuelve como tipo de datos un 'any' o 'interface{}' debemos escribir un tipo
		de aserción, es decir, una especificación del tipo de datos que se espera recibir desde esta
		sesión, de manera de mantener un check sobre el tipo de datos que quiero obtener para poder
		manipularlo dentro del método, cuya sintaxis x.(T), donde x es el valor que quieres verificar
		y T es el tipo que estás comprobando.
	*/
	reserv, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	// Comprobamos que el tipo de aserción sea coincidente
	if !ok {
		log.Println("Cannot get item from session")

		// Agrego un mensaje de error a la sesión para mostrarle al usuario
		m.App.Session.Put(r.Context(), "error", "Can't got reservation from session")

		// Redirigimos a la home page con un código HTTP 307
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// Quitamos el valor "reservation" de la sesión de nuestro context
	m.App.Session.Remove(r.Context(), "reservation")

	// Instanciamos un mapa de clave-valor para cargar nuestro formulario con datos ya verificados
	data := make(map[string]any)

	// Cargamos el formulario a mostrar
	data["reservation"] = reserv

	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{

		// Lo incluimos como parámetro para mostrar en respuesta a la petición HTTP.
		Data: data,
	})

}

func (m *Repository) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
