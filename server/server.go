package server

import (
	"net/http"
)

type helloHandler struct{}

func (helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello\n"))
}

func serveMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", helloHandler{})
	return mux
}

func Start() (*http.ServeMux, error) {
	return serveMux(), nil
}
