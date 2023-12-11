package common

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type ErrorCode string

const (
	MethodNotAllowed   ErrorCode = "METHOD_NOT_ALLOWED"
	NotFound           ErrorCode = "NOT_FOUND"
	InvalidRequest     ErrorCode = "INVALID_REQUEST"
	EncodingError      ErrorCode = "ENCODING_ERROR"
	InteralServerError ErrorCode = "INTERNAL_SERVER_ERROR"
	UnhandledError     ErrorCode = "UNHANDLED_ERROR"
)

type ErrorResponse struct {
	ErrorCode     ErrorCode `json:"errorCode"`
	ErrorMessages []string  `json:"errorMessages"`
}

func HandleResponse(w http.ResponseWriter, body interface{}, successStatus int, err error) {
	if err != nil {
		HandleError(w, err)
		return
	}

	HandleSuccess(w, body, successStatus)
}

func HandleSuccess(w http.ResponseWriter, body interface{}, status int) {
	respondWithJson(w, body, status)
}

func HandleError(w http.ResponseWriter, err error) {
	var status int
	var errorCode ErrorCode

	if errors.Is(err, NotFoundError) {
		status = http.StatusNotFound
		errorCode = NotFound
	} else if errors.Is(err, InvalidRequestError) {
		status = http.StatusBadRequest
		errorCode = InvalidRequest
	} else if errors.Is(err, MethodNotAllowedError) {
		status = http.StatusMethodNotAllowed
		errorCode = MethodNotAllowed
	} else if errors.Is(err, InternalServerError) {
		status = http.StatusInternalServerError
		errorCode = InteralServerError
	} else if errors.Is(err, InvalidIDError) {
		status = http.StatusBadRequest
		errorCode = InvalidRequest
	} else {
		status = http.StatusInternalServerError
		errorCode = UnhandledError
	}

	resp := ErrorResponse{
		ErrorCode:     errorCode,
		ErrorMessages: []string{err.Error()},
	}

	respondWithJson(w, resp, status)
}

func respondWithJson(w http.ResponseWriter, body interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		respondWithJson(w, ErrorResponse{
			ErrorCode:     EncodingError,
			ErrorMessages: []string{err.Error()},
		}, http.StatusInternalServerError)
	}
}

func GetParams(r *http.Request, prefix string) []string {
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
