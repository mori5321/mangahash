package handlers

import (
	"fmt"
	"net/http"
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
		respondWithError(w, ErrorResponse{
			MethodNotAllowed,
			[]string{
				fmt.Sprintf("Method %s not allowed for /healthz", r.Method),
			},
		})
	}
}
