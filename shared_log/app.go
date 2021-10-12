package shared_log

import (
	pb "github.com/nathanieltornow/PMLog/shared_log/shared_logpb"
)

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
	sl.lsnTokenMu.RLock()
	tu, ok := sl.lsnToken[localToken]
	sl.lsnTokenMu.RUnlock()
	if ok {
		defer func() {
			sl.lsnTokenMu.Lock()
			delete(sl.lsnToken, localToken)
			sl.lsnTokenMu.Unlock()
		}()
		sl.clientChsMu.RLock()
		sl.clientChs[tu.clientID] <- &pb.AppendResponse{Gsn: globalToken, Token: tu.token}
		sl.clientChsMu.RUnlock()
		return nil
	}

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
