package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tsyrdev/chirpy/utils"
)

type User struct {
	ID 			uuid.UUID 	`json:"id"`
	CreatedAt	time.Time 	`json:"created_at"`
	UpdatedAt	time.Time 	`json:"updated_at"`
	Email		string		`json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	var params struct {
		Email	string	`json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "JSON not valid")
		return 
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), params.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return 
	}
	
	response := User{
		ID: 		user.ID,
		CreatedAt:	user.CreatedAt,
		UpdatedAt: 	user.UpdatedAt,
		Email:		user.Email,
	}

	utils.RespondWithJSON(w, http.StatusCreated, response)
}
