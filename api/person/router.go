package person

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter ...
func NewRouter(mainRouter *mux.Router, handler Handler) http.Handler {

	mainRouter.Path("/persons").Methods("GET").Handler(http.HandlerFunc(handler.Get))
	mainRouter.Path("/persons").Methods("POST").Handler(http.HandlerFunc(handler.Post))

	return mainRouter
}
