package sender

import (
	"WeNeedGameServer/network/receiver"
	"net"
)

type TcpSender struct {
	listenUdpAddr string
	tcpConnTable  map[uint32]*net.Conn
}

func NewTcpSender(listenUdpAddr string, tcpConnTable map[uint32]*net.Conn) *TcpSender {
	return &TcpSender{
		listenUdpAddr: listenUdpAddr,
		tcpConnTable:  tcpConnTable,
	}
}

func (s *TcpSender) ProcessMessage(message receiver.TcpReceiverMessage) {

}
