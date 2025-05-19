package receiver

import "net"

type TcpReceiver struct {
	tcpListener *net.TCPListener
}

func NewTcpReceiver(tcpReceiver *net.TCPListener) *TcpReceiver {
	return &TcpReceiver{tcpListener: tcpReceiver}
}
