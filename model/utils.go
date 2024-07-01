package model

import (
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github/shivasaicharanruthala/backend-engineer-takehome/errors"
)

// IsValidUUID checks if a given string is a valid UUID format.
// It returns true if the UUID string is valid, otherwise false.
func IsValidUUID(uuidStr string) bool {
	_, err := uuid.Parse(uuidStr)
	return err == nil
}

// StringPointer takes a string and returns a pointer to that string.
func StringPointer(s string) *string {
	return &s
}

// CountAlphanumericCharacters Rule-1: One point for every alphanumeric character in the retailer name.
// CountAlphanumericCharacters counts the number of alphanumeric characters in a string.
// It iterates through each character in the string and counts alphabetic characters (letters).
func (receipt *Receipt) CountAlphanumericCharacters() {
	for _, char := range *receipt.Retailer {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			receipt.Points += 1
		}
	}
}

// FiftyPointRule is Rule-2: 50 points if the total is a round dollar amount with no cents.
func (receipt *Receipt) FiftyPointRule() {
	total := *receipt.Total
	centsOfTotal := total[len(total)-2:]
	if centsOfTotal == "00" {
		receipt.Points += 50
	}
}

// TwentyFivePointRule is Rule-3: 25 points if the total is a multiple of 0.25.
func (receipt *Receipt) TwentyFivePointRule() {
	convertedTotal, _ := strconv.ParseFloat(*receipt.Total, 64)
	if math.Mod(convertedTotal, 0.25) == 0 {
		receipt.Points += 25
	}
}

// FivePointRule is Rule-4: 5 points for every two items on the receipt.
func (receipt *Receipt) FivePointRule() {
	receipt.Points += 5 * (len(receipt.Items) / 2)
}

// CountTrimmedItemDescriptionPoints is Rule-5: If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and
// round up to the nearest integer. The result is the number of points earned.
// CountTrimmedItemDescriptionPoints calculates points based on trimmed item descriptions.
// It iterates through a list of items, trims each item's short description,
// calculates points based on specific conditions, and accumulates the total points earned.
func (receipt *Receipt) CountTrimmedItemDescriptionPoints() {
	var count int
	for _, item := range receipt.Items {
		trimmedItemName := strings.Trim(*item.ShortDescription, " ") // Trim leading and trailing spaces from the item's short description.
		trimmedNameLength := len(trimmedItemName)
		if math.Mod(float64(trimmedNameLength), 3) == 0 { // Check if the length of the trimmed name is divisible by 3.
			itemPrice, _ := strconv.ParseFloat(*item.Price, 64) // Parse item price from string to float64.
			pointsEarned := int(math.Ceil(itemPrice * 0.2))
			count += pointsEarned
		}
	}

	receipt.Points += count
}

// SixPointRule Rule-6: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
func (receipt *Receipt) SixPointRule() error {
	parsedPurchaseDate, err := time.Parse("2006-01-02", *receipt.PurchaseDate) // parse PurchaseDate in "YYYY-MM-DD" format
	if err != nil {
		return errors.NewInvalidParam(errors.InvalidParam{Param: "purchaseDate"})
	}

	purchaseDay := parsedPurchaseDate.Day()
	if purchaseDay%2 != 0 {
		receipt.Points += 6
	}

	return nil
}

// TenPointRule Rule-7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
func (receipt *Receipt) TenPointRule() error {
	parsedPurchaseTime, err := time.Parse("15:04", *receipt.PurchaseTime) // parse PurchaseTime in 24hrs format
	if err != nil {
		return errors.NewInvalidParam(errors.InvalidParam{Param: "purchaseTime"})
	}

	// Create time objects for 2 PM and 4 PM
	twoPM, _ := time.Parse("15:04", "14:00")
	fourPM, _ := time.Parse("15:04", "16:00")

	// Check if the parsed time is between 2 PM and 4 PM
	if parsedPurchaseTime.After(twoPM) && parsedPurchaseTime.Before(fourPM) {
		receipt.Points += 10
	}

	return nil
}

// CalculateTotalReceiptPoints calculates total points for a receipt based on various criteria.
// It computes points from retailer name, total amount, item descriptions, purchase date,
// purchase time, and specific time conditions.
// It updates the Points field of the Receipt struct and returns an error if there are parsing issues.
func (receipt *Receipt) CalculateTotalReceiptPoints() error {
	// Rule-1: One point for every alphanumeric character in the retailer name.
	receipt.CountAlphanumericCharacters()

	// Rule-2: 50 points if the total is a round dollar amount with no cents.
	receipt.FiftyPointRule()

	// Rule-3: 25 points if the total is a multiple of 0.25.
	receipt.TwentyFivePointRule()

	// Rule-4: 5 points for every two items on the receipt.
	receipt.FivePointRule()

	// Rule-5: if the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and
	// round up to the nearest integer. The result is the number of points earned.
	receipt.CountTrimmedItemDescriptionPoints()

	// Rule-6: 6 points if the day in the purchase date is odd.
	if err := receipt.SixPointRule(); err != nil {
		return err
	}

	// Rule-7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	if err := receipt.TenPointRule(); err != nil {
		return err
	}

	return nil
}
