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

func (log *PMLog) Append(record string, lsn uint64) error {
	cstr := C.CString(record)
	C.Append(log.pmLog, cstr, C.ulong(lsn))
	C.free(unsafe.Pointer(cstr))

	return nil
}

func (log *PMLog) Commit(lsn uint64, color uint32, gsn uint64) error {
	C.Commit(log.pmLog, C.ulong(lsn), C.uint(color), C.ulong(gsn))
	return nil
}

func (log *PMLog) Read(color uint32, gsn uint64) (string, error) {
	var ret = C.GoString(C.Read(log.pmLog, C.uint(color), C.ulong(gsn)))
	
	return ret, nil
}

func (log *PMLog) Trim(color uint32, gsn uint64) error {
	C.Trim(log.pmLog, C.uint(color), C.ulong(gsn))
	return nil
}
