package main

import (
	"fmt"
	"net/http"
)

func HandRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}
