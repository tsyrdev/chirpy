package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body	string	`json:"body"`
	}

	type resError struct {
		Error 	string	`json:"error"`
	}

	type resValid struct {
		CleanedBody	string	`json:"cleaned_body"`
	}

	respondWithError := func(status int, message string) {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(resError{Error: message})
	}

	respondWithJSON := func(status int, payload interface{}) {
		w.WriteHeader(status)	
		json.NewEncoder(w).Encode(payload)
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


	res := cleanChirp(par.Body)
	respondWithJSON(http.StatusOK, resValid{CleanedBody: res})
}

func cleanChirp(chirp string) string {
	words := strings.Fields(chirp)
	for i, word := range words {
		word := strings.ToLower(word)
		if word == "kerfuffle" || word == "sharbert" || word == "fornax" {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
