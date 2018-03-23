package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/raymasson/go-mongodb-minikube-kubectl-helm/api/database"
	"github.com/raymasson/go-mongodb-minikube-kubectl-helm/api/person"
)

const welcomeMessage = "Welcome to the person API!"

// NewRouter api router
func NewRouter(mainRouter *mux.Router, mongoSession database.Session) http.Handler {

	apiRouter := mainRouter.PathPrefix("/").Subrouter().StrictSlash(true)

	// Root route
	apiRouter.HandleFunc("/", index).Methods("GET")

	// Managers
	databaseManager := database.NewManager(mongoSession)

	// Handlers
	personHandler := person.NewHandler(databaseManager)

	//Endpoints
	person.NewRouter(apiRouter, personHandler)

	return apiRouter
}

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, welcomeMessage)
}
