package mem_log

//#cgo LDFLAGS: -L. -lstorage -ltbb
//#include "Log.h"
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

type MemLog struct {
	memLog C.Log
}

func NewMemLog() (*MemLog, error) {
	var ret MemLog
	ret.memLog = C.LogNew()

	return &ret, nil
}

func (log *MemLog) Append(record string, lsn uint64) error {
	var cRecord = C.CString(record)
	C.cAppend(log.memLog, cRecord, C.ulong(lsn))
	C.free(unsafe.Pointer(cRecord))
	fmt.Println("Appending", record, lsn, time.Now().Nanosecond())

	return nil
}

func (log *MemLog) Commit(lsn uint64, gsn uint64) error {
	C.cCommit(log.memLog, C.ulong(lsn), C.ulong(gsn))
	fmt.Println("Committing", lsn, time.Now().Nanosecond())
	return nil
}

func (log *MemLog) Read(gsn uint64) (string, error) {
	var ret string
	var cRecord = C.CString(ret)
	C.cRead(log.memLog, C.ulong(gsn), cRecord)
	ret = C.GoString(cRecord)

	return ret, nil
}

func (log *MemLog) Trim(gsn uint64) error {
	C.cTrim(log.memLog, C.ulong(gsn))

	return nil
}
