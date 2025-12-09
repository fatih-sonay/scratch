package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {

	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {

	if statusCode > 499 {
		log.Printf("Responding with 5XX error: %s", message)
	}

	type errorStruct struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, statusCode, errorStruct{Error: message})
}
