// main.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/gorilla/mux"
)

// Article - Our struct for all articles

type Config struct {
	Class          string `json:"connector.class"`
	Tasks          string `json:"tasks.max"`
	Topics         string `json:"topics"`
	File           string `json:"file"`
	FilePattern    string `json:"file.pattern"`
	Converter      string `json:"value.converter"`
	SchemaRegistry string `json:"value.converter.schema.registry.url"`
}
type Connector struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Config Config `json:"config"`
}

var connectors []Connector

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Api Gateway For Cp4D F055 Streaming!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllConnectors(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get("http://kafka-connect-simulator-5-streamingtest.apps.cp4d-poc.cp4d.ichp.nietsnel.nu/connectors")
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(body)

	json.NewEncoder(w).Encode(connectors)
}

func returnSingleConenctor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	for _, connector := range connectors {
		if connector.Id == key {
			json.NewEncoder(w).Encode(connector)
		}
	}
}

func createNewConnector(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	reqBody, _ := ioutil.ReadAll(r.Body)

	var connector Connector
	json.Unmarshal(reqBody, &connector)
	// // update our global Articles array to include
	// // our new Article
	// sanitize connector
	if connector.Config.Class == "" {
		panic("ERROR: Connector class not specified!")
	}
	if connector.Config.Topics == "" {
		panic("ERROR: Missing topic information!")
	}
	if connector.Config.Converter == "" {
		panic("ERROR: Missing Converter information!")
	}
	if connector.Config.SchemaRegistry == "" {
		panic("ERROR: Missing SchemaRegistry information!")
	}
	if strings.Contains(connector.Config.Class, "FileStreamSinkConnector") {
		if connector.Config.File == "" {
			panic("ERROR: Output file not specified")
		}
		if connector.Config.FilePattern == "" {
			connector.Config.FilePattern = "'.'yyyy-MM-dd-HH-mm"
		}
	}

	jsonValue, _ := json.Marshal(connector)
	resp, err := http.Post("http://kafka-connect-simulator-5-streamingtest.apps.cp4d-poc.cp4d.ichp.nietsnel.nu/connectors", "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var responseConnector Connector
	//Convert the body to type string
	json.Unmarshal(body, &responseConnector)

	spew.Dump(responseConnector)
	fmt.Println(resp.Status)
	json.NewEncoder(w).Encode(connectors)

}

func deleteConnector(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for index, connector := range connectors {
		if connector.Id == id {
			connectors = append(connectors[:index], connectors[index+1:]...)
		}
	}

}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/connectors", returnAllConnectors)
	myRouter.HandleFunc("/connector", createNewConnector).Methods("POST")
	myRouter.HandleFunc("/connector/{id}", deleteConnector).Methods("DELETE")
	myRouter.HandleFunc("/connector/{id}", returnSingleConenctor)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {

	handleRequests()
}
