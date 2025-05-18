package receiver

import (
	"WeNeedGameServer/packet"
)

type ChanManager struct {
	CmChan chan packet.PacketI
	nChan  *chan packet.PacketI
}

func NewChanManager(nChan *chan packet.PacketI) *ChanManager {
	return &ChanManager{CmChan: make(chan packet.PacketI), nChan: nChan}
}

func (cm *ChanManager) StartChanManager() {
	for {
		select {
		case p := <-cm.CmChan:
			if p != nil {
				*cm.nChan <- p
			}
		}
	}
}
