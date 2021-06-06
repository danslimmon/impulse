package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/danslimmon/impulse/server"
)

func main() {
	mux, err := server.Start()
	if err != nil {
		log.WithField("error", err.Error()).Error("error starting server")
		os.Exit(1)
	}

	w := new(httptest.ResponseRecorder)
	buf := bytes.NewBuffer([]byte{})
	w.Body = buf

	r, _ := http.NewRequest("GET", "http://localhost/hello", nil)

	mux.ServeHTTP(w, r)
	body, _ := buf.ReadString('\n')
	log.WithField("body", body).Info("response body")
}
