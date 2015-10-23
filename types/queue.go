package types

import (
	"errors"
	"sync"
)

var ErrEmpty = errors.New("empty queue")

type Queue struct {
	*sync.RWMutex
	buffer [][]byte
}

func NewQueue() *Queue {
	return &Queue{RWMutex: new(sync.RWMutex)}
}

func (q *Queue) Put(payload []byte) {
	q.Lock()
	defer q.Unlock()
	q.buffer = append(q.buffer, payload)
}

func (q *Queue) Get() ([]byte, error) {
	q.Lock()
	defer q.Unlock()
	if len(q.buffer) == 0 {
		return nil, ErrEmpty
	}
	item := q.buffer[0]
	q.buffer = q.buffer[1:]
	return item, nil
}

func (q *Queue) Size() int {
	q.RLock()
	defer q.RUnlock()
	return len(q.buffer)
}

func (q *Queue) Empty() bool {
	q.RLock()
	defer q.RUnlock()
	return len(q.buffer) == 0
}
