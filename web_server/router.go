package main

import (
	"net/http"
)

// Creamos una estructura Router con las reglas
// o el path donde se realizaran las consultas

type Router struct {
	rules map[string]http.HandlerFunc
}

// Retornar un nuevo Route
func NewRouter() *Router {
	return &Router{
		rules: make(map[string]http.HandlerFunc),
	}
}

// Creamos un buscador de handler para registrar todas las rutas que se encuentran trabajando
func (r *Router) FindHandler(path string) (http.HandlerFunc, bool) {
	handler, exist := r.rules[path]
	return handler, exist
}

//Para ser parte de Handler http debemos implementar el metodo ServeHTTP
//Esta funcion recibe un writer y un request.
func (r *Router) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	handler, exist := r.FindHandler(request.URL.Path)

	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	handler(w, request)

}
