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

func Start(p int) {
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
	Err       error  `json:"error"`
	Status    string `json:"status"`
}

var currentConfig Config

func runTerraformHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("INFO: Server: POST Request")
	if currentConfig.RequestId != 0 && (currentConfig.Status == "RUNNING" || currentConfig.Status == "") {
		fmt.Println("ERROR: 503 Server busy")
		w.WriteHeader(503)
		json.NewEncoder(w).Encode(currentConfig)
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	currentConfig = Config{}
	json.Unmarshal(reqBody, &currentConfig)
	if currentConfig.Version == "" {
		fmt.Println("ERROR: Server: 400 Bad Request: version missing")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(currentConfig)
		return
	}

	currentConfig.RequestId = rand.Uint32()
	go runTerraform()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentConfig)
}

func runTerraform() {
	currentConfig.Err = terraform.Run(currentConfig.Version, currentConfig.Path)
	if currentConfig.Err == nil {
		currentConfig.Status = "SUCCESS"
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("INFO: Server: GET Request")
	if currentConfig.Err != nil {
		currentConfig.Status = "ERROR"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(currentConfig)
		return
	}

	if currentConfig.RequestId != 0 && currentConfig.Err == nil && currentConfig.Status == "" {
		currentConfig.Status = "RUNNING"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(currentConfig)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentConfig)
}
