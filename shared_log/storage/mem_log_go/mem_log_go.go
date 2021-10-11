package mem_log_go

import (
	"fmt"
	"sync"
)

type MemLogGo struct {
	lsnToRecord sync.Map
	gsnToRecord sync.Map
}

func NewMemLogGo() *MemLogGo {
	mlg := new(MemLogGo)
	return mlg
}

func (mlg *MemLogGo) Append(record string, lsn uint64) error {
	mlg.lsnToRecord.Store(lsn, record)
	return nil
}

// Commit commits the records which is stored on local-sequence-number lsn with the global-sequence-number gsn on
// on the log of color. The record can be committed for multiple colors.
func (mlg *MemLogGo) Commit(lsn uint64, gsn uint64) error {
	rec, ok := mlg.lsnToRecord.LoadAndDelete(lsn)
	if !ok {
		return fmt.Errorf("couldn't find record")
	}
	mlg.gsnToRecord.Store(gsn, rec)
	return nil
}

// Read returns the record stored on global-sequence-number gsn on the log of color
func (mlg *MemLogGo) Read(gsn uint64) (string, error) {
	rec, ok := mlg.gsnToRecord.Load(gsn)
	if !ok {
		return "", fmt.Errorf("failed to find record")
	}
	return rec.(string), nil
}

// Trim deletes all records of the log of color before global-sequence-number gsn
func (mlg *MemLogGo) Trim(gsn uint64) error {
	mlg.gsnToRecord.Range(func(key, value interface{}) bool {
		if key.(uint64) < gsn {
			mlg.gsnToRecord.Delete(key)
		}
		return true
	})
	return nil
}
