package batching

import (
	"fmt"
	"testing"
	"time"
)

func TestBatcher(t *testing.T) {
	b := NewBatcher(&toBatch{}, 30, bb, time.Second)
	for i := 0; i < 30; i++ {
		time.Sleep(time.Millisecond * 200)
		b.Add(&toBatch{2})
	}
}

type toBatch struct {
	n int
}

func (tb *toBatch) Merge(other BatchSubject) {
	tb.n += other.(*toBatch).n
}

func (tb *toBatch) Size() int {
	return tb.n
}

func bb(tb BatchSubject) {
	fmt.Println("batching", tb)
}
