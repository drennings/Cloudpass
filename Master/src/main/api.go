package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// API is the object that is responsible for serving the API
type API struct {
	Port string
	Man  *Manager
}

// NewAPI creates a new instance of the API
func NewAPI(port string, manager *Manager) *API {
	return &API{
		Port: port,
		Man:  manager,
	}
}

// Serve starts a webserver with the different handlers
func (api *API) Serve() error {
	// Register handlers
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/status", statusHandler)

	// Start serving
	err := http.ListenAndServe(api.Port, nil)
	return err
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from /")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	res, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Err occurred in status: %v", err)
	}
	fmt.Printf("Received status update:%s", res)
	fmt.Fprintf(w, "Status update received")

	// See if it contains a solution
	c := make(map[string]interface{})

	// unmarschal JSON
	err = json.Unmarshal(res, &c)
	if err != nil {
		fmt.Printf("Could not unmarshal: %v", err)
	}

	solution, ok := c["solution"].(string)
	if !ok {
		fmt.Printf("Could not convert to string: %v", err)
	}

	if len(solution) > 6 {
		fmt.Printf("Found solution: %v", solution)
		os.Exit(0)
	}

}
