package receiver

import (
	"WeNeedGameServer/packet/udp"
)

type UdpChanManager struct {
	CmChan chan udp.PacketI
	nChan  *chan udp.PacketI
}

func NewUdpChanManager(nChan *chan udp.PacketI) *UdpChanManager {
	return &UdpChanManager{CmChan: make(chan udp.PacketI), nChan: nChan}
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
