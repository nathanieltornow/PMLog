package mem_log

type MemLog struct {
}

func (log *MemLog) Append(color uint32, record string) (uint64, error) {
	return 0, nil
}

func (log *MemLog) Commit(color uint32, lsn uint64, gsn uint64) error {
	return nil
}

func (log *MemLog) Read(color uint32, gsn uint64) (string, error) {
	return "", nil
}

func (log *MemLog) Trim(color uint32, gsn uint64) error {
	return nil
}
