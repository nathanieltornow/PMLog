package dummy_log

type DummyLog struct {
}

func NewDummyLog() (*DummyLog, error) {
	return &DummyLog{}, nil
}

func (log *DummyLog) Append(record string, lsn uint64) error {
	return nil
}

func (log *DummyLog) Commit(lsn uint64, gsn uint64) error {
	return nil
}

func (log *DummyLog) Read(gsn uint64) (string, error) {
	return "", nil
}

func (log *DummyLog) Trim(gsn uint64) error {
	return nil
}
