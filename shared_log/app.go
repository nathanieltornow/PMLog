package shared_log

import (
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
)

type newRecord struct {
	findToken uint64
	record    string
	color     uint32
	gsn       chan uint64
}

func (sl *SharedLog) MakeCommitRequests(ch chan *frame.CommitRequest) error {
	for newRec := range sl.newAppends {
		comReq := &frame.CommitRequest{
			Color:     newRec.color,
			Content:   newRec.record,
			FindToken: newRec.findToken,
		}
		ch <- comReq
	}
	return nil
}

func (sl *SharedLog) Prepare(localToken uint64, color uint32, content string, findToken uint64) error {
	// TODO Color
	if err := sl.log.Append(content, localToken); err != nil {
		return err
	}
	sl.localToFindTokenMu.Lock()
	sl.localToFindToken[localToken] = findToken
	sl.localToFindTokenMu.Unlock()
	return nil
}

func (sl *SharedLog) IsPrepared(localToken uint64) bool {
	return true
}

func (sl *SharedLog) Commit(localToken uint64, color uint32, globalToken uint64) error {
	// TODO color
	if err := sl.log.Commit(localToken, globalToken); err != nil {
		return err
	}

	sl.pendingAppendsMu.Lock()
	sl.pendingAppends[localToken] <- globalToken
	delete(sl.pendingAppends, localToken)
	sl.pendingAppendsMu.Unlock()
	return nil
}
