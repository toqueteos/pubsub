package types

import (
	"sync/atomic"
)

type ChanQueue struct {
	buffer chan []byte
	count  int32
}

func NewChanQueue(size int) *ChanQueue {
	return &ChanQueue{
		buffer: make(chan []byte, size),
	}
}

func (q *ChanQueue) Put(payload []byte) {
	q.buffer <- payload
	atomic.AddInt32(&q.count, 1)
}

func (q *ChanQueue) Get() ([]byte, error) {
	item := <-q.buffer
	atomic.AddInt32(&q.count, -1)
	return item, nil
}

func (q *ChanQueue) Size() int {
	count := atomic.LoadInt32(&q.count)
	return int(count)
}

func (q *ChanQueue) Empty() bool {
	return q.Size() == 0
}
