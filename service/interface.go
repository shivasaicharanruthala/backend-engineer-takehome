package service

import (
	"github/shivasaicharanruthala/backend-engineer-takehome/model"
)

type Receipts interface {
	Get(receiptID string) (*model.ReceiptGetResponse, error)
	Insert(receipt *model.Receipt) (*model.ReceiptPostResponse, error)
}
