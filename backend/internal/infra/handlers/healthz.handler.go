package handlers

import (
	"net/http"

	"github.com/mori5321/mangahash/backend/internal/infra/errs"
)

type HealthzResponse struct {
	Status string `json:"status"`
}

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		respondWithJson(w, HealthzResponse{Status: "OK"}, http.StatusOK)
		w.WriteHeader(http.StatusOK)
	default:
		handleError(w, errs.MethodNotAllowedError)
	}
}
