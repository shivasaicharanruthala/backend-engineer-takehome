package responder

import (
	"encoding/json"
	"github/shivasaicharanruthala/backend-engineer-takehome/errors"
	"net/http"
)

// SetErrorResponse sends an error response with the appropriate status code and error message.
// It accepts an error and an http.ResponseWriter as parameters and handles different types of custom errors.
func SetErrorResponse(err error, w http.ResponseWriter) {
	switch val := err.(type) {
	case errors.InvalidParam:
		errJson, _ := json.Marshal(val)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(val.StatusCode)
		_, _ = w.Write(errJson)
	case errors.MissingParam:
		errJson, _ := json.Marshal(val)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(val.StatusCode)
		_, _ = w.Write(errJson)
	case errors.EntityNotFound:
		errJson, _ := json.Marshal(val)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(val.StatusCode)
		_, _ = w.Write(errJson)
	case errors.CustomError:
		errJson, _ := json.Marshal(val)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(val.StatusCode)
		_, _ = w.Write(errJson)
	default:
		errJson, _ := json.Marshal(val)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		_, _ = w.Write(errJson)
	}
}

// SetResponse sends a successful response with the specified status code and response body.
// It accepts a response object, a status code, and an http.ResponseWriter as parameters.
func SetResponse(resp interface{}, statusCode int, w http.ResponseWriter) {
	respJson, _ := json.Marshal(resp)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(respJson)
}
