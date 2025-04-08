package main 

import (
	"encoding/json"
	"net/http"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body	string	`json:"body"`
	}

	type resError struct {
		Error 	string	`json:"error"`
	}

	type resValid struct {
		Valid 	bool	`json:"valid"`
	}
	respondWithError := func(status int, message string) {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(resError{Error: message})
	}

	w.Header().Set("Content-Type", "application/json") 
	defer r.Body.Close()

	var par parameters
	if err := json.NewDecoder(r.Body).Decode(&par); err != nil {
		respondWithError(http.StatusBadRequest, "Invalid JSON")
		return 
	}

	if len(par.Body) > 140 {
		respondWithError(http.StatusBadRequest, "Chirp is too long")
		return 
	}

	json.NewEncoder(w).Encode(resValid{Valid: true})
}

