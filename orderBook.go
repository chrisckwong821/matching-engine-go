package main

import (
	"sync"
)

// maximum limits per orderbook side to pre-allocate memory
const MaxLimitsNum int = 10000

type Orderbook struct {
	Bids *redBlackBST
	Asks *redBlackBST

	bidLimitsCache map[float32]*OrdersQueue
	askLimitsCache map[float32]*OrdersQueue
	pool           *sync.Pool
}

func NewOrderbook() Orderbook {
	bids := NewRedBlackBST()
	asks := NewRedBlackBST()
	return Orderbook{
		Bids: &bids,
		Asks: &asks,

		bidLimitsCache: make(map[float32]*OrdersQueue, MaxLimitsNum),
		askLimitsCache: make(map[float32]*OrdersQueue, MaxLimitsNum),
		pool: &sync.Pool{
			New: func() interface{} {
				orderqueue := NewOrdersQueue(0.0, Order{}.size())
				return &orderqueue
			},
		},
	}
}

// entry point for order to either be queued or matched
func (this *Orderbook) Execute(o *Order) float32 {
	if o.Order.BidOrAsk {
		return this.ExecuteBid(o)
	} else {
		return this.ExecuteAsk(o)
	}
}

// entry point for bid
func (this *Orderbook) ExecuteBid(o *Order) float32 {
	// no best ask
	if this.ALength() == 0 {
		this.Add(o.Order.Price, o)
		return o.ExecutedQuantity
	}
	best_ask := this.GetBestOffer()
	leftQuantity := o.Order.Quantity - o.ExecutedQuantity
	// order fully filled, exit
	if leftQuantity == 0 {
		return o.ExecutedQuantity
	}
	// if order can be matched
	if o.Order.Price >= best_ask || o.Order.OrderType == MARKET {
		v := this.GetVolumeAtAskLimit(best_ask)
		// if the best ask can be swept
		if leftQuantity >= v {
			o.ExecutedQuantity += v
			// execute each order in the ringbuffer
			this.askLimitsCache[best_ask].Execute(v)
			// remove the ask from the cache&BST, return the ringbuffer to the pool
			this.DeleteAskLimit(best_ask)
			// recursive call on the next best ask
			this.ExecuteBid(o)
		} else
		// if the order would be fully filled in the current ask
		{
			this.askLimitsCache[best_ask].Execute(leftQuantity)
			o.ExecutedQuantity = o.Order.Quantity
			return o.ExecutedQuantity
		}

	} else
	// if order can NOT be matched
	{
		this.Add(o.Order.Price, o)
	}

	return o.ExecutedQuantity
}

// entry point for ask
func (this *Orderbook) ExecuteAsk(o *Order) float32 {
	// no best bid
	if this.BLength() == 0 {
		this.Add(o.Order.Price, o)
		return o.ExecutedQuantity
	}
	best_bid := this.GetBestBid()
	leftQuantity := o.Order.Quantity - o.ExecutedQuantity
	// order fully filled, exit
	if leftQuantity == 0 {
		return o.ExecutedQuantity
	}
	// if order can be matched
	if o.Order.Price <= best_bid || o.Order.OrderType == MARKET {
		// if the best bid can be swept
		v := this.GetVolumeAtBidLimit(best_bid)
		if leftQuantity >= v {
			o.ExecutedQuantity += v
			// execute each order in the ringbuffer
			this.bidLimitsCache[best_bid].Execute(v)
			// remove the bid from the cache&BST, return the ringbuffer to the pool
			this.DeleteBidLimit(best_bid)
			// recursive call
			this.ExecuteAsk(o)
		} else {
			// if the order would be fully filled in the current ask
			this.bidLimitsCache[best_bid].Execute(leftQuantity)
			o.ExecutedQuantity = o.Order.Quantity
			return o.ExecutedQuantity
		}
	} else
	// if order can NOT be matched
	{
		this.Add(o.Order.Price, o)
	}
	return o.ExecutedQuantity
}

func (this *Orderbook) Add(price float32, o *Order) {
	var orderqueue *OrdersQueue
	if o.Order.BidOrAsk {
		orderqueue = this.bidLimitsCache[price]
	} else {
		orderqueue = this.askLimitsCache[price]
	}

	if orderqueue == nil {
		// getting a new limit from pool
		orderqueue = this.pool.Get().(*OrdersQueue)
		orderqueue.price = price

		// insert into the corresponding BST and cache
		if o.Order.BidOrAsk {
			this.Bids.Put(price, orderqueue)
			this.bidLimitsCache[price] = orderqueue
		} else {
			this.Asks.Put(price, orderqueue)
			this.askLimitsCache[price] = orderqueue
		}
	}

	// add order to the limit
	orderqueue.PlaceOrder(o)
}

func (this *Orderbook) DeleteBidLimit(price float32) {
	limit := this.bidLimitsCache[price]
	if limit == nil {
		return
	}

	this.deleteLimit(price, true)
	delete(this.bidLimitsCache, price)

	// put limit back to the pool
	limit.Clear()
	this.pool.Put(limit)

}

func (this *Orderbook) DeleteAskLimit(price float32) {
	orderqueue := this.askLimitsCache[price]
	if orderqueue == nil {
		return
	}

	this.deleteLimit(price, false)
	delete(this.askLimitsCache, price)

	// put limit back to the pool
	orderqueue.Clear()
	this.pool.Put(orderqueue)
}

func (this *Orderbook) deleteLimit(price float32, bidOrAsk bool) {
	if bidOrAsk {
		this.Bids.Delete(price)
	} else {
		this.Asks.Delete(price)
	}
}

func (this *Orderbook) GetVolumeAtBidLimit(price float32) float32 {
	orderqueue := this.bidLimitsCache[price]
	if orderqueue == nil {
		return 0
	}
	return orderqueue.TotalVolume()
}

func (this *Orderbook) GetVolumeAtAskLimit(price float32) float32 {
	orderqueue := this.askLimitsCache[price]
	if orderqueue == nil {
		return 0
	}
	return orderqueue.TotalVolume()
}

func (this *Orderbook) GetBestBid() float32 {
	return this.Bids.Max()
}

func (this *Orderbook) GetBestOffer() float32 {
	return this.Asks.Min()
}

func (this *Orderbook) BLength() int {
	return len(this.bidLimitsCache)
}

func (this *Orderbook) ALength() int {
	return len(this.askLimitsCache)
}
