package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"regexp"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github/shivasaicharanruthala/backend-engineer-takehome/data"
	"github/shivasaicharanruthala/backend-engineer-takehome/handler"
	"github/shivasaicharanruthala/backend-engineer-takehome/log"
	"github/shivasaicharanruthala/backend-engineer-takehome/model"
	"github/shivasaicharanruthala/backend-engineer-takehome/service"
)

func TestIntegrations_Insert(t *testing.T) {
	server := httptest.NewServer(setUpRouter())
	defer server.Close()

	testCases := []struct {
		id               int
		useCase          string
		reqBody          string
		expectedResponse string
		statusCode       int
	}{
		{
			id: 1, useCase: "Positive case: request having query parameters",
			reqBody:          `{"retailer": "Target", "purchaseDate": "2022-01-01", "purchaseTime": "13:01", "items": [{"shortDescription": "Mountain Dew 12PK", "price": "5.00"}], "total": "5.00"}`,
			expectedResponse: "",
			statusCode:       201,
		},
		{
			id: 2, useCase: "Positive case: Duplicate data",
			reqBody:          `{"retailer": "Target", "purchaseDate": "2022-01-01", "purchaseTime": "13:01", "items": [{"shortDescription": "Mountain Dew 12PK", "price": "5.00"}], "total": "5.00"}`,
			expectedResponse: "",
			statusCode:       201,
		},
		{
			id: 3, useCase: "Negative case: missing price",
			reqBody:          `{"retailer": "Target", "purchaseDate": "2022-01-01", "purchaseTime": "13:01", "items": [{"shortDescription": "Mountain Dew 12PK"}], "total": "5.00"}`,
			expectedResponse: "Parameter price is required for this request",
			statusCode:       400,
		},
		{
			id: 4, useCase: "Negative case: missing shortDescription",
			reqBody:          `{"retailer": "Target", "purchaseDate": "2022-01-01", "purchaseTime": "13:01", "items": [{"price": "5.00"}], "total": "5.00"}`,
			expectedResponse: "Parameter shortDescription is required for this request",
			statusCode:       400,
		},
		{
			id: 5, useCase: "Negative case: missing total",
			reqBody:          `{"retailer": "Target", "purchaseDate": "2022-01-01", "purchaseTime": "13:01", "items": [{"shortDescription": "Mountain Dew 12PK", "price": "5.00"}]}`,
			expectedResponse: "Parameter total is required for this request",
			statusCode:       400,
		},
		{
			id: 6, useCase: "Negative case: missing items",
			reqBody:          `{"retailer": "Target", "purchaseDate": "2022-01-01", "purchaseTime": "13:01", "total": "5.00"}`,
			expectedResponse: "Parameter items is required for this request",
			statusCode:       400,
		},
		{
			id: 7, useCase: "Negative case: missing purchaseTime",
			reqBody:          `{"retailer": "Target", "purchaseDate": "2022-01-01", "items": [{"shortDescription": "Mountain Dew 12PK", "price": "5.00"}], "total": "5.00"}`,
			expectedResponse: "Parameter purchaseTime is required for this request",
			statusCode:       400,
		},
		{
			id: 8, useCase: "Negative case: missing purchaseDate",
			reqBody:          `{"retailer": "Target", "purchaseTime": "13:01", "items": [{"shortDescription": "Mountain Dew 12PK", "price": "5.00"}], "total": "5.00"}`,
			expectedResponse: "Parameter purchaseDate is required for this request",
			statusCode:       400,
		},
		{
			id: 9, useCase: "Negative case: missing retailer",
			reqBody:          `{"purchaseDate": "2022-01-01", "purchaseTime": "13:01", "items": [{"shortDescription": "Mountain Dew 12PK", "price": "5.00"}], "total": "5.00"}`,
			expectedResponse: "Parameter retailer is required for this request",
			statusCode:       400,
		},
	}

	for _, tc := range testCases {
		result, _ := http.Post(server.URL+"/v1/receipts/process", "application/json", bytes.NewBuffer([]byte(tc.reqBody)))
		resp, _ := io.ReadAll(result.Body)

		assert.Equal(t, tc.statusCode, result.StatusCode)
		if result.StatusCode >= 400 {
			regex, err := regexp.Compile(tc.expectedResponse)
			assert.NoError(t, err)

			assert.Regexp(t, regex, string(resp))
		}

	}
}

func TestIntegrations_Get(t *testing.T) {
	server := httptest.NewServer(setUpRouter())
	defer server.Close()

	// Insert test record
	var receiptPostResp model.ReceiptPostResponse

	insertRes, _ := http.Post(server.URL+"/v1/receipts/process", "application/json", bytes.NewBuffer([]byte(`{"retailer": "M&M Corner Market 45", "purchaseDate": "2022-01-01", "purchaseTime": "15:01", "items": [{"shortDescription": "Mountain Dew 12PK", "price": "5.00"}, {"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ", "price": "15.00"}, {"shortDescription": "Emils Cheese Pizza", "price": "5.00"}], "total": "25.00"}`)))
	insertResp, _ := io.ReadAll(insertRes.Body)
	_ = json.Unmarshal(insertResp, &receiptPostResp)

	testCases := []struct {
		id               int
		useCase          string
		receiptID        string
		expectedResponse string
		statusCode       int
	}{
		{
			id: 1, useCase: "Positive case: valid case",
			receiptID:        receiptPostResp.Id,
			expectedResponse: "{\"points\":116}",
			statusCode:       200,
		},
		{
			id: 2, useCase: "Negative case: receipt not found for given id",
			receiptID:        "4a77ec9d-5334-43d0-a9e1-4fca8807bf8f",
			expectedResponse: "No 'receipts' found for Id: '4a77ec9d-5334-43d0-a9e1-4fca8807bf8f'",
			statusCode:       404,
		},
		{
			id: 3, useCase: "Negative case: invalid receipt id",
			receiptID:        "1234",
			expectedResponse: "Incorrect value for parameter: id",
			statusCode:       400,
		},
	}

	for _, tc := range testCases {
		result, _ := http.Get(fmt.Sprintf("%v/v1/receipts/%v/points", server.URL, tc.receiptID))
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

func setUpRouter() *mux.Router {
	logger, _ := log.NewCustomLogger("test.log")
	receiptStore := data.New(logger)
	receiptSvc := service.New(logger, receiptStore)
	receiptHandler := handler.New(logger, receiptSvc)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1/receipts/{id}/points", receiptHandler.Get).Methods("GET")
	router.HandleFunc("/v1/receipts/process", receiptHandler.Insert).Methods("POST")

	return router
}

// TestLiveServer runs integration tests using Newman (Postman CLI).
// It starts the server, waits for it to initialize, and then executes the tests.
func TestLiveServer(t *testing.T) {
	// Start the server in a separate goroutine to allow it to run concurrently with the tests.
	go func() {
		main()
	}()

	// Allow some time for the server to start
	time.Sleep(1 * time.Second)

	// Define the command and arguments for running the integration tests with Newman.
	cmd := "npx"
	args := []string{
		"newman",
		"run",
		"./tests/integration_tests.json",
		"--reporters",
		"cli,junit",
		"--reporter-junit-export",
		"integration_report.xml",
		"--insecure",
	}

	//Create a command
	command := exec.Command(cmd, args...)
	command.Stdout = os.Stdout // Set the standard output to the console.
	command.Stderr = os.Stderr // Set the standard error to the console.

	// Run the command and check for errors.
	err := command.Run()
	if err != nil {
		t.Errorf("Expected All Integration Tests to pass but got error")
	}
}
