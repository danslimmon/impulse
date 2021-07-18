package server

import (
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/danslimmon/impulse/api"
)

type listHandler struct {
	fs Filesystem
}

func (handler listHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		blopList, err := handler.fs.GetBlopList("_")
		if err != nil {
			sendResponse(
				w,
				http.StatusInternalServerError,
				api.ErrorResponse{err},
			)
		}
		sendResponse(
			w,
			http.StatusOK,
			api.ListResponse{blopList},
		)
		return
	}

	sendResponse(
		w,
		http.StatusMethodNotAllowed,
		api.ErrorResponse{errors.New("Method " + r.Method + " not supported")},
	)
	return
}

func sendResponse(w http.ResponseWriter, header int, response interface{}) error {
	b, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.WithField("error", err.Error()).Error("error marshaling HTTP response to JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	w.WriteHeader(header)
	w.Write(b)
	return nil
}

func serveMux() *http.ServeMux {
	diskFS := &DiskFilesystem{
		rootDir: "/Users/danslimmon/j_workspace/impulse",
	}
	mux := http.NewServeMux()
	mux.Handle("/", listHandler{fs: diskFS})
	return mux
}

func Start() (*http.ServeMux, error) {
	return serveMux(), nil
}
