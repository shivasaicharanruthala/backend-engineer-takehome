package data

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github/shivasaicharanruthala/backend-engineer-takehome/errors"
	"github/shivasaicharanruthala/backend-engineer-takehome/log"
	"github/shivasaicharanruthala/backend-engineer-takehome/model"
)

func NewTest() *receiptStore {
	logger, _ := log.NewCustomLogger("test.log") // Create a real logger for consistency in tests.
	return &receiptStore{
		logger:             logger,
		inMemoryReceiptMap: make(map[string]model.Receipt),
	}
}

func TestDataStoreGet(t *testing.T) {
	store := NewTest()

	receiptID := "4a77ec9d-5334-43d0-a9e1-4fca8807bf8f"
	receipt := model.Receipt{Id: receiptID, Points: 10}
	store.inMemoryReceiptMap[receiptID] = receipt

	testcases := []struct {
		id                      int
		useCase                 string
		receiptID               string
		expectedReceiptResponse *model.ReceiptGetResponse
		expectedError           error
	}{
		{
			id: 1, useCase: "Positive case: Fetch Existing Receipt",
			receiptID:               receiptID,
			expectedReceiptResponse: &model.ReceiptGetResponse{Points: 10},
			expectedError:           nil,
		},
		{
			id: 2, useCase: "Negative case: Fetch Non Existing Receipt",
			receiptID:               "5a77ec9d-5334-43d0-a9e1-4fca8807bf8f",
			expectedReceiptResponse: nil,
			expectedError:           errors.EntityNotFound{Entity: "receipts", ID: "5a77ec9d-5334-43d0-a9e1-4fca8807bf8f"},
		},
	}

	for _, tc := range testcases {
		resp, err := store.Get(tc.receiptID)

		assert.Equal(t, tc.expectedReceiptResponse, resp, fmt.Sprintf("Test %v Failed with use case %v", tc.id, tc.useCase))
		if err != nil {
			// timestamps in the error cant can be compared to comparing error strings
			assert.Equal(t, tc.expectedError.Error(), err.Error(), fmt.Sprintf("Test %v Failed with use case %v", tc.id, tc.useCase))
		} else {
			assert.NoError(t, tc.expectedError, err, fmt.Sprintf("Test %v Failed with use case %v", tc.id, tc.useCase))
		}
	}
}

func TestDataStoreInsert(t *testing.T) {
	store := NewTest()

	testcases := []struct {
		id                      int
		useCase                 string
		receipt                 *model.Receipt
		expectedReceiptResponse *model.ReceiptPostResponse
	}{
		{
			id: 1, useCase: "New Receipt",
			receipt:                 &model.Receipt{Id: "4a77ec9d-5334-43d0-a9e1-4fca8807bf8f"},
			expectedReceiptResponse: &model.ReceiptPostResponse{Id: "4a77ec9d-5334-43d0-a9e1-4fca8807bf8f"},
		},
		{
			id: 2, useCase: "Duplicate Receipt",
			receipt:                 &model.Receipt{Id: "4a77ec9d-5334-43d0-a9e1-4fca8807bf8f"},
			expectedReceiptResponse: &model.ReceiptPostResponse{Id: "4a77ec9d-5334-43d0-a9e1-4fca8807bf8f"},
		},
	}

	for _, tc := range testcases {
		resp := store.Insert(tc.receipt)

		assert.Equal(t, tc.expectedReceiptResponse, resp, fmt.Sprintf("Test %v Failed with use case %v", tc.id, tc.useCase))
	}
}

// TestConcurrencyInsert tests the concurrent insertion of receipts into the store.
func TestConcurrencyInsert(t *testing.T) {
	// Create a new test store
	store := NewTest()

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Define a list of receipts to be inserted concurrently
	receipts := []*model.Receipt{
		{Id: "1a77ec9d-5334-43d0-a9e1-4fca8807bf8f", Points: 10},
		{Id: "2a77ec9d-5334-43d0-a9e1-4fca8807bf8f", Points: 20},
		{Id: "3a77ec9d-5334-43d0-a9e1-4fca8807bf8f", Points: 30},
		{Id: "4a77ec9d-5334-43d0-a9e1-4fca8807bf8f", Points: 40},
		{Id: "5a77ec9d-5334-43d0-a9e1-4fca8807bf8f", Points: 50},
	}

	wg.Add(len(receipts))

	// concurrently insert different receipts
	go func() {
		defer wg.Done()
		_ = store.Insert(receipts[0])
	}()
	go func() {
		defer wg.Done()
		_ = store.Insert(receipts[1])
	}()
	go func() {
		defer wg.Done()
		_ = store.Insert(receipts[2])
	}()
	go func() {
		defer wg.Done()
		_ = store.Insert(receipts[3])
	}()
	go func() {
		defer wg.Done()
		_ = store.Insert(receipts[4])
	}()

	wg.Wait()

	// Verify that all receipts have been inserted correctly
	for _, receipt := range receipts {
		storedReceipt, exists := store.inMemoryReceiptMap[receipt.Id]
		assert.True(t, exists)
		assert.Equal(t, receipt, &storedReceipt)
	}
}
