package mem_log

type MemLog struct {
}

func NewMemLog() (*MemLog, error) {
	return &MemLog{}, nil
}

func (log *MemLog) Append(record string) (uint64, error) {
	return 0, nil
}

func (log *MemLog) Commit(lsn uint64, color uint32, gsn uint64) error {
	return nil
}

func (log *MemLog) Read(color uint32, gsn uint64) (string, error) {
	return "", nil
}

func (log *MemLog) Trim(color uint32, gsn uint64) error {
	return nil
}
