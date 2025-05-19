package receiver

import (
	"WeNeedGameServer/packet"
)

type UdpChanManager struct {
	CmChan chan packet.PacketI
	nChan  *chan packet.PacketI
}

func NewUdpChanManager(nChan *chan packet.PacketI) *UdpChanManager {
	return &UdpChanManager{CmChan: make(chan packet.PacketI), nChan: nChan}
}

func (cm *UdpChanManager) StartChanManager() {
	for {
		select {
		case p := <-cm.CmChan:
			if p != nil {
				*cm.nChan <- p
			}
		}
	}
}
