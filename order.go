package main

import (
	"unsafe"
)

type OrderType uint8

const (
	MARKET OrderType = iota
	LIMIT
)

// single order in the orderqueue
type Order struct {
	Order IncomingOrder // 16
	//each order is given a unique increasing sequenceId for deterministically handling
	SequenceId       uint32  // 4
	ExecutedQuantity float32 // 4
}

type IncomingOrder struct {
	Price    float32 //4
	Quantity float32 //4
	BidOrAsk bool    // 1
	// market / limit
	OrderType OrderType // 1
	AccountId uint32    // 4

}

func NewIncomingOrder(price float32, quantity float32, BidOrAsk bool, orderType OrderType, accountId uint32) IncomingOrder {
	return IncomingOrder{price, quantity, BidOrAsk, orderType, accountId}
}

func NewOrder(incomingOrder IncomingOrder, sequenceId uint32) Order {
	return Order{incomingOrder, sequenceId, 0}
}

func (o Order) size() int {
	return int(unsafe.Sizeof(o))
}

func (o IncomingOrder) size() int {
	return int(unsafe.Sizeof(o))
}
