package app_node

import (
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	commitWorkers = 2
)

type comManager struct {
	mu sync.Mutex

	app frame.Application

	numOfAcks map[uint64]uint32
	peers     uint32

	comCh chan *comMsgIsCoor

	waitCs   map[uint64]chan bool
	waitCsMu sync.Mutex
}

type comMsgIsCoor struct {
	comMsg *nodepb.Com
	isCoor bool
}

func newComManager(app frame.Application) *comManager {
	comMan := new(comManager)
	comMan.app = app
	comMan.numOfAcks = make(map[uint64]uint32)
	comMan.comCh = make(chan *comMsgIsCoor, 1024)
	comMan.waitCs = make(map[uint64]chan bool)

	for i := 0; i < commitWorkers; i++ {
		go comMan.executeCommits()
	}

	return comMan
}

func (cm *comManager) addPeer() {
	cm.mu.Lock()
	cm.peers++
	cm.mu.Unlock()
}

func (cm *comManager) waitForCom(localToken uint64) {
	cm.waitCsMu.Lock()
	waitC, ok := cm.waitCs[localToken]
	if !ok {
		waitC = make(chan bool, 1)
		cm.waitCs[localToken] = waitC
	}
	cm.waitCsMu.Unlock()
	<-waitC
}

func (cm *comManager) receiveAck(ackMsg *nodepb.Ack) {
	cm.mu.Lock()
	cm.numOfAcks[ackMsg.LocalToken]++
	num := cm.numOfAcks[ackMsg.LocalToken]
	peers := cm.peers
	cm.mu.Unlock()
	if num == peers {
		cm.comCh <- &comMsgIsCoor{
			comMsg: &nodepb.Com{
				LocalToken:  ackMsg.LocalToken,
				Color:       ackMsg.Color,
				GlobalToken: ackMsg.GlobalToken,
			},
			isCoor: true,
		}
	}
}

func (cm *comManager) receiveCom(comMsg *nodepb.Com) {
	cm.comCh <- &comMsgIsCoor{
		comMsg: &nodepb.Com{
			LocalToken:  comMsg.LocalToken,
			Color:       comMsg.Color,
			GlobalToken: comMsg.GlobalToken,
		},
		isCoor: false,
	}
}

func (cm *comManager) executeCommits() {
	for comIsCoor := range cm.comCh {
		if err := cm.app.Commit(comIsCoor.comMsg.LocalToken, comIsCoor.comMsg.Color, comIsCoor.comMsg.GlobalToken, comIsCoor.isCoor); err != nil {
			logrus.Fatalln(err)
		}

		cm.waitCsMu.Lock()
		waitC, ok := cm.waitCs[comIsCoor.comMsg.LocalToken]
		if !ok {
			waitC = make(chan bool, 1)
			cm.waitCs[comIsCoor.comMsg.LocalToken] = waitC
		}
		cm.waitCsMu.Unlock()

		waitC <- true
	}
}
