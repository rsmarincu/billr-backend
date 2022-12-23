package main

import (
	"github.com/rsmarincu/billr/api"
)

type BillrService struct {
	httpHandler api.Handler
}

func NewBillrService(handler api.Handler) *BillrService {
	return &BillrService{httpHandler: handler}
}

func (b *BillrService) Init() {
	b.httpHandler.SetupRouting()
}
