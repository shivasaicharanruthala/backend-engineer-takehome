package service

import (
	"github.com/google/uuid"
	"github/shivasaicharanruthala/backend-engineer-takehome/data"
	"github/shivasaicharanruthala/backend-engineer-takehome/log"
	"github/shivasaicharanruthala/backend-engineer-takehome/model"
)

// receiptsService is a service layer structure for handling receipt-related operations.
type receiptsService struct {
	logger    *log.CustomLogger
	dataStore data.Receipts // Data layer interface for interacting with the receipt data store.
}

// New creates and returns a new instance of receiptsService which implements all methods of the interface service.Receipts.
func New(l *log.CustomLogger, ds data.Receipts) Receipts {
	return &receiptsService{
		logger:    l,
		dataStore: ds,
	}
}

// Get retrieves a receipt from the data store by its ID.
// It returns a ReceiptGetResponse containing the points if the receipt is found,
// otherwise, it returns an error indicating that the receipt was not found.
func (rs receiptsService) Get(receiptID string) (*model.ReceiptGetResponse, error) {
	return rs.dataStore.Get(receiptID)
}

// Insert adds a new receipt to the data store after validating and calculating its points.
// It validates the receipt payload, calculates the receipt points, generates a new UUID for the receipt,
// and then inserts it into the data store. It returns a ReceiptPostResponse containing the ID of the newly inserted receipt.
func (rs receiptsService) Insert(receipt *model.Receipt) (*model.ReceiptPostResponse, error) {
	// Validates the receipt payload.
	err := receipt.PayloadValidation()
	if err != nil {
		return nil, err
	}

	// Calculates the points for the receipt.
	if err = receipt.CalculateTotalReceiptPoints(); err != nil {
		return nil, err
	}

	// Generates a new UUID for the receipt.
	receipt.Id = uuid.New().String()

	return rs.dataStore.Insert(receipt), nil
}
