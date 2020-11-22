package main

import "net/http"

// Creamos una estructura que actuara como servidor con los parametros
// del puerto y el router o ruta
type Server struct {
	port   string
	router *Router
}

// En una funcion que modifica el servidor (desde el apuntador *)
// para asi instaciar los valores de un nuevo servidor
func NewServer(port string) *Server {
	return &Server{
		port:   port,
		router: NewRouter(),
	}
}

// Agregamos la funcion que permitira a nuestro servidor agregar rutas a un handler especifico
func (s *Server) Handle(path string, handler http.HandlerFunc) {
	s.router.rules[path] = handler
}

// Esta funcion sera propia del servidor
// Devolvera un error en caso se encuentre un error
func (s *Server) Listen() error {
	http.Handle("/", s.router)
	err := http.ListenAndServe(s.port, nil)
	if err != nil {
		return nil
	}
	return nil
}
