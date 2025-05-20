package receiver

import (
	"net"
)

type TcpReceiverMessage struct {
}

type TcpReceiver struct {
	tcpListener           *net.TCPListener
	tcpConnTable          map[uint32]*net.Conn
	loginFunc             func(uint32, *net.TCPAddr) (uint32, error)
	communicateSenderFunc func(TcpReceiverMessage)
}

func NewTcpReceiver(tcpReceiver *net.TCPListener, tcpConnTable map[uint32]*net.Conn, loginFunc func(uint32, *net.TCPAddr) (uint32, error), communicateSenderFunc func(TcpReceiverMessage)) *TcpReceiver {
	return &TcpReceiver{
		tcpListener:           tcpReceiver,
		tcpConnTable:          tcpConnTable,
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
