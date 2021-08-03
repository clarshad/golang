package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/clarshad/golang/terraform-as-service/terraform"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

//Handle function handles all the HTTP requests
func Handle(p int) {
	port := ":" + strconv.Itoa(p)
	fmt.Printf("INFO: Started HTTP Server, listening at port %v\n", p)

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/apply", runTerraformHandler).Methods("POST")
	r.HandleFunc("/destroy", runTerraformHandler).Methods("POST")
	r.HandleFunc("/job/{id}", statusHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(port, r))

}

type Config struct {
	Path    string `json:"path"`
	Version string `json:"version"`
	Err     string `json:"error"`
	Status  string `json:"status"`
	Action  string `json:"action"`
	PostResp
}

type PostResp struct {
	RequestId string `json:"request_id"`
}

var currentConfig Config

//runTerraformHandler handles POST request for running terraform configuration
func runTerraformHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("INFO: Server: POST Request at path %v\n", r.URL.Path)

	if currentConfig.RequestId != "" && currentConfig.Status == "RUNNING" {
		fmt.Println("ERROR: Server: 503 Server Busy")
		w.WriteHeader(503)
		json.NewEncoder(w).Encode(currentConfig.PostResp)
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	currentConfig = Config{}
	json.Unmarshal(reqBody, &currentConfig)
	if currentConfig.Version == "" {
		fmt.Println("ERROR: Server: 400 Bad Request")
		w.WriteHeader(400)
		return
	}

	currentConfig.Action = strings.Trim(r.URL.Path, "/")
	currentConfig.RequestId = uuid.NewString()
	currentConfig.Status = "RUNNING"

	go runTerraform()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentConfig.PostResp)
}

func runTerraform() {
	err := terraform.Run(currentConfig.Version, currentConfig.Action, currentConfig.Path)
	if err != nil {
		currentConfig.Status = "ERROR"
		currentConfig.Err = err.Error()
		return
	}

	currentConfig.Status = "SUCCESS"
}

//statusHandler handles GET request to check the status
func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("INFO: Server: GET Request at path %v\n", r.URL.Path)

	params := mux.Vars(r)

	if currentConfig.RequestId == params["id"] {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(currentConfig)
		return
	} else {
		fmt.Println("ERROR: Server: 404 Not Found")
		w.WriteHeader(404)
		return
	}
}
