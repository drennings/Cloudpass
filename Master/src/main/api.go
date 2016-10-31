package main

import (
	"fmt"
	"net/http"
)

// API is the object that is responsible for serving the API
type API struct {
	Port string
}

// NewAPI creates a new instance of the API
func NewAPI(port string) *API {
	return &API{
		Port: port,
	}
}

// Serve starts a webserver with the different handlers
func (api *API) Serve() error {
	// Register handlers
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/found", foundHandler)

	// Start serving
	err := http.ListenAndServe(api.Port, nil)
	return err
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from /")
}

func foundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from /found")
}
