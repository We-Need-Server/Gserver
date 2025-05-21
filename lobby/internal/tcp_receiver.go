package internal

import (
	"WeNeedGameServer/protocol/tcp"
	"net"
)

type TcpReceiver struct {
	tcpListener           *net.TCPListener
	loginFunc             func(uint32, net.Conn) error
	communicateSenderFunc func(*tcp.ReceiverMessage)
}

func NewTcpReceiver(tcpListener *net.TCPListener, loginFunc func(uint32, net.Conn) error, communicateSenderFunc func(*tcp.ReceiverMessage)) *TcpReceiver {
	return &TcpReceiver{
		tcpListener:           tcpListener,
		loginFunc:             loginFunc,
		communicateSenderFunc: communicateSenderFunc,
	}
}

func (r *TcpReceiver) StartTcp() {
	//for {
	//	// 클라이언트 연결 수락
	//	conn, err := r.tcpListener.Accept()
	//	if err != nil {
	//		fmt.Println("연결 수락 오류:", err)
	//		continue
	//	}
	//
	//}
}
