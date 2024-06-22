package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var tests = []struct {
	name          string
	url           string
	method        string
	params        []postData
	expStatusCode int
}{
	// GET requests tests
	{"home",	"/",					"GET", []postData{}, http.StatusOK},
	{"about",	"/about",				"GET", []postData{}, http.StatusOK},
	{"gq",		"/generals-quarters",	"GET", []postData{}, http.StatusOK},
	{"ms", 		"/majors-suite",		"GET", []postData{}, http.StatusOK},
	{"sa",		"/search-availability",	"GET", []postData{}, http.StatusOK},
	{"contact",	"/contact",				"GET", []postData{}, http.StatusOK},
	{"mr",		"/make-reservation",	"GET", []postData{}, http.StatusOK},

	// POST requests tests
	{"p_sa",	"/search-availability-json",	"POST", []postData{
		{key: "start",	value: "2024-07-07"},
		{key: "end",	value: "2024-07-08"},
	}, http.StatusOK},
	{"p_mr",	"/make-reservation",			"POST", []postData{
		{key: "first-name",	value: "John"},
		{key: "last-name",	value: "Smith"},
		{key: "email",	value: "me@here.com"},
		{key: "phone",	value: "555-555-5555"},
	}, http.StatusOK},
	
	
}

// Prueba unitaria donde se crea el servicio HTTP para probar con un cliente automatizado
func TestHandlers(t *testing.T) {

	// Obtenemos las rutas que armamos en el setup de las pruebas
	routes := getRoutes()

	/* Definimos un servicio que emula un cliente para llevar a cabo nuestras pruebas.
	   Le pasamos como parámetros, las rutas a probar */ 
	ts := httptest.NewTLSServer(routes)
	defer ts.Close() // Cuando terminan las pruebas, debemos cerrar el server.

	// Ejecutamos las pruebas con un bucle
	for _, e := range tests {

		// Controlamos el verbo del request HTTP de la prueba
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			//Controlamos el código HTTP que nos devuelve el endpoint que estamos probando
			if resp.StatusCode != e.expStatusCode {
				t.Errorf("For %s expected %d but got %d", e.name, e.expStatusCode, resp.StatusCode)
			}
		} else {
			// Construimos una variable cuya configuración sea la del dato que esperamos recibir
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL + e.url, values)

			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			//Controlamos el código HTTP que nos devuelve el endpoint que estamos probando
			if resp.StatusCode != e.expStatusCode {
				t.Errorf("For %s expected %d but got %d", e.name, e.expStatusCode, resp.StatusCode)
			}
		}
	}

}