package model

// ReceiptPostResponse represents the response structure after successfully posting a receipt.
type ReceiptPostResponse struct {
	Id string `json:"id"`
}

// ReceiptGetResponse represents the response structure when retrieving receipt details.
type ReceiptGetResponse struct {
	Points int `json:"points"`
}
