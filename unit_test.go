package main

import (
	"fmt"
	"testing"
)

func NewMockBuyOrder() Order {
	price := float32(10)
	quantity := float32(100)
	accountId := uint32(2)
	sequenceId := uint32(1)
	BidOrAsk := true // true = bid
	orderType := LIMIT
	ni := NewIncomingOrder(price, quantity, BidOrAsk, orderType, accountId)
	no := NewOrder(ni, sequenceId)
	return no
}

func NewMockSellOrder() Order {
	price := float32(10)
	quantity := float32(100)
	accountId := uint32(2)
	sequenceId := uint32(1)
	BidOrAsk := false // true = bid
	orderType := LIMIT
	ni := NewIncomingOrder(price, quantity, BidOrAsk, orderType, accountId)
	no := NewOrder(ni, sequenceId)
	return no
}

func NewCustomOrder(price float32, quantity float32, BidOrAsk bool, orderType OrderType, accountId uint32, sequenceId uint32) Order {
	ni := NewIncomingOrder(price, quantity, BidOrAsk, orderType, accountId)
	no := NewOrder(ni, sequenceId)
	return no
}

// OrderBook
func TestOrderbookEmpty(t *testing.T) {
	b := NewOrderbook()
	if b.BLength() != 0 || b.ALength() != 0 {
		t.Errorf("orderbook should be empty")
	}
}

func TestPlaceOrder(t *testing.T) {
	b := NewOrderbook()
	o := NewMockBuyOrder()
	oo := &o
	b.Execute(oo)
}

func TestGetBestBid(t *testing.T) {
	b := NewOrderbook()
	o := NewMockBuyOrder()
	oo := &o
	b.Execute(oo)
	if b.GetBestBid() != o.Order.Price {
		t.Errorf("best bid should be order price")
	}

}

func TestGetVolumeAtBidLimit(t *testing.T) {
	b := NewOrderbook()
	o := NewMockBuyOrder()
	oo := &o
	b.Execute(oo)
	if b.GetVolumeAtBidLimit(o.Order.Price) != o.Order.Quantity {
		t.Errorf("best offer volume should be order quantity")
	}
}

func TestGetBestOffer(t *testing.T) {
	b := NewOrderbook()
	o := NewMockSellOrder()
	oo := &o
	b.Execute(oo)
	if b.GetBestOffer() != o.Order.Price {
		t.Errorf("best offer should be order price")
	}
}

func TestGetVolumeAtAskLimit(t *testing.T) {
	b := NewOrderbook()
	o := NewMockSellOrder()
	oo := &o
	b.Execute(oo)
	if b.GetVolumeAtAskLimit(o.Order.Price) != o.Order.Quantity {
		t.Errorf("best offer volume should be order quantity")
	}
}

func TestMatchingEqualVolume(t *testing.T) {
	b := NewOrderbook()
	o := NewMockBuyOrder()
	oo := &o
	b.Execute(oo)
	// fmt.Println("FIRST EXECUTION :  ", exc)
	// fmt.Println("Best Bid :  ", b.GetBestBid())
	// fmt.Println("Best Bid Volume :  ", b.GetVolumeAtBidLimit(o.Order.Price))
	quantity := float32(o.Order.Quantity)
	price := float32(o.Order.Price)
	accountId := uint32(o.Order.AccountId)
	sequenceId := uint32(o.SequenceId + 1)
	BidOrAsk := false // true = bid
	orderType := LIMIT
	sellOrder := NewCustomOrder(price, quantity, BidOrAsk, orderType, accountId, sequenceId)
	s := &sellOrder
	b.Execute(s)

	if b.BLength() != 0 || b.ALength() != 0 {
		fmt.Println(b.BLength(), b.ALength())
		t.Errorf("orderbook should be empty")
	}
	// fmt.Println("SECOND EXECUTION : ", excc)

}

func TestMatchingHalfVolume(t *testing.T) {
	b := NewOrderbook()
	o := NewMockBuyOrder()
	oo := &o
	b.Execute(oo)
	// fmt.Println("FIRST EXECUTION :  ", exc)
	// fmt.Println("Best Bid :  ", b.GetBestBid())
	// fmt.Println("Best Bid Volume :  ", b.GetVolumeAtBidLimit(o.Order.Price))
	quantity := float32(o.Order.Quantity) / 2
	price := float32(o.Order.Price)
	accountId := uint32(o.Order.AccountId)
	sequenceId := uint32(o.SequenceId + 1)
	BidOrAsk := false // true = bid
	orderType := LIMIT
	sellOrder := NewCustomOrder(price, quantity, BidOrAsk, orderType, accountId, sequenceId)
	s := &sellOrder
	b.Execute(s)

	if b.GetVolumeAtBidLimit(o.Order.Price) != (o.Order.Quantity / 2) {
		fmt.Println(b.GetVolumeAtBidLimit(o.Order.Price))
		t.Errorf("best offer volume should be order quantity")
	}
	// fmt.Println("SECOND EXECUTION : ", excc)

}

func TestMatchingDoubleVolume(t *testing.T) {
	b := NewOrderbook()
	o := NewMockBuyOrder()
	oo := &o
	b.Execute(oo)
	// fmt.Println("FIRST EXECUTION :  ", exc)
	// fmt.Println("Best Bid :  ", b.GetBestBid())
	// fmt.Println("Best Bid Volume :  ", b.GetVolumeAtBidLimit(o.Order.Price))
	quantity := float32(o.Order.Quantity) * 2
	price := float32(o.Order.Price)
	accountId := uint32(o.Order.AccountId)
	sequenceId := uint32(o.SequenceId + 1)
	BidOrAsk := false // true = bid
	orderType := LIMIT
	sellOrder := NewCustomOrder(price, quantity, BidOrAsk, orderType, accountId, sequenceId)
	s := &sellOrder
	b.Execute(s)

	if b.GetVolumeAtAskLimit(o.Order.Price) != o.Order.Quantity-o.ExecutedQuantity {
		fmt.Println(b.GetVolumeAtAskLimit(o.Order.Price))
		t.Errorf("best offer volume should be order quantity")
	}
	// fmt.Println("SECOND EXECUTION : ", excc)

}

func TestOrderbookAddMultiple(t *testing.T) {
	b := NewOrderbook()
	for i := 0; i < 100; i += 1 {
		bid := NewMockBuyOrder()
		bid.Order.Price = float32(i)
		b.Execute(&bid)
	}

	for i := 100; i < 200; i += 1 {
		ask := NewMockSellOrder()
		ask.Order.Price = float32(i)
		b.Execute(&ask)
	}

	if b.BLength() != 100 {
		t.Errorf("book should have 100 bids")
	}
	if b.ALength() != 100 {
		t.Errorf("book should have 100 asks")
	}

}

func TestOrderqueueResizing(t *testing.T) {
	b := NewOrderbook()
	var bid, ask Order
	no_Order := RINGBUF_INI_SIZE + 1
	for i := 0; i < no_Order; i += 1 {
		bid = NewMockBuyOrder()
		b.Execute(&bid)
	}

	for i := 0; i < no_Order; i += 1 {
		ask = NewMockSellOrder()
		// to avoid matching
		ask.Order.Price = bid.Order.Price + 1
		b.Execute(&ask)
	}

	if b.BLength() != 1 {
		t.Errorf("book should have 1 bids")
	}

	if b.bidLimitsCache[bid.Order.Price].ringbuffer.Len() != no_Order {
		t.Errorf("ringbuffer has not resized")
	}

}
