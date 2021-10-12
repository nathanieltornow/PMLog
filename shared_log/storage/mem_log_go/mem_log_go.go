package mem_log_go

import (
	"fmt"
	"sync"
)

type MemLogGo struct {
	lsnMu       sync.Mutex
	lsnToRecord map[uint64]string
	gsnMu       sync.RWMutex
	gsnToRecord map[uint64]string
}

func NewMemLogGo() *MemLogGo {
	mlg := new(MemLogGo)
	mlg.lsnToRecord = make(map[uint64]string)
	mlg.gsnToRecord = make(map[uint64]string)
	return mlg
}

func (mlg *MemLogGo) Append(record string, lsn uint64) error {
	mlg.lsnMu.Lock()
	mlg.lsnToRecord[lsn] = record
	mlg.lsnMu.Unlock()
	return nil
}

// Commit commits the records which is stored on local-sequence-number lsn with the global-sequence-number gsn on
// on the log of color. The record can be committed for multiple colors.
func (mlg *MemLogGo) Commit(lsn uint64, gsn uint64) error {
	mlg.lsnMu.Lock()
	rec, ok := mlg.lsnToRecord[lsn]
	delete(mlg.lsnToRecord, lsn)
	mlg.lsnMu.Unlock()
	if !ok {
		return fmt.Errorf("couldn't find record")
	}
	mlg.gsnMu.Lock()
	mlg.gsnToRecord[gsn] = rec
	mlg.gsnMu.Unlock()
	return nil
}

// Read returns the record stored on global-sequence-number gsn on the log of color
func (mlg *MemLogGo) Read(gsn uint64) (string, error) {
	mlg.gsnMu.RLock()
	defer mlg.gsnMu.RUnlock()
	rec, ok := mlg.gsnToRecord[gsn]
	if !ok {
		return "", fmt.Errorf("failed to find record")
	}
	return rec, nil
}

// Trim deletes all records of the log of color before global-sequence-number gsn
func (mlg *MemLogGo) Trim(gsn uint64) error {
	mlg.gsnMu.Lock()
	for recGsn := range mlg.gsnToRecord {
		if recGsn < gsn {
			delete(mlg.gsnToRecord, recGsn)
		}
	}
	mlg.gsnMu.Unlock()
	return nil
}
