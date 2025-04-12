package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tsyrdev/chirpy/internal/auth"
	"github.com/tsyrdev/chirpy/internal/database"
	"github.com/tsyrdev/chirpy/utils"
)

type Chirp struct {
	ID 			uuid.UUID 	`json:"id"`
	CreatedAt 	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	Body		string		`json:"body"`
	UserID		uuid.UUID	`json:"user_id"`
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	path := r.PathValue("chirpID") 
	chirpID, err := uuid.Parse(path)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Server could not parse the chirp ID")
		return 
	}

	authHeader, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Server couldn't extract the authorization header")
		return 
	}

	userID, err := auth.ValidateJWT(authHeader, cfg.secret)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Server couldn't validate the access token")
		return 
	}

	dbChirp, err := cfg.dbQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Chirp could not be found")
		return 
	}
	if userID != dbChirp.UserID {
		utils.RespondWithError(w, http.StatusForbidden, "Chirp does not belong to user")
		return 
	}

	if err := cfg.dbQueries.DeleteChirp(r.Context(), chirpID); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Server couldn't delete chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := r.PathValue("chirpID")
	id, err := uuid.Parse(path)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Server could not parse the chirp ID")
		return 
	}

	dbChirp, err := cfg.dbQueries.GetChirp(r.Context(), id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not get the Chirp")
		return
	}

	chirp := Chirp{
		ID: dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body: dbChirp.Body,
		UserID: dbChirp.UserID,
	}
	utils.RespondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dbChirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not get Chirps")
		return 
	}

	chirps := make([]Chirp, 0, len(dbChirps))
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID: dbChirp.ID, 
			CreatedAt: 	dbChirp.CreatedAt,
			UpdatedAt: 	dbChirp.UpdatedAt,
			Body: 		dbChirp.Body,
			UserID: 	dbChirp.UserID,
		})
	}
	
	utils.RespondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	var params struct {
		Body	string		`json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "User does not possess a login token")
		return 
	}

	tokenUUID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "401 Unauthorized")
		return 
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if len(params.Body) > 140 {
		utils.RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return 
	}

	cleanChirp := cleanChirp(params.Body)

	chirp, err := cfg.dbQueries.CreateChirps(r.Context(), database.CreateChirpsParams{
		Body:	cleanChirp,
		UserID:	tokenUUID,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not create the chirp: %s", err))
		return 
	}

	utils.RespondWithJSON(w, http.StatusCreated, chirp)	
}

func cleanChirp(chirp string) string {
	badwords := map[string]bool{
		"kerfuffle": 	true,
		"sharbert":		true,
		"fornax":		true,
	}

	words := strings.Fields(chirp)
	for i, word := range words {
		word := strings.ToLower(word)
		if badwords[word] {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
