package types

import (
	"fmt"
	"sync"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type QueueSuite struct{}

var _ = Suite(&QueueSuite{})

func (s *QueueSuite) TestBasicQueue(c *C) {
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
	c.Assert(result, DeepEquals, A)
	c.Assert(err, IsNil)
}

func (s *QueueSuite) TestConcurrentQueue(c *C) {
	var q = NewQueue()
	c.Assert(q.Size(), Equals, 0)
	c.Assert(q.Empty(), Equals, true)

	var n = 10000
	for i := 0; i < n; i++ {
		q.Put([]byte(fmt.Sprint(i)))
	}
	c.Assert(q.Empty(), Equals, false)
	c.Assert(q.Size(), Equals, n)

	wg := new(sync.WaitGroup)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var result, err = q.Get()
			c.Assert(result, NotNil)
			c.Assert(err, IsNil)
		}()
	}
	wg.Wait()
	c.Assert(q.Empty(), Equals, true)
}
