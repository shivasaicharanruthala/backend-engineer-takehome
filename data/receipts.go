package data

import (
	"sync"

	"github/shivasaicharanruthala/backend-engineer-takehome/errors"
	"github/shivasaicharanruthala/backend-engineer-takehome/log"
	"github/shivasaicharanruthala/backend-engineer-takehome/model"
)

// receiptStore is a thread-safe structure for storing and managing receipts in memory.
type receiptStore struct {
	logger             *log.CustomLogger
	mu                 sync.Mutex               // Mutex to ensure thread-safe access to the in-memory receipt map.
	inMemoryReceiptMap map[string]model.Receipt // In-memory map to store receipts with their IDs as keys.
}

// New creates and returns a new instance of receiptStore which implements methods of the interface Receipts.
func New(l *log.CustomLogger) Receipts {
	return &receiptStore{
		logger:             l,
		inMemoryReceiptMap: make(map[string]model.Receipt),
	}
}

// Get retrieves a receipt from the in-memory store by its ID.
// It returns a ReceiptGetResponse containing the points if the receipt is found,
// otherwise, it returns an error indicating that the receipt was not found.
func (rs *receiptStore) Get(receiptID string) (*model.ReceiptGetResponse, error) {
	rs.mu.Lock()         // Locks the receipt store to ensure thread-safe access.
	defer rs.mu.Unlock() // Unlocks the receipt store after the operation.

	receipt, exists := rs.inMemoryReceiptMap[receiptID]
	if !exists {
		// Returns an error if the receipt is not found.
		return nil, errors.NewEntityNotFound(errors.EntityNotFound{Entity: "receipts", ID: receiptID})
	}

	return &model.ReceiptGetResponse{
		Points: receipt.Points,
	}, nil
}

// Insert adds a new receipt to the in-memory store.
// It returns a ReceiptPostResponse containing the ID of the newly inserted receipt.
func (rs *receiptStore) Insert(receipt *model.Receipt) *model.ReceiptPostResponse {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	rs.inMemoryReceiptMap[receipt.Id] = *receipt

	return &model.ReceiptPostResponse{
		Id: receipt.Id,
	}
}
