package forms

/*
	Se define una estructura tipo mapa clave/valor, cuyo valor son slices de strings.
*/

// Lista de errores para cada campo
type errors map[string][]string

// Agrega un nuevo error para el campo especificado
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get devuelve el primer error de cada campo. En caso de que no hayan errores
// precargados, se devuelve un string vacío. 
func (e errors) Get(field string) string {

	// Buscamos el slice de strings que corresponden a ese campo
	es := e[field]
	
	// Si no se encuentra el campo cargado, se devuelve un mensaje vacío.
	if len(es) == 0 {
		return ""
	}

	// Devolvemos el primer resultado encontrado.
	return es[0]

}