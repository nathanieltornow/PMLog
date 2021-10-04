package local_node

import (
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	"github.com/sirupsen/logrus"
	"sync"
)

type prepManager struct {
	mu  sync.Mutex
	app frame.Application

	waitCs map[uint64]chan bool

	prepQueue chan *nodepb.Prep
}

type prepMsgFindToken struct {
	prepMsg   *nodepb.Prep
	findToken uint64
}

func newPrepManager(app frame.Application) *prepManager {
	pm := new(prepManager)
	pm.app = app
	pm.waitCs = make(map[uint64]chan bool)
	pm.prepQueue = make(chan *nodepb.Prep, 1024)

	go pm.executePreparations()

	return pm
}

func (pm *prepManager) waitForPrep(localToken uint64) {
	pm.mu.Lock()
	waitC, ok := pm.waitCs[localToken]
	if !ok {
		waitC = make(chan bool, 1)
		pm.waitCs[localToken] = waitC
	}
	pm.mu.Unlock()

	defer func() {
		pm.mu.Lock()
		delete(pm.waitCs, localToken)
		pm.mu.Unlock()
	}()
	<-waitC
}

func (pm *prepManager) prepare(prepMsg *nodepb.Prep) {
	pm.prepQueue <- prepMsg
}

func (pm *prepManager) executePreparations() {
	for prepMsg := range pm.prepQueue {
		for i, content := range prepMsg.Contents {
			if err := pm.app.Prepare(prepMsg.LocalToken+uint64(i), prepMsg.Color, content); err != nil {
				logrus.Fatalln(err)
			}
			pm.mu.Lock()
			waitC, ok := pm.waitCs[prepMsg.LocalToken]
			if !ok {
				waitC = make(chan bool, 1)
				pm.waitCs[prepMsg.LocalToken] = waitC
			}
			pm.mu.Unlock()

			waitC <- true
		}

	}
}
