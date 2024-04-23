package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pablom07/go-course/internal/config"
	"github.com/pablom07/go-course/internal/forms"
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

	// Usaremos la sesión para imprimir la IP desde donde se llama a la webapp
	remoteIP := r.RemoteAddr

	// Colocaremos el valor asignándoselo a una clave dentro de la sesión levantada por el cliente
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

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
	
}

// MakeReservation es la renderización de la página de reservas de habitaciones
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
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
