package handler

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github/shivasaicharanruthala/backend-engineer-takehome/errors"
	"github/shivasaicharanruthala/backend-engineer-takehome/log"
	"github/shivasaicharanruthala/backend-engineer-takehome/model"
	"github/shivasaicharanruthala/backend-engineer-takehome/service"
	"io"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestHandlerGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	receiptService := service.NewMockReceipts(ctrl)
	logger, _ := log.NewCustomLogger("test.log")
	handler := New(logger, receiptService)

	testCases := []struct {
		id               int
		useCase          string
		path             string
		receiptID        string
		expectedResponse string
		statusCode       int
		mockCall         *gomock.Call
	}{
		{
			id: 1, useCase: "Negative case: request having query parameters",
			path:             "/v1/receipts/4a77ec9d-5334-43d0-a9e1-4fca8807bf8f/points?id=1",
			receiptID:        "4a77ec9d-5334-43d0-a9e1-4fca8807bf8f",
			expectedResponse: "query parameter is not accepted",
			statusCode:       400,
			mockCall:         nil,
		},
		{
			id: 2, useCase: "Negative case: invalid receipt id",
			path:             "/v1/receipts/4a77ec9d-5334-43d0-a9e1/points",
			receiptID:        "4a77ec9d-5334-43d0-a9e1",
			expectedResponse: "Incorrect value for parameter: id",
			statusCode:       400,
			mockCall:         nil,
		},
		{
			id: 3, useCase: "Negative case: receipt with given id not found",
			path:             "/v1/receipts/5a77ec9d-5334-43d0-a9e1-4fca8807bf8f/points",
			receiptID:        "5a77ec9d-5334-43d0-a9e1-4fca8807bf8f",
			expectedResponse: "No 'receipts' found for Id: '5a77ec9d-5334-43d0-a9e1-4fca8807bf8f'",
			statusCode:       404,
			mockCall: receiptService.EXPECT().Get("5a77ec9d-5334-43d0-a9e1-4fca8807bf8f").
				Return(nil, errors.NewEntityNotFound(errors.EntityNotFound{Entity: "receipts", ID: "5a77ec9d-5334-43d0-a9e1-4fca8807bf8f"})),
		},
		{
			id: 4, useCase: "Positive Case: fetch points",
			path:             "/v1/receipts/4a77ec9d-5334-43d0-a9e1-4fca8807bf8f/points",
			receiptID:        "4a77ec9d-5334-43d0-a9e1-4fca8807bf8f",
			expectedResponse: `{"points":10}`,
			statusCode:       200,
			mockCall: receiptService.EXPECT().Get("4a77ec9d-5334-43d0-a9e1-4fca8807bf8f").
				Return(&model.ReceiptGetResponse{Points: 10}, nil),
		},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", tc.path, nil)

		vars := map[string]string{
			"id": tc.receiptID,
		}
		r = mux.SetURLVars(r, vars)

		handler.Get(w, r)
		result := w.Result()
		resp, _ := io.ReadAll(result.Body)
		assert.Equal(t, tc.statusCode, result.StatusCode)
		if result.StatusCode >= 400 {
			regex, err := regexp.Compile(tc.expectedResponse)
			assert.NoError(t, err)

			assert.Regexp(t, regex, string(resp))
		} else {
			assert.Equal(t, tc.expectedResponse, string(resp))
		}
	}
}

func TestHandlerInsert(t *testing.T) {
	ctrl := gomock.NewController(t)
	receiptService := service.NewMockReceipts(ctrl)
	logger, _ := log.NewCustomLogger("test.log")
	handler := New(logger, receiptService)

	testCases := []struct {
		id               int
		useCase          string
		path             string
		reqBody          string
		statusCode       int
		expectedResponse string
		mockCall         *gomock.Call
	}{
		{
			id: 1, useCase: "Negative case: no body",
			path:             "/v1/receipts/process",
			reqBody:          ``,
			statusCode:       400,
			expectedResponse: "unexpected end of JSON input",
			mockCall:         nil,
		},
		{
			id: 2, useCase: "Negative case: invalid body",
			path:             "/v1/receipts/process",
			reqBody:          `{""}`,
			statusCode:       400,
			expectedResponse: "invalid character",
			mockCall:         nil,
		},
		{
			id: 3, useCase: "Negative case: empty body",
			path:             "/v1/receipts/process",
			reqBody:          `{}`,
			statusCode:       400,
			expectedResponse: "Parameter retailer is required for this request",
			mockCall: receiptService.EXPECT().Insert(&model.Receipt{}).
				Return(nil, errors.NewMissingParam(errors.MissingParam{Param: "retailer"}))},
		{
			id: 4, useCase: "Positive case: valid body",
			path:             "/v1/receipts/process",
			reqBody:          `{"retailer": "Target", "purchaseDate": "2022-01-01", "purchaseTime": "13:01", "items": [{"shortDescription": "Mountain Dew 12PK", "price": "5.00"}], "total": "5.00"}`,
			statusCode:       201,
			expectedResponse: "",
			mockCall: receiptService.EXPECT().Insert(&model.Receipt{
				Retailer:     model.StringPointer("Target"),
				PurchaseDate: model.StringPointer("2022-01-01"),
				PurchaseTime: model.StringPointer("13:01"),
				Total:        model.StringPointer("5.00"),
				Items: []model.Item{
					{ShortDescription: model.StringPointer("Mountain Dew 12PK"), Price: model.StringPointer("5.00")}}}).
				Return(&model.ReceiptPostResponse{Id: "4a77ec9d-5334-43d0-a9e1-4fca8807bf8f"}, nil),
		},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", tc.path, bytes.NewBuffer([]byte(tc.reqBody)))

		handler.Insert(w, r)
		result := w.Result()
		resp, _ := io.ReadAll(result.Body)

		assert.Equal(t, tc.statusCode, result.StatusCode)
		if result.StatusCode >= 400 {
			regex, err := regexp.Compile(tc.expectedResponse)
			assert.NoError(t, err)

			assert.Regexp(t, regex, string(resp))
		}
	}
}
