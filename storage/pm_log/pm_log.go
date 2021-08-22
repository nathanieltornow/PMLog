package pm_log

type PMLog struct {
}

func NewPMLog() (*PMLog, error) {
	return &PMLog{}, nil
}

func (log *PMLog) Append(color uint32, record string) (uint64, error) {
	return 0, nil
}

func (log *PMLog) Commit(color uint32, lsn uint64, gsn uint64) error {
	return nil
}

func (log *PMLog) Read(color uint32, gsn uint64) (string, error) {
	return "", nil
}

func (log *PMLog) Trim(color uint32, gsn uint64) error {
	return nil
}
