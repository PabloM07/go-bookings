package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/pablom07/go-course/internal/config"
	"github.com/pablom07/go-course/internal/models"
)

// Funciones que le queremos atribuir a nuestros templates html
var funcs = template.FuncMap{}

var app *config.AppConfig
var tmplPath string = "./templates"

/*
	En este ejemplo, utilizaremos el paquete 'http' propio de Golang para renderizar los templates
	html con los que vamos a trabajar en esta Webapp. Aún así, una sugerencia por parte del docente
	es jugar y aprender con la siguiente librería:

	http://github.com/CloudyKit/jet
*/

// NewTemplate setea la config para el paquete template
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData agrega datos adicionales para todos los templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {

	// Parámetro que se obtiene del TemplateData para verificar si mostrar alertas de distintos alias
	// información que se muestra y desaparece cuando la página que lo muestra se refresca.
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")

	// Generación del token CSRF para prevenir accesos maliciosos
	td.CSRFToken = nosurf.Token(r)

	return td
}

// Función que renderiza los templates HTML para Golang. Para poder armar la prueba unitaria, debemos
// evitar que la ejecución se interrumpa a causa de algún error y hacer que devuelva el error para
// analizarlo.
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	// Crear una caché de templates
	var tc map[string]*template.Template

	// Generamos un check para poder utilizar la caché de templates
	if app.UseCache {
		// obtenemos el caché de la variable de config de la webapp
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	// Obtener el template requerido de caché
	t, ok := tc[tmpl]

	if !ok {
		// Imprimimos el mensaje en el log
		log.Println("can't get template from cache!")

		// Devolvemos el error para analizarlo en pruebas unitarias.
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)

	/* Cuando el template se ejecuta, le pasamos como parámetros los datos que almacenamos en
	   el buffer, y como segundo parámetro, le pasamos una struct que posee las variables que
	   contienen los resultados de las operaciones de negocio implementadas en los controladores
	   o handlers
	*/

	// Agregamos datos adicionales a los templates a usar
	td = AddDefaultData(td, r)

	err := t.Execute(buf, td)

	if err != nil {
		log.Println(err)
		return errors.New("can't get template from cache")
	}

	// Renderizar el template leído del disco
	_, err = buf.WriteTo(w)

	if err != nil {
		log.Fatal(err)
	}

	// Al final de la ejecución, concluimos devolviendo nulo para errores.
	return nil
}

// CreateTemplateCache se usa para generar una caché de templates de 0, leyendo los
// templates HTML
func CreateTemplateCache() (map[string]*template.Template, error) {
	// myCache := make(map[string]*template.Template)

	// Se crea un caché en blanco para popular desde el principio.
	myCache := map[string]*template.Template{}

	// Obtener todos los archivos nombrados como *.page.tmpl desde ./templates

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", tmplPath))

	if err != nil {
		return myCache, err
	}

	// Se recorre el slice de nombres de archivos que se obtienen en la variable *pages*

	for _, page := range pages {
		// Se obtienen los nombres de archivos de los templates requeridos
		name := filepath.Base(page)

		// ts: template-set
		ts, err := template.New(name).Funcs(funcs).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", tmplPath))

		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", tmplPath))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}
