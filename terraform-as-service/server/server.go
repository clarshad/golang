package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/clarshad/golang/terraform-as-service/terraform"
	"github.com/gorilla/mux"
)

//Handle function handles all the HTTP requests
func Handle(p int) {
	port := ":" + strconv.Itoa(p)
	fmt.Printf("INFO: Started HTTP Server, listening at port %v\n", p)

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", runTerraformHandler).Methods("POST")
	r.HandleFunc("/", statusHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(port, r))

}

type Config struct {
	Path      string `json:"path"`
	Version   string `json:"version"`
	RequestId uint32 `json:"request_id"`
	Err       string `json:"error"`
	Status    string `json:"status"`
}

var currentConfig Config

//runTerraformHandler handles POST request for running terraform configuration
func runTerraformHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("INFO: Server: POST Request")

	if currentConfig.RequestId != 0 && currentConfig.Status == "RUNNING" {
		fmt.Println("ERROR: Server: 503 Server Busy")
		w.WriteHeader(503)
		json.NewEncoder(w).Encode(currentConfig)
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	currentConfig = Config{}
	json.Unmarshal(reqBody, &currentConfig)
	if currentConfig.Version == "" {
		fmt.Println("ERROR: Server: 400 Bad Request")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(currentConfig)
		return
	}

	go runTerraform()

	currentConfig.RequestId = rand.Uint32()
	currentConfig.Status = "RUNNING"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentConfig)
}

func runTerraform() {
	err := terraform.Run(currentConfig.Version, currentConfig.Path)
	if err != nil {
		currentConfig.Status = "ERROR"
		currentConfig.Err = err.Error()
		return
	}

	currentConfig.Status = "SUCCESS"
}

//statusHandler handles GET request to check the status
func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("INFO: Server: GET Request")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentConfig)
}
