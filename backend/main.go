package main

import (
	"io"
	"net/http"

	"github.com/gopheracademy/manager/backend/server"
)

func main() {

	server.Run()

}

// statusCodeHandler is useful for testing the server by returning a
// specific HTTP status code.
//  http.Handle("/", statusCodeHandler(http.StatusInternalServerError))
type statusCodeHandler int

func (c statusCodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(int(c))
	io.WriteString(w, http.StatusText(int(c)))
}
