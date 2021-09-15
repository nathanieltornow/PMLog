package mem_log

import (
	"fmt"
)

type MemLog struct {
}

func NewMemLog() (*MemLog, error) {
	return &MemLog{}, nil
}

func (log *MemLog) Append(record string, lsn uint64) error {
	fmt.Println("Appending", record, lsn)
	return nil
}

func (log *MemLog) Commit(lsn uint64, gsn uint64) error {
	fmt.Println("Committing", gsn)
	return nil
}

func (log *MemLog) Read(gsn uint64) (string, error) {
	return "", nil
}

func (log *MemLog) Trim(gsn uint64) error {
	return nil
}
