package main

// default order per tick before resizing the ringbuffer
const RINGBUF_INI_SIZE = 1 << 13 // 8192 order

// r is used for matching, while w is used for placing order
//
//	type Deque[T any] struct {
//		buf    []T
//		head   int
//		tail   int
//		count  int
//		minCap int
//	}
type OrdersQueue struct {
	price         float32
	totalVolume   float32
	ringbuffer    *Deque[Order]
	orderByteSize int
}

func (this *OrdersQueue) Price() float32 {
	return this.price
}

func (this *OrdersQueue) TotalVolume() float32 {
	return this.totalVolume
}

func NewOrdersQueue(price float32, orderByteSize int) OrdersQueue {
	var r = New[Order](RINGBUF_INI_SIZE)
	return OrdersQueue{price, 0, r, orderByteSize}
}

func (this *OrdersQueue) Size() int {
	return this.ringbuffer.Len()
}

func (this *OrdersQueue) IsEmpty() bool {
	return this.ringbuffer.Len() == 0
}

func (this *OrdersQueue) PlaceOrder(o *Order) {
	q := o.Order.Quantity - o.ExecutedQuantity
	this.totalVolume += float32(q)
	// if the oldest order is not matched and the ringbuffer filled up, the ringbuffer would resize.
	this.ringbuffer.PushBack(*o)
}

// the queue doesnt care price level.
func (this *OrdersQueue) Execute(quantity float32) float32 {
	if this.ringbuffer.Len() == 0 {
		return 0
	}

	o := quantity
	var order Order
	var q float32 = 0
	// execute logic
	for quantity >= q {
		order = this.ringbuffer.Front()
		q = order.Order.Quantity
		//fully filled
		if quantity >= q {
			quantity -= q
			// move the read pointer by flushing
			// @TODO check order status to implement cancel order
			this.ringbuffer.PopFront()
			// partial filled
			this.totalVolume -= q
			order.ExecutedQuantity += q
		} else {
			// write the updated quantity back to the buf
			// a better way can be a cache quantity of top order in the orderqueue level
			order.ExecutedQuantity = quantity
			this.ringbuffer.PopFront()
			this.ringbuffer.PushFront(order)
			this.totalVolume -= quantity
		}

	}
	return o
}

func (this *OrdersQueue) Clear() {
	this.ringbuffer.Clear()
}
