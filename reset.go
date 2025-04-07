package main 

import (
	"net/http"
)

func (ac *apiConfig) handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	ac.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
