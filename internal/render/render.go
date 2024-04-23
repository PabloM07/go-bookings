package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/pablom07/go-course/internal/config"
	"github.com/pablom07/go-course/internal/models"
)

var app *config.AppConfig

// NewTemplate setea la config para el paquete template
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData agrega datos adicionales para todos los templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
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
		log.Fatal("Couldn't get template from cache")
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
	}

	// Renderizar el template leído del disco

	_, err = buf.WriteTo(w)

	if err != nil {
		log.Println(err)
	}
}

// CreateTemplateCache se usa para generar una caché de templates de 0, leyendo los
// templates HTML
func CreateTemplateCache() (map[string]*template.Template, error) {
	// myCache := make(map[string]*template.Template)

	// Se crea un caché en blanco para popular desde el principio.
	myCache := map[string]*template.Template{}

	// Obtener todos los archivos nombrados como *.page.tmpl desde ./templates

	pages, err := filepath.Glob("./templates/*.page.tmpl")

	if err != nil {
		return myCache, err
	}

	// Se recorre el slice de nombres de archivos que se obtienen en la variable *pages*

	for _, page := range pages {
		// Se obtienen los nombres de archivos de los templates requeridos
		name := filepath.Base(page)

		// ts: template-set
		ts, err := template.New(name).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")

		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}
