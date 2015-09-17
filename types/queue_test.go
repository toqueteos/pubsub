package types

import . "gopkg.in/check.v1"

type QueueSuite struct{}

var _ = Suite(&QueueSuite{})

func (s *QueueSuite) TestQueue(c *C) {
	var q = NewQueue()
	var A = []byte("A")
	c.Assert(q.Size(), Equals, 0)
	c.Assert(q.Empty(), Equals, true)

	q.Put(A)
	q.Put(A)
	c.Assert(q.Size(), Equals, 2)
	c.Assert(q.Empty(), Equals, false)

	var result, err = q.Get()
	c.Assert(q.Size(), Equals, 1)
	c.Assert(q.Empty(), Equals, false)
	c.Assert(result, Equals, A)
	c.Assert(err, NotNil)
}
