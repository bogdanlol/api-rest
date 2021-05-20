// main.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"

	"github.com/gorilla/mux"
)

// Article - Our struct for all articles

type Config struct {
	Class  string `json:"class"`
	Tasks  string `json:"tasks"`
	Topics string `json:"topics"`
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
	fmt.Println("Endpoint Hit: returnAllArticles")
	spew.Dump(connectors)
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
	spew.Dump(connector)
	jsonValue, _ := json.Marshal(connector)
	resp, err := http.Post("localhost:8081/connectors", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err)
	}
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
