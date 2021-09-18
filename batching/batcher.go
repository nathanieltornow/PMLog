package batching

import (
	"time"
)

type OnBatch func(batchS BatchSubject)

type Batcher struct {
	batchS BatchSubject

	onB OnBatch

	maxSize int

	mergeCh chan BatchSubject
}

type BatchSubject interface {
	Merge(other BatchSubject)
	Size() int
}

func NewBatcher(batchS BatchSubject, maxSize int, onB OnBatch, interval time.Duration) *Batcher {
	b := new(Batcher)
	b.batchS = batchS
	b.onB = onB
	b.maxSize = maxSize
	b.mergeCh = make(chan BatchSubject, 1024)
	go b.batch(interval)
	return b
}

func (b *Batcher) Add(subject BatchSubject) {
	b.mergeCh <- subject
}

func (b *Batcher) batch(interval time.Duration) {
	ticker := time.Tick(interval)

	var current BatchSubject

	newBatch := true

	for {
		select {

		case batchSub := <-b.mergeCh:
			if newBatch {
				current = batchSub
				newBatch = false
				continue
			}

			current.Merge(batchSub)

			if current.Size() > b.maxSize {
				b.onB(current)
				newBatch = true
			}

		case <-ticker:
			if newBatch {
				continue
			}
			b.onB(current)
			newBatch = true
		}
	}
}
