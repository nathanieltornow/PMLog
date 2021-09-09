package pm_log

//#cgo LDFLAGS: -L. -lstorage -lpmemobj
//#include "PMLog.h"
import "C"
import "unsafe"

type PMLog struct {
	pmLog C.PMLog
}

func NewPMLog() (*PMLog, error) {
	var ret PMLog
	ret.pmLog = C.startUp()
	
	return &ret, nil
}

func (log *PMLog) Append(record []byte, lsn uint64) error {
	C.cAppend(log.pmLog, (*C.char)(unsafe.Pointer(&record)), C.ulong(lsn))

	return nil
}

func (log *PMLog) Commit(lsn uint64,gsn uint64) error {
	C.cCommit(log.pmLog, C.ulong(lsn), C.ulong(gsn))
	return nil
}

func (log *PMLog) Read(gsn uint64) ([]byte, error) {
	var ret []byte
	C.cRead(log.pmLog, C.ulong(gsn), (*C.char)(unsafe.Pointer(&ret)))
	
	return ret, nil
}

func (log *PMLog) Trim(gsn uint64) error {
	C.cTrim(log.pmLog, C.ulong(gsn))
	return nil
}






