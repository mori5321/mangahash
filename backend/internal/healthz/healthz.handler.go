package healthz

import (
	"net/http"

	"github.com/mori5321/mangahash/backend/internal/common"
)

type HealthzResponse struct {
	Status string `json:"status"`
}

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		common.HandleSuccess(w, HealthzResponse{Status: "OK"}, http.StatusOK)
	default:
		common.HandleError(w, common.MethodNotAllowedError)
	}
}
