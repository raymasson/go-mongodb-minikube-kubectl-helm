package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"github.com/rs/cors"

	"github.com/raymasson/go-mongodb-minikube-kubectl-helm/api/config"
	"github.com/raymasson/go-mongodb-minikube-kubectl-helm/api/database"
	"github.com/raymasson/go-mongodb-minikube-kubectl-helm/api/logger"
	"github.com/raymasson/go-mongodb-minikube-kubectl-helm/api/router"
)

func init() {
	//Get configuration
	config.Get()

	//Create a new logger
	logger.New()
}

func main() {
	// Connect to database
	session := database.Connect()
	// Defer the close of the database
	defer session.Close()

	// Setup the API routes
	r := mux.NewRouter().StrictSlash(true)

	// Create Api routes
	router.NewRouter(r, session)

	// AllowAll create a new Cors handler with permissive configuration allowing all
	handler := cors.AllowAll().Handler(r)

	go func() {
		logger.Log.Info("The news service has been started on port 8000")

		log.Fatal(http.ListenAndServe(":8000", handler))
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-signals:
			session.Close()
			log.Fatal("Interruption is detected")
			os.Exit(0)
		}
	}
}
