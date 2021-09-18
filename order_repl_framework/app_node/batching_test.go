package app_node

import (
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	"testing"
	"time"
)

func TestBatching(t *testing.T) {
	as := &taaa{}
	b := newPrepBatch(as.onB, 128, time.Second*10)
	for i := 0; i < 10; i++ {
		time.Sleep(2 * time.Second)
		b.add(&nodepb.Prep{Content: "asd", LocalToken: 123, Color: 12})
	}
}

type taaa struct {
}

func (t *taaa) onB(prep *nodepb.BatchedPrep) {

}
