package responder

import (
	"encoding/json"
	"net/http"

	"github/shivasaicharanruthala/backend-engineer-takehome/errors"
	"github/shivasaicharanruthala/backend-engineer-takehome/log"
)

// SetErrorResponse sends an error response with the appropriate status code and error message.
// It accepts an error and an http.ResponseWriter as parameters and handles different types of custom errors.
func SetErrorResponse(logger *log.CustomLogger, err error, w http.ResponseWriter, r *http.Request) {
	switch val := err.(type) {
	case errors.InvalidParam:
		errJson, _ := json.Marshal(val)

		lm := log.Message{Level: "ERROR", Method: r.Method, URI: r.RequestURI, StatusCode: val.StatusCode, ErrorMessage: val.Error()}
		logger.Log(&lm)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(val.StatusCode)
		_, _ = w.Write(errJson)
	case errors.MissingParam:
		errJson, _ := json.Marshal(val)

		lm := log.Message{Level: "ERROR", Method: r.Method, URI: r.RequestURI, StatusCode: val.StatusCode, ErrorMessage: val.Error()}
		logger.Log(&lm)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(val.StatusCode)
		_, _ = w.Write(errJson)
	case errors.EntityNotFound:
		errJson, _ := json.Marshal(val)

		lm := log.Message{Level: "ERROR", Method: r.Method, URI: r.RequestURI, StatusCode: val.StatusCode, ErrorMessage: val.Error()}
		logger.Log(&lm)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(val.StatusCode)
		_, _ = w.Write(errJson)
	case errors.CustomError:
		errJson, _ := json.Marshal(val)

		lm := log.Message{Level: "ERROR", Method: r.Method, URI: r.RequestURI, StatusCode: val.StatusCode, ErrorMessage: val.Error()}
		logger.Log(&lm)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(val.StatusCode)
		_, _ = w.Write(errJson)
	default:
		errJson, _ := json.Marshal(val)

		lm := log.Message{Level: "ERROR", Method: r.Method, URI: r.RequestURI, StatusCode: 500, ErrorMessage: val.Error()}
		logger.Log(&lm)

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
