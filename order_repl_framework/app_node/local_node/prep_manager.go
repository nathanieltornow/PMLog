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

	prepQueue chan *prepMsgFindToken
}

type prepMsgFindToken struct {
	prepMsg   *nodepb.Prep
	findToken uint64
}

func newPrepManager(app frame.Application) *prepManager {
	pm := new(prepManager)
	pm.app = app
	pm.waitCs = make(map[uint64]chan bool)
	pm.prepQueue = make(chan *prepMsgFindToken, 1024)

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

func (pm *prepManager) prepare(prepMsg *nodepb.Prep, findToken uint64) {
	pm.prepQueue <- &prepMsgFindToken{prepMsg: prepMsg, findToken: findToken}
}

func (pm *prepManager) executePreparations() {
	for prepMsgFT := range pm.prepQueue {
		for i, content := range prepMsgFT.prepMsg.Contents {
			if err := pm.app.Prepare(prepMsgFT.prepMsg.LocalToken+uint64(i), prepMsgFT.prepMsg.Color, content, prepMsgFT.findToken); err != nil {
				logrus.Fatalln(err)
			}
			pm.mu.Lock()
			waitC, ok := pm.waitCs[prepMsgFT.prepMsg.LocalToken]
			if !ok {
				waitC = make(chan bool, 1)
				pm.waitCs[prepMsgFT.prepMsg.LocalToken] = waitC
			}
			pm.mu.Unlock()

			waitC <- true
		}

	}
}
