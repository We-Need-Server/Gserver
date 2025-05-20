package receiver

import (
	"net"
)

type TcpReceiver struct {
	tcpListener *net.TCPListener
}

func NewTcpReceiver(tcpReceiver *net.TCPListener) *TcpReceiver {
	return &TcpReceiver{tcpListener: tcpReceiver}
}

func (r *TcpReceiver) StartTcp() {
	for {
		// 클라이언트 연결 수락
		//conn, err := r.tcpListener.Accept()
		//if err != nil {
		//	fmt.Println("연결 수락 오류:", err)
		//	continue
		//}
		//
		//// 각 연결을 고루틴으로 처리 (동시에 여러 클라이언트 처리 가능)
		//go handleConnection(conn)
	}
}
