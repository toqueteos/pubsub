package types

import (
	"reflect"
	"testing"
)

func TestChanQueue(t *testing.T) {
	var A = []byte("A")

	q := NewChanQueue(32)
	size := q.Size()
	if size != 0 || !q.Empty() {
		t.Errorf("expected 0 ChanQueue.Size, got %d\n", size)
	}

	q.Put(A)
	q.Put(A)
	size = q.Size()
	if size != 2 {
		t.Errorf("expected size 2, got %d\n", size)
	}

	item, err := q.Get()
	size = q.Size()
	if size != 1 {
		t.Errorf("expected size 1, got %d\n", size)
	}

	if !reflect.DeepEqual(A, item) {
		t.Errorf("expected %x, got %x\n", A, item)
	}
	if err != nil {
		t.Errorf("expected nil err, got %s\n", err)
	}
}

func BenchmarkChanQueue(b *testing.B) {
	b.ReportAllocs()

	q := NewChanQueue(32)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q.Put([]byte{0xde, 0xad, 0xbe, 0xef, 0xc0, 0xff, 0xee})
		q.Get()
	}
}
