package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ErrorCode string

const (
	MethodNotAllowed ErrorCode = "METHOD_NOT_ALLOWED"
	NotFound         ErrorCode = "NOT_FOUND"
	InvalidRequest   ErrorCode = "INVALID_REQUEST"
	EncodingError    ErrorCode = "ENCODING_ERROR"
)

type ErrorResponse struct {
	ErrorCode     ErrorCode `json:"errorCode"`
	ErrorMessages []string  `json:"errorMessages"`
}

func respondWithJson(w http.ResponseWriter, body interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		respondWithError(w, ErrorResponse{
			ErrorCode:     EncodingError,
			ErrorMessages: []string{err.Error()},
		})
	}
}

func respondWithError(w http.ResponseWriter, e ErrorResponse) {
	var status int
	switch e.ErrorCode {
	case MethodNotAllowed:
		status = http.StatusMethodNotAllowed
	case NotFound:
		status = http.StatusNotFound
	case InvalidRequest:
		status = http.StatusBadRequest
	default:
		status = http.StatusInternalServerError
	}

	respondWithJson(w, e, status)
}

func respondWithEmpty(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func getParams(r *http.Request, prefix string) []string {
	url := r.URL.String()
	path := strings.TrimPrefix(url, prefix)

	splitteds := strings.Split(path, "/")

	var params []string
	for _, splitted := range splitteds {
		if len(splitted) == 0 {
			continue
		}
		params = append(params, splitted)
	}

	return params
}
