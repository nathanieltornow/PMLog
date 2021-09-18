package app_node

import (
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	"google.golang.org/protobuf/proto"
	"time"
)

type onBatch func(prep *nodepb.BatchedPrep)

type prepBatch struct {
	onB     onBatch
	maxSize int

	prepCh chan *nodepb.Prep
}

func newPrepBatch(onB onBatch, maxSize int, interval time.Duration) *prepBatch {
	b := new(prepBatch)
	b.maxSize = maxSize
	b.onB = onB
	b.prepCh = make(chan *nodepb.Prep, 1024)
	go b.batch(interval)
	return b
}

func (b *prepBatch) add(prepMsg *nodepb.Prep) {
	b.prepCh <- prepMsg
}

func (b *prepBatch) batch(interval time.Duration) {

	var current *nodepb.BatchedPrep
	newBatch := true

	var send <-chan time.Time

	for {
		select {
		case prepMsg := <-b.prepCh:
			if newBatch {
				send = time.After(interval)
				current = &nodepb.BatchedPrep{Preps: make([]*nodepb.Prep, 0)}
				newBatch = false
			}
			current.Preps = append(current.Preps, prepMsg)

			if proto.Size(current) > b.maxSize {
				b.onB(current)
				newBatch = true
			}

		case <-send:
			b.onB(current)
			newBatch = true
		}
	}
}
