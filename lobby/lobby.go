package lobby

import (
	"WeNeedGameServer/external/db"
	"WeNeedGameServer/lobby/internal"
	"log"
	"net"
)

type Lobby struct {
	userDb        *db.Db
	tcpListener   *net.TCPListener
	tcpReceiver   *internal.TcpReceiver
	tcpSender     *internal.TcpSender
	listenUdpAddr string
	listenTcpAddr string
}

func NewLobby(listenUdpAddr string, listenTcpAddr string) *Lobby {
	return &Lobby{
		userDb:        db.NewUserDb(),
		tcpReceiver:   nil,
		tcpSender:     nil,
		listenUdpAddr: listenUdpAddr,
		listenTcpAddr: listenTcpAddr,
	}
}

func (l *Lobby) communicateTcpSender(message internal.TcpReceiverMessage) {
	l.tcpSender.ProcessMessage(message)
}

func (l *Lobby) ReadyTcp() (*internal.TcpReceiver, *internal.TcpSender) {
	tcpServerPoint, tcpResolveErr := net.ResolveTCPAddr("tcp", l.listenTcpAddr)
	if tcpResolveErr != nil {
		log.Panicln("네트워크 리졸버 오류")
	}
	tcpLn, tcpListenErr := net.ListenTCP("tcp", tcpServerPoint)
	if tcpListenErr != nil {
		log.Panicln("리슨 오류")
	}
	l.tcpListener = tcpLn
	l.tcpReceiver = internal.NewTcpReceiver(tcpLn, l.userDb.Login, l.communicateTcpSender)
	l.tcpSender = internal.NewTcpSender(l.listenUdpAddr, l.userDb.BlueTeamDb, l.userDb.RedTeamDb)
	return l.tcpReceiver, l.tcpSender
}
