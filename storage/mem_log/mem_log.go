package mem_log

import (
	"fmt"
	"time"
)

type MemLog struct {
}

func NewMemLog() (*MemLog, error) {
	return &MemLog{}, nil
}

func (log *MemLog) Append(record []byte, lsn uint64) error {
	fmt.Println("Appending", record, lsn, time.Now().Nanosecond())
	return nil
}

func (log *MemLog) Commit(lsn uint64, gsn uint64) error {
	fmt.Println("Committing", lsn, time.Now().Nanosecond())
	return nil
}

func (log *MemLog) Read(gsn uint64) ([]byte, error) {
	return []byte{}, nil
}

func (log *MemLog) Trim(gsn uint64) error {
	return nil
}
