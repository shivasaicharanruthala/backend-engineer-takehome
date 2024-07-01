package handler

import (
	"encoding/json"
	er "errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github/shivasaicharanruthala/backend-engineer-takehome/errors"
	"github/shivasaicharanruthala/backend-engineer-takehome/log"
	"github/shivasaicharanruthala/backend-engineer-takehome/model"
	"github/shivasaicharanruthala/backend-engineer-takehome/responder"
	"github/shivasaicharanruthala/backend-engineer-takehome/service"
)

// receiptsHandler is a HTTP handler for receipt-related endpoints.
type receiptsHandler struct {
	logger *log.CustomLogger
	svc    service.Receipts
}

// New creates and returns a new instance of receiptsHandler.
func New(l *log.CustomLogger, svc service.Receipts) *receiptsHandler {
	return &receiptsHandler{
		logger: l,
		svc:    svc,
	}
}

// Get handles HTTP GET requests to retrieve a receipt by its ID.
// It validates the receipt ID and retrieves the receipt points from the service layer.
func (rh *receiptsHandler) Get(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()) > 0 {
		responder.SetErrorResponse(rh.logger, errors.NewCustomError(er.New("query parameter is not accepted"), 400), w, r)

		return
	}

	//	 TODO: check for body

	// Fetch receiptId from path param
	receiptID := mux.Vars(r)["id"]
	if !model.IsValidUUID(receiptID) { // validate id, if not valid throw an error
		responder.SetErrorResponse(rh.logger, errors.NewInvalidParam(errors.InvalidParam{Param: "id"}), w, r)

		return
	}

	// Get service call to fetch points to a receiptID
	receiptPoints, err := rh.svc.Get(receiptID)
	if err != nil {
		responder.SetErrorResponse(rh.logger, err, w, r)

		return
	}

	// Responds with the receipt points if successfully retrieved.
	responder.SetResponse(receiptPoints, 200, w)
	return
}

// Insert handles HTTP POST requests to insert a new receipt.
// It reads and unmarshals the request body, then inserts the receipt through the service layer.
func (rh *receiptsHandler) Insert(w http.ResponseWriter, r *http.Request) {
	var receipt model.Receipt

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// Responds with an error if there's an issue reading the request body.
		err = errors.NewCustomError(err, 400)
		responder.SetErrorResponse(rh.logger, err, w, r)

		return
	}

	err = json.Unmarshal(body, &receipt)
	if err != nil {
		// Responds with an error if there's an issue unmarshalling the request body.
		err = errors.NewCustomError(err, 400)
		responder.SetErrorResponse(rh.logger, err, w, r)

		return
	}

	// service call to insert receipt
	receiptResponse, err := rh.svc.Insert(&receipt)
	if err != nil {
		responder.SetErrorResponse(rh.logger, err, w, r)

		return
	}

	responder.SetResponse(receiptResponse, 201, w)
	return
}

// Health handles HTTP GET requests to check the health status of the service.
// It responds with a simple health check message.
func (rh *receiptsHandler) Health(w http.ResponseWriter, r *http.Request) {
	resp := "{'health': 'ok'}"
	responder.SetResponse(resp, 200, w)
}
