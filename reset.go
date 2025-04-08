package main

import (
	"net/http"

	"github.com/tsyrdev/chirpy/utils"
)

const DEVELOPMENT = "dev"

func (ac *apiConfig) handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	if ac.platform != DEVELOPMENT {
		utils.RespondWithError(w, http.StatusForbidden, "Cannot reset database.")
		return 
	}

	ac.dbQueries.ResetUsers(r.Context())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
