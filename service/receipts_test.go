package service

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	store "github/shivasaicharanruthala/backend-engineer-takehome/data"
	"github/shivasaicharanruthala/backend-engineer-takehome/errors"
	"github/shivasaicharanruthala/backend-engineer-takehome/log"
	"github/shivasaicharanruthala/backend-engineer-takehome/model"
)

func TestServiceGet(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger, _ := log.NewCustomLogger("test.log")
	receiptStore := store.NewMockReceipts(ctrl)
	receiptService := New(logger, receiptStore)

	testCases := []struct {
		id                    int
		useCase               string
		receiptID             string
		actualReceiptOutput   *model.ReceiptGetResponse
		actualError           error
		expectedReceiptOutput *model.ReceiptGetResponse
		expectedError         error
	}{
		{
			id: 1, useCase: "Fetch Existing Receipt",
			receiptID:             "4a77ec9d-5334-43d0-a9e1-4fca8807bf8f",
			actualReceiptOutput:   &model.ReceiptGetResponse{Points: 10},
			actualError:           nil,
			expectedReceiptOutput: &model.ReceiptGetResponse{Points: 10},
			expectedError:         nil,
		},
		{
			id: 2, useCase: "Fetch Non Existing Receipt",
			receiptID:             "5a77ec9d-5334-43d0-a9e1-4fca8807bf8f",
			actualReceiptOutput:   nil,
			actualError:           errors.EntityNotFound{Entity: "receipts", ID: "5a77ec9d-5334-43d0-a9e1-4fca8807bf8f"},
			expectedReceiptOutput: nil,
			expectedError:         errors.EntityNotFound{Entity: "receipts", ID: "5a77ec9d-5334-43d0-a9e1-4fca8807bf8f"},
		},
	}

	for _, tc := range testCases {
		receiptStore.EXPECT().Get(tc.receiptID).Return(tc.actualReceiptOutput, tc.actualError)

		receiptResp, err := receiptService.Get(tc.receiptID)
		assert.Equal(t, tc.expectedReceiptOutput, receiptResp)
		if err != nil {
			assert.Equal(t, tc.expectedError.Error(), err.Error())
		} else {
			assert.Equal(t, tc.expectedError, err)
		}
	}
}

func TestServiceInsert_Failure(t *testing.T) {
	logger, _ := log.NewCustomLogger("test.log")
	receiptStore := store.New(logger)
	receiptService := New(logger, receiptStore)

	testCases := []struct {
		id            int
		useCase       string
		receipt       *model.Receipt
		expectedError error
	}{
		{
			id: 1, useCase: "Missing Retailer in the payload",
			receipt: &model.Receipt{
				PurchaseDate: model.StringPointer("2024-05-20"),
				PurchaseTime: model.StringPointer("15:04"),
				Total:        model.StringPointer("35.35"),
				Items: []model.Item{
					{ShortDescription: model.StringPointer("Mountain Dew 12PK"), Price: model.StringPointer("6.49")},
				},
			},
			expectedError: errors.MissingParam{Param: "retailer"},
		},
		{
			id: 2, useCase: "Missing PurchaseDate in the payload",
			receipt: &model.Receipt{
				Retailer:     model.StringPointer("Target"),
				PurchaseTime: model.StringPointer("15:04"),
				Total:        model.StringPointer("35.35"),
				Items: []model.Item{
					{ShortDescription: model.StringPointer("Mountain Dew 12PK"), Price: model.StringPointer("6.49")},
				},
			},
			expectedError: errors.MissingParam{Param: "purchaseDate"},
		},
		{
			id: 3, useCase: "Missing PurchaseTime in the payload",
			receipt: &model.Receipt{
				Retailer:     model.StringPointer("Target"),
				PurchaseDate: model.StringPointer("2024-05-20"),
				Total:        model.StringPointer("35.35"),
				Items: []model.Item{
					{ShortDescription: model.StringPointer("Mountain Dew 12PK"), Price: model.StringPointer("6.49")},
				},
			},
			expectedError: errors.MissingParam{Param: "purchaseTime"},
		},
		{
			id: 4, useCase: "Missing Items in the payload",
			receipt: &model.Receipt{
				Retailer:     model.StringPointer("Target"),
				PurchaseDate: model.StringPointer("2024-05-20"),
				PurchaseTime: model.StringPointer("15:04"),
				Total:        model.StringPointer("35.35"),
			},
			expectedError: errors.MissingParam{Param: "items"},
		},
		{
			id: 5, useCase: "Missing shortDescription in the payload",
			receipt: &model.Receipt{
				Retailer:     model.StringPointer("Target"),
				PurchaseDate: model.StringPointer("2024-05-20"),
				PurchaseTime: model.StringPointer("15:04"),
				Total:        model.StringPointer("35.35"),
				Items:        []model.Item{{Price: model.StringPointer("6.49")}},
			},
			expectedError: errors.MissingParam{Param: "shortDescription"},
		},
		{
			id: 6, useCase: "Missing price in the payload",
			receipt: &model.Receipt{
				Retailer:     model.StringPointer("Target"),
				PurchaseDate: model.StringPointer("2024-05-20"),
				PurchaseTime: model.StringPointer("15:04"),
				Total:        model.StringPointer("35.35"),
				Items:        []model.Item{{ShortDescription: model.StringPointer("Mountain Dew 12PK")}}},
			expectedError: errors.MissingParam{Param: "price"},
		},
		{
			id: 7, useCase: "Missing total in the payload",
			receipt: &model.Receipt{
				Retailer:     model.StringPointer("Target"),
				PurchaseDate: model.StringPointer("2022-01-01"),
				PurchaseTime: model.StringPointer("13:01"),
				Items: []model.Item{
					{ShortDescription: model.StringPointer("Mountain Dew 12PK"), Price: model.StringPointer("6.49")},
				},
			},
			expectedError: errors.MissingParam{Param: "total"},
		},
		{
			id: 8, useCase: "Invalid PurchaseDate in the payload",
			receipt: &model.Receipt{
				Retailer:     model.StringPointer("Target"),
				PurchaseDate: model.StringPointer("01-01"),
				PurchaseTime: model.StringPointer("13:01"),
				Total:        model.StringPointer("35.35"),
				Items: []model.Item{
					{ShortDescription: model.StringPointer("Mountain Dew 12PK"), Price: model.StringPointer("6.49")},
				},
			},
			expectedError: errors.InvalidParam{Param: "purchaseDate"},
		},
		{
			id: 9, useCase: "Invalid PurchaseTime in the payload",
			receipt: &model.Receipt{
				Retailer:     model.StringPointer("Target"),
				PurchaseDate: model.StringPointer("2022-01-01"),
				PurchaseTime: model.StringPointer("01"),
				Total:        model.StringPointer("35.35"),
				Items: []model.Item{
					{ShortDescription: model.StringPointer("Mountain Dew 12PK"), Price: model.StringPointer("6.49")},
				},
			},
			expectedError: errors.InvalidParam{Param: "purchaseTime"},
		},
	}

	for _, tc := range testCases {
		_, err := receiptService.Insert(tc.receipt)
		assert.Equal(t, tc.expectedError.Error(), err.Error())
	}
}

func TestServiceInsert_Success(t *testing.T) {
	logger, _ := log.NewCustomLogger("test.log")
	receiptStore := store.New(logger)
	receiptService := New(logger, receiptStore)

	testCases := []struct {
		id             int
		useCase        string
		receipt        *model.Receipt
		expectedPoints int
	}{
		{
			id: 1, useCase: "New Receipt: purchaseTime before 2pm, purchase day as odd, odd number of items",
			receipt: &model.Receipt{
				Retailer:     model.StringPointer("Target"),
				PurchaseDate: model.StringPointer("2022-01-01"),
				PurchaseTime: model.StringPointer("13:01"),
				Total:        model.StringPointer("35.35"),
				Items: []model.Item{
					{ShortDescription: model.StringPointer("Mountain Dew 12PK"), Price: model.StringPointer("6.49")},
					{ShortDescription: model.StringPointer("Emils Cheese Pizza"), Price: model.StringPointer("12.25")},
					{ShortDescription: model.StringPointer("Knorr Creamy Chicken"), Price: model.StringPointer("1.26")},
					{ShortDescription: model.StringPointer("Doritos Nacho Cheese"), Price: model.StringPointer("3.35")},
					{ShortDescription: model.StringPointer("   Klarbrunn 12-PK 12 FL OZ  "), Price: model.StringPointer("12.00")},
				},
			},
			expectedPoints: 28,
		},
		{
			id: 2, useCase: "New Receipt: purchaseTime after 4pm, purchase day as odd",
			receipt: &model.Receipt{
				Retailer:     model.StringPointer("Target"),
				PurchaseDate: model.StringPointer("2022-01-01"),
				PurchaseTime: model.StringPointer("17:01"),
				Total:        model.StringPointer("35.35"),
				Items: []model.Item{
					{ShortDescription: model.StringPointer("Mountain Dew 12PK"), Price: model.StringPointer("6.49")},
					{ShortDescription: model.StringPointer("Emils Cheese Pizza"), Price: model.StringPointer("12.25")},
					{ShortDescription: model.StringPointer("Knorr Creamy Chicken"), Price: model.StringPointer("1.26")},
					{ShortDescription: model.StringPointer("Doritos Nacho Cheese"), Price: model.StringPointer("3.35")},
					{ShortDescription: model.StringPointer("   Klarbrunn 12-PK 12 FL OZ  "), Price: model.StringPointer("12.00")},
				},
			},
			expectedPoints: 28,
		},
		{
			id: 3, useCase: "New Receipt: purchaseTime between 2pm to 4pm, purchase day as odd",
			receipt: &model.Receipt{
				Retailer:     model.StringPointer("Target"),
				PurchaseDate: model.StringPointer("2022-01-01"),
				PurchaseTime: model.StringPointer("15:01"),
				Total:        model.StringPointer("35.35"),
				Items: []model.Item{
					{ShortDescription: model.StringPointer("Mountain Dew 12PK"), Price: model.StringPointer("6.49")},
					{ShortDescription: model.StringPointer("Emils Cheese Pizza"), Price: model.StringPointer("12.25")},
					{ShortDescription: model.StringPointer("Knorr Creamy Chicken"), Price: model.StringPointer("1.26")},
					{ShortDescription: model.StringPointer("Doritos Nacho Cheese"), Price: model.StringPointer("3.35")},
					{ShortDescription: model.StringPointer("   Klarbrunn 12-PK 12 FL OZ  "), Price: model.StringPointer("12.00")},
				},
			},
			expectedPoints: 38,
		},
		{
			id: 4, useCase: "New Receipt: purchase day as even",
			receipt: &model.Receipt{
				Retailer:     model.StringPointer("Target"),
				PurchaseDate: model.StringPointer("2022-01-02"),
				PurchaseTime: model.StringPointer("15:01"),
				Total:        model.StringPointer("35.35"),
				Items: []model.Item{
					{ShortDescription: model.StringPointer("Mountain Dew 12PK"), Price: model.StringPointer("6.49")},
					{ShortDescription: model.StringPointer("Emils Cheese Pizza"), Price: model.StringPointer("12.25")},
					{ShortDescription: model.StringPointer("Knorr Creamy Chicken"), Price: model.StringPointer("1.26")},
					{ShortDescription: model.StringPointer("Doritos Nacho Cheese"), Price: model.StringPointer("3.35")},
					{ShortDescription: model.StringPointer("   Klarbrunn 12-PK 12 FL OZ  "), Price: model.StringPointer("12.00")},
				},
			},
			expectedPoints: 32,
		},
		{
			id: 5, useCase: "New Receipt: trimmed length of item short description has multiple & non multiple of 3, with extra spaces",
			receipt: &model.Receipt{
				Retailer:     model.StringPointer("Target"),
				PurchaseDate: model.StringPointer("2022-01-02"),
				PurchaseTime: model.StringPointer("15:01"),
				Total:        model.StringPointer("35.35"),
				Items: []model.Item{
					{ShortDescription: model.StringPointer("Mountain Dew 12PK"), Price: model.StringPointer("6.49")},
					{ShortDescription: model.StringPointer("Emils Cheese Pizza"), Price: model.StringPointer("12.25")},
					{ShortDescription: model.StringPointer("Knorr Creamy Chicken"), Price: model.StringPointer("1.26")},
					{ShortDescription: model.StringPointer("Doritos Nacho Cheese"), Price: model.StringPointer("3.35")}, {
						ShortDescription: model.StringPointer("   Klarbrunn 12-PK 12 FL OZ  "), Price: model.StringPointer("12.00")},
				},
			},
			expectedPoints: 32,
		},
		{
			id: 6, useCase: "New Receipt: only one item on receipt",
			receipt: &model.Receipt{
				Retailer:     model.StringPointer("Target"),
				PurchaseDate: model.StringPointer("2022-01-02"),
				PurchaseTime: model.StringPointer("13:01"),
				Total:        model.StringPointer("6.49"),
				Items:        []model.Item{{ShortDescription: model.StringPointer("Mountain Dew 12PK"), Price: model.StringPointer("6.49")}},
			},
			expectedPoints: 6,
		},
		{
			id: 7, useCase: "New Receipt: even number of items on receipt",
			receipt: &model.Receipt{
				Retailer:     model.StringPointer("Target"),
				PurchaseDate: model.StringPointer("2022-01-02"),
				PurchaseTime: model.StringPointer("15:01"),
				Total:        model.StringPointer("18.74"),
				Items: []model.Item{
					{ShortDescription: model.StringPointer("Mountain Dew 12PK"), Price: model.StringPointer("6.49")},
					{ShortDescription: model.StringPointer("Emils Cheese Pizza"), Price: model.StringPointer("12.25")},
				},
			},
			expectedPoints: 24,
		},
		{
			id: 8, useCase: "New Receipt: total is multiple of 0.25, has no cents",
			receipt: &model.Receipt{
				Retailer:     model.StringPointer("Target"),
				PurchaseDate: model.StringPointer("2022-01-02"),
				PurchaseTime: model.StringPointer("13:01"),
				Total:        model.StringPointer("5.00"),
				Items: []model.Item{
					{ShortDescription: model.StringPointer("Mountain Dew 12PK"), Price: model.StringPointer("5.00")},
				},
			},
			expectedPoints: 81,
		},
		{
			id: 9, useCase: "New Receipt: total is non multiple of 0.25, has cents",
			receipt: &model.Receipt{
				Retailer:     model.StringPointer("Target"),
				PurchaseDate: model.StringPointer("2022-01-02"),
				PurchaseTime: model.StringPointer("13:01"),
				Total:        model.StringPointer("5.12"),
				Items: []model.Item{
					{ShortDescription: model.StringPointer("Mountain Dew 12PK"), Price: model.StringPointer("5.12")},
				},
			},
			expectedPoints: 6,
		},
		{
			id: 10, useCase: "New Receipt: retailer has alphanumeric characters & special character",
			receipt: &model.Receipt{
				Retailer:     model.StringPointer("M&M Corner Market 32"),
				PurchaseDate: model.StringPointer("2022-01-02"),
				PurchaseTime: model.StringPointer("13:01"),
				Total:        model.StringPointer("5.12"),
				Items: []model.Item{
					{ShortDescription: model.StringPointer("Mountain Dew 12PK"), Price: model.StringPointer("5.12")},
				},
			},
			expectedPoints: 16,
		},
	}

	for _, tc := range testCases {
		receiptResp, _ := receiptService.Insert(tc.receipt)

		receipt, _ := receiptStore.Get(receiptResp.Id)
		assert.Equal(t, tc.expectedPoints, receipt.Points, fmt.Sprintf("Test %v Failed with use case %v", tc.id, tc.useCase))
	}
}
