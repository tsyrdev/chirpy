package utils

import (
	"net/http"
	"encoding/json"
)

type resError struct {
	Error 	string 	`json:"error"`
}

func RespondWithError (w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resError{Error: message})
}

func RespondWithJSON (w http.ResponseWriter, status int, payload any) {
	w.WriteHeader(status)	
	json.NewEncoder(w).Encode(payload)
}

