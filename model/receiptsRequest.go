package model

import "github/shivasaicharanruthala/backend-engineer-takehome/errors"

// Item represents an item in a receipt
type Item struct {
	ShortDescription *string `json:"shortDescription"`
	Price            *string `json:"price"`
}

// Receipt represents a receipt with its details including items purchased.
type Receipt struct {
	Id           string
	Retailer     *string `json:"retailer"`
	PurchaseDate *string `json:"purchaseDate"`
	PurchaseTime *string `json:"purchaseTime"`
	Items        []Item  `json:"items"`
	Total        *string `json:"total"`
	Points       int
}

// PayloadValidation performs validation on the receipt's payload fields.
// It checks for required fields like retailer, purchase date, purchase time, and items.
// It also delegates validation of each item in the receipt.
func (receipt *Receipt) PayloadValidation() error {
	if receipt == nil {
		return errors.NewMissingParam(errors.MissingParam{Param: "receipt"})
	}

	if receipt.Retailer == nil {
		return errors.NewMissingParam(errors.MissingParam{Param: "retailer"})
	}

	if receipt.PurchaseDate == nil {
		return errors.NewMissingParam(errors.MissingParam{Param: "purchaseDate"})
	}

	if receipt.PurchaseTime == nil {
		return errors.NewMissingParam(errors.MissingParam{Param: "purchaseTime"})
	}

	if receipt.Total == nil {
		return errors.NewMissingParam(errors.MissingParam{Param: "total"})
	}

	if receipt.Items == nil {
		return errors.NewMissingParam(errors.MissingParam{Param: "items"})
	}

	// Validate each item in the receipt.
	for _, item := range receipt.Items {
		if err := item.PayloadValidation(); err != nil {
			return err
		}
	}

	return nil
}

// PayloadValidation performs validation on the item's payload fields.
// It checks for required fields like short description and price.
func (i *Item) PayloadValidation() error {
	if i.ShortDescription == nil {
		return errors.NewMissingParam(errors.MissingParam{Param: "shortDescription"})
	}

	if i.Price == nil {
		return errors.NewMissingParam(errors.MissingParam{Param: "price"})
	}

	return nil
}
