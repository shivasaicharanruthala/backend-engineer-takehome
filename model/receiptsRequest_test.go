package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github/shivasaicharanruthala/backend-engineer-takehome/errors"
)

func TestReceiptPayloadValidation(t *testing.T) {
	testCases := []struct {
		id        int
		useCase   string
		receipt   *Receipt
		expectErr error
	}{
		{
			id: 1, useCase: "Negative case: missing retailer",
			receipt: &Receipt{
				PurchaseDate: StringPointer("2022-02-01"),
				PurchaseTime: StringPointer("13:01"),
				Total:        StringPointer("5.00"),
				Items: []Item{
					{ShortDescription: StringPointer("Mountain Dew 12PK"), Price: StringPointer("5.00")},
				},
			},
			expectErr: errors.NewMissingParam(errors.MissingParam{Param: "retailer"}),
		},
		{
			id: 2, useCase: "Negative case: missing purchaseDate",
			receipt: &Receipt{
				Retailer:     StringPointer("Target"),
				PurchaseTime: StringPointer("13:01"),
				Total:        StringPointer("5.00"),
				Items: []Item{
					{ShortDescription: StringPointer("Mountain Dew 12PK"), Price: StringPointer("5.00")},
				},
			},
			expectErr: errors.NewMissingParam(errors.MissingParam{Param: "purchaseDate"}),
		},
		{
			id: 3, useCase: "Negative case: missing purchaseTime",
			receipt: &Receipt{
				Retailer:     StringPointer("Target"),
				PurchaseDate: StringPointer("2022-02-01"),
				Total:        StringPointer("5.00"),
				Items: []Item{
					{ShortDescription: StringPointer("Mountain Dew 12PK"), Price: StringPointer("5.00")},
				},
			},
			expectErr: errors.NewMissingParam(errors.MissingParam{Param: "purchaseTime"}),
		},
		{
			id: 4, useCase: "Negative case: missing shortDescription",
			receipt: &Receipt{
				Retailer:     StringPointer("Target"),
				PurchaseDate: StringPointer("2022-02-01"),
				PurchaseTime: StringPointer("13:01"),
				Total:        StringPointer("5.00"),
				Items:        []Item{{Price: StringPointer("5.00")}},
			},
			expectErr: errors.NewMissingParam(errors.MissingParam{Param: "shortDescription"})},
		{
			id: 5, useCase: "Negative case: missing price",
			receipt: &Receipt{
				Retailer:     StringPointer("Target"),
				PurchaseDate: StringPointer("2022-02-01"),
				PurchaseTime: StringPointer("13:01"),
				Total:        StringPointer("5.00"),
				Items:        []Item{{ShortDescription: StringPointer("Mountain Dew 12PK")}},
			},
			expectErr: errors.NewMissingParam(errors.MissingParam{Param: "price"}),
		},
		{
			id: 6, useCase: "Negative case: missing items",
			receipt: &Receipt{
				Retailer:     StringPointer("Target"),
				PurchaseDate: StringPointer("2022-02-01"),
				PurchaseTime: StringPointer("13:01"),
				Total:        StringPointer("5.00"),
			},
			expectErr: errors.NewMissingParam(errors.MissingParam{Param: "items"}),
		},
		{
			id: 7, useCase: "Negative case: missing total",
			receipt: &Receipt{
				Retailer:     StringPointer("Target"),
				PurchaseDate: StringPointer("2022-02-01"),
				PurchaseTime: StringPointer("13:01"),
				Items: []Item{
					{ShortDescription: StringPointer("Mountain Dew 12PK"), Price: StringPointer("5.00")},
				},
			},
			expectErr: errors.NewMissingParam(errors.MissingParam{Param: "total"}),
		},
		{
			id: 8, useCase: "Negative case: missing receipt",
			receipt:   nil,
			expectErr: errors.NewMissingParam(errors.MissingParam{Param: "receipt"}),
		},
		{
			id: 9, useCase: "Positive case: valid receipt",
			receipt: &Receipt{
				Retailer:     StringPointer("Target"),
				PurchaseDate: StringPointer("2022-02-01"),
				PurchaseTime: StringPointer("13:01"),
				Total:        StringPointer("5.00"),
				Items: []Item{
					{ShortDescription: StringPointer("Mountain Dew 12PK"), Price: StringPointer("5.00")},
				},
			},
			expectErr: nil,
		},
	}

	for _, tc := range testCases {
		err := tc.receipt.PayloadValidation()
		if err != nil {
			assert.Equal(t, tc.expectErr.Error(), err.Error(), fmt.Sprintf("Test %v Failed with use case %v", tc.id, tc.useCase))
		} else {
			assert.Equal(t, tc.expectErr, err, fmt.Sprintf("Test %v Failed with use case %v", tc.id, tc.useCase))
		}
	}
}
func TestItemPayloadValidation(t *testing.T) {
	testCases := []struct {
		id        int
		useCase   string
		item      *Item
		expectErr error
	}{
		{
			id: 1, useCase: "Negative case: missing short description",
			item:      &Item{Price: StringPointer("5.00")},
			expectErr: errors.NewMissingParam(errors.MissingParam{Param: "shortDescription"}),
		},
		{
			id: 2, useCase: "Negative case: missing price",
			item:      &Item{ShortDescription: StringPointer("Mountain Dew 12PK")},
			expectErr: errors.NewMissingParam(errors.MissingParam{Param: "price"}),
		},
		{
			id: 3, useCase: "Positive case: valid item",
			item: &Item{
				ShortDescription: StringPointer("Mountain Dew 12PK"),
				Price:            StringPointer("5.00"),
			},
			expectErr: nil,
		},
	}

	for _, tc := range testCases {
		err := tc.item.PayloadValidation()
		if err != nil {
			assert.Equal(t, tc.expectErr.Error(), err.Error(), fmt.Sprintf("Test %v Failed with use case %v", tc.id, tc.useCase))
		} else {
			assert.Equal(t, tc.expectErr, err, fmt.Sprintf("Test %v Failed with use case %v", tc.id, tc.useCase))
		}
	}
}
