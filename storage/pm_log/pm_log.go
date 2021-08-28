package pm_log

type PMLog struct {
}

func NewPMLog() (*PMLog, error) {
	return &PMLog{}, nil
}

func (log *PMLog) Append(record string, lsn uint64) error {
	return nil
}

func (log *PMLog) Commit(lsn uint64, color uint32, gsn uint64) error {
	return nil
}

func (log *PMLog) Read(color uint32, gsn uint64) (string, error) {
	return "", nil
}

func (log *PMLog) Trim(color uint32, gsn uint64) error {
	return nil
}
