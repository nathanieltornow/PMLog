package color_service

import (
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	preparationWorkers = 2
)

type prepManager struct {
	mu  sync.Mutex
	app frame.Application

	color uint32

	waitCs map[uint64]chan bool

	prepQueue chan *singlePrepMsg
}

type singlePrepMsg struct {
	content    string
	localToken uint64
	findToken  uint64
}

func newPrepManager(app frame.Application, color uint32) *prepManager {
	pm := new(prepManager)
	pm.app = app
	pm.waitCs = make(map[uint64]chan bool)
	pm.prepQueue = make(chan *singlePrepMsg, 1024)
	pm.color = color
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

func (pm *prepManager) prepareReplica(prepMsg *nodepb.Prep) {
	for i, content := range prepMsg.Contents {
		pm.prepQueue <- &singlePrepMsg{content: content, localToken: prepMsg.LocalToken + uint64(i)}
	}
}

func (pm *prepManager) prepareCoordinator(content string, localToken uint64, findToken uint64) {
	pm.prepQueue <- &singlePrepMsg{content: content, localToken: localToken, findToken: findToken}
}

func (pm *prepManager) executePreparations() {
	for singlePrep := range pm.prepQueue {
		if err := pm.app.Prepare(singlePrep.localToken, pm.color, singlePrep.content, singlePrep.findToken); err != nil {
			logrus.Fatalln(err)
		}

		pm.mu.Lock()
		waitC, ok := pm.waitCs[singlePrep.localToken]
		if !ok {
			waitC = make(chan bool, 1)
			pm.waitCs[singlePrep.localToken] = waitC
		}
		pm.mu.Unlock()

		waitC <- true

	}
}
