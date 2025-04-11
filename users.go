package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tsyrdev/chirpy/internal/auth"
	"github.com/tsyrdev/chirpy/internal/database"
	"github.com/tsyrdev/chirpy/utils"
)

type User struct {
	ID 			uuid.UUID 	`json:"id"`
	CreatedAt	time.Time 	`json:"created_at"`
	UpdatedAt	time.Time 	`json:"updated_at"`
	Email		string		`json:"email"`
}

func (cfg *apiConfig) handlerRevokeRefresh(w http.ResponseWriter, r *http.Request) {
	authHeader, err := auth.GetBearerToken(r.Header)	
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "User does not possess a refresh token")
		return 
	}
	
	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), authHeader)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Unable to revoke refresh")
		return 
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "User does not possess a refresh token")
		return 
	}
	dbToken, err := cfg.dbQueries.GetRefreshToken(r.Context(), refreshToken)	
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Couldn't find the token in the DB")
		return 
	}
	if dbToken.RevokedAt.Valid {
		utils.RespondWithError(w, http.StatusUnauthorized, "The refresh token has been revoked")
		return 
	}

	accessToken, err := auth.MakeJWT(dbToken.UserID, cfg.secret)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Coudln't create a new access token")
		return
	}
		
	utils.RespondWithJSON(w, http.StatusOK, accessToken)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json") 

	var params struct {
		Password 	string	`json:"password"`
		Email		string	`json:"email"`
		ExpiresIn	int		`json:"expires_in_seconds"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return 
	}

	if params.ExpiresIn == 0 {
		params.ExpiresIn = 3600 // default value is 1 hour
	}

	dbUser, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error in the server")
		return 
	}

	if err := auth.CheckPasswordHash(dbUser.HashedPassword, params.Password); err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return 
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.secret)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating login token")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating refresh token")
		return 
	}

	_, err = cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: dbUser.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not create refresh token")
	}

	response := struct{
		ID				uuid.UUID	`json:"id"`
		CreatedAt		time.Time	`json:"created_at"`
		UpdatedAt		time.Time	`json:"updated_at"`
		Email			string		`json:"email"`
		Token			string		`json:"token"`
		RefreshToken	string		`json:"refresh_token"`
	}{
		ID:				dbUser.ID,
		CreatedAt:		dbUser.CreatedAt,
		UpdatedAt:		dbUser.UpdatedAt,
		Email:			dbUser.Email,
		Token:			token,
		RefreshToken: 	refreshToken,
	}
	
	utils.RespondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")

	var params struct {
		Password 	string 	`json:"password"`
		Email		string	`json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "JSON not valid")
		return 
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error saving the password")
		return 
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
	})

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
