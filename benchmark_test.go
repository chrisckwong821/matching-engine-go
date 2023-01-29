package main

import (
	"math/rand"
	"testing"
)

func benchmarkOrderbookLimitedRandomInsert(n int, b *testing.B) {
	book := NewOrderbook()

	// maximum number of levels is MaxLimitsNum
	limitslist := make([]float32, n)
	for i := range limitslist {
		limitslist[i] = rand.Float32()
	}

	// preallocate empty orders
	orders := make([]*Order, 0, b.N)
	for i := 0; i < b.N; i += 1 {
		orders = append(orders, &Order{})
	}

	// initialize the ringbuffer for each price
	for i := 0; i < b.N; i += 1 {
		price := limitslist[rand.Intn(len(limitslist))]
		// create a new order
		o := orders[i]
		o.Order.Quantity = 1
		o.Order.BidOrAsk = price < 0.5

		// add to the book
		book.Add(price, o)
	}
	// measure insertion time after all ringbuffer are initialized
	b.ResetTimer()
	for i := 0; i < b.N; i += 1 {
		price := limitslist[rand.Intn(len(limitslist))]
		// create a new order
		o := orders[i]
		o.Order.Quantity = 1
		o.Order.BidOrAsk = price < 0.5

		// add to the book
		book.Add(price, o)
	}

	//fmt.Printf("bid size %d, ask size %d\n", book.BLength(), book.ALength())
}

// average 300-400ns/ops or 3m op/s
func BenchmarkOrderbook5kLevelsRandomInsert(b *testing.B) {
	benchmarkOrderbookLimitedRandomInsert(10000, b)
}
