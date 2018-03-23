package person

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/raymasson/go-mongodb-minikube-kubectl-helm/api/database"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Handler : handles http requests on persons
type Handler struct {
	databaseManager database.Manager
}

// NewHandler : initialize the Handler go struct
func NewHandler(d database.Manager) Handler {
	return Handler{
		databaseManager: d,
	}
}

// Person ...
type Person struct {
	ID        bson.ObjectId `bson:"_id" json:"id,omitempty"`
	FirstName string        `json:"first_name"`
	LastName  string        `json:"last_name"`
}

// Get : Gets all the persons
func (h Handler) Get(w http.ResponseWriter, r *http.Request) {

	persons := []Person{}

	// Build the get query
	query := func(c database.Collection) error {
		return c.Find(nil).All(&persons)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Get the persons from the person collection
	err := h.databaseManager.ExecuteQuery(query)
	if err != nil {
		if err == mgo.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(persons); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Post : inserts a person in the person collection
func (h Handler) Post(w http.ResponseWriter, r *http.Request) {

	person := Person{}

	// Decode the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &person); err != nil {
		w.WriteHeader(422)
		return
	}

	// Generate the person ID
	person.ID = bson.NewObjectId()

	// Build the insert query
	query := func(c database.Collection) error {
		return c.Insert(person)
	}

	// Insert the person in the person collection
	err = h.databaseManager.ExecuteQuery(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
