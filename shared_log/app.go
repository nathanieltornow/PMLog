package shared_log

import (
	"github.com/sirupsen/logrus"
)

type newRecord struct {
	findToken uint64
	record    string
	color     uint32
	gsn       chan uint64
}

func (sl *SharedLog) Prepare(localToken uint64, color uint32, content string, findToken uint64) error {
	if err := sl.log.Append(content, localToken); err != nil {
		return err
	}
	sl.localToFindTokenMu.Lock()
	sl.localToFindToken[localToken] = findToken
	sl.localToFindTokenMu.Unlock()
	return nil
}

func (sl *SharedLog) Commit(localToken uint64, color uint32, globalToken uint64, isCoordinator bool) error {
	// TODO color
	if err := sl.log.Commit(localToken, globalToken); err != nil {
		return err
	}

	if isCoordinator {
		sl.localToFindTokenMu.Lock()
		findToken, ok := sl.localToFindToken[localToken]
		sl.localToFindTokenMu.Unlock()

		if !ok {
			logrus.Fatalln("Failed to find token")
		}

		sl.pendingAppendsMu.Lock()
		sl.pendingAppends[findToken] <- globalToken
		delete(sl.pendingAppends, localToken)
		sl.pendingAppendsMu.Unlock()
	}

	return nil
}
