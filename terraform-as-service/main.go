package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	Path    string `json:"path"`
	Version string `json:"version"`
}

func runTerraform(config Config) error {
	if config.Version != "1.0" {
		return fmt.Errorf("version in not 1.0, instead it is %v", config.Version)
	} else {
		return nil
	}
}

func runTerraformHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var config Config
	json.Unmarshal(reqBody, &config)
	err := runTerraform(config)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		json.NewEncoder(w).Encode(config)
	}
}

func startServer() {
	fmt.Println("starting http server...")
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", runTerraformHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func main() {
	startServer()
}
