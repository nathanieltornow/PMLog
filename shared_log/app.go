package shared_log

func (sl *SharedLog) Prepare(localToken uint64, color uint32, content string) error {
	if err := sl.log.Append(content, localToken); err != nil {
		return err
	}
	return nil
}

func (sl *SharedLog) Commit(localToken uint64, color uint32, globalToken uint64) error {
	// TODO color
	if err := sl.log.Commit(localToken, globalToken); err != nil {
		return err
	}
	return nil
}

func (sl *SharedLog) Acknowledge(localToken uint64, color uint32, globalToken uint64) error {
	defer func() {
		sl.pendingAppendsMu.Lock()
		delete(sl.pendingAppends, localToken)
		sl.pendingAppendsMu.Unlock()
	}()
	sl.pendingAppendsMu.Lock()
	defer sl.pendingAppendsMu.Unlock()
	p, ok := sl.pendingAppends[localToken]
	if !ok {
		p = make(chan uint64, 1)
		sl.pendingAppends[localToken] = p
	}
	p <- globalToken
	return nil
}
