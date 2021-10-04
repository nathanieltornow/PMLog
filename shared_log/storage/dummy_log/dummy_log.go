package dummy_log

import (
	"time"
)

type DummyLog struct {
}

func NewDummyLog() (*DummyLog, error) {
	return &DummyLog{}, nil
}

func (log *DummyLog) Append(record string, lsn uint64) error {
	time.Sleep(time.Microsecond * 30)
	return nil
}

func (log *DummyLog) Commit(lsn uint64, gsn uint64) error {
	time.Sleep(time.Microsecond * 20)
	return nil
}

func (log *DummyLog) Read(gsn uint64) (string, error) {
	return "", nil
}

func (log *DummyLog) Trim(gsn uint64) error {
	return nil
}
