package mem_log

//#cgo LDFLAGS: -L. -ltbb
//#include "Log.h"
import "C"
import "unsafe"

type Log struct {
	Log C.Log
}

func NewLog() (*Log, error) {
	var ret Log
	ret.Log = C.LogNew()

	return &ret, nil
}

func (log *Log) Append(record string, lsn uint64) error {
	var cRecord = C.CString(record)
	C.cAppend(log.Log, cRecord, C.ulong(lsn))
	C.free(unsafe.Pointer(cRecord))

	return nil
}

func (log *Log) Commit(lsn uint64, gsn uint64) error {
	C.cCommit(log.Log, C.ulong(lsn), C.ulong(gsn))
	return nil
}

func (log *Log) Read(gsn uint64) (string, error) {
	var ret string
	var cRecord = C.CString(ret)
	C.cRead(log.Log, C.ulong(gsn), cRecord)
	ret = C.GoString(cRecord)

	return ret, nil
}

func (log *Log) Trim(gsn uint64) error {
	C.cTrim(log.Log, C.ulong(gsn))
	return nil
}
