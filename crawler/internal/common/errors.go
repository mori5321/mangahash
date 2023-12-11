package common

import (
	"errors"
)

var NotFoundError = errors.New("NOT_FOUND")
var InvalidRequestError = errors.New("INVALID_REQUEST")
var MethodNotAllowedError = errors.New("METHOD_NOT_ALLOWED")
var EncodeError = errors.New("ENCODING_ERROR")
var InternalServerError = errors.New("INTERNAL_SERVER_ERROR")
var InvalidIDError = errors.New("INVALID_ID")
