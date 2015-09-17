package types

import "errors"

var ErrEmpty = errors.New("empty queue")

type Queue struct {
	buffer [][]byte
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Put(payload []byte) {
	q.buffer = append(q.buffer, payload)
}

func (q *Queue) Get() ([]byte, error) {
	if len(q.buffer) == 0 {
		return nil, ErrEmpty
	}
	item := q.buffer[0]
	q.buffer = q.buffer[1:]
	return item, nil
}

func (q *Queue) Size() int   { return len(q.buffer) }
func (q *Queue) Empty() bool { return len(q.buffer) == 0 }
