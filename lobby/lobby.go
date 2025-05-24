package lobby

import (
	"WeNeedGameServer/external/db"
	"WeNeedGameServer/game_manager"
	"WeNeedGameServer/lobby/internal"
	"WeNeedGameServer/protocol/tcp"
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
	matchScore    uint16
	gameManager   *game_manager.GameManager
}

func NewLobby(userDb *db.Db, listenUdpAddr string, listenTcpAddr string, matchScore uint16) *Lobby {
	return &Lobby{
		userDb:        userDb,
		tcpReceiver:   nil,
		tcpSender:     nil,
		listenUdpAddr: listenUdpAddr,
		listenTcpAddr: listenTcpAddr,
		matchScore:    matchScore,
		gameManager:   nil,
	}
}

func (l *Lobby) communicateTcpSender(message *tcp.Message) {
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
	l.tcpSender = internal.NewTcpSender(l.listenUdpAddr, l.userDb.BlueTeamDb, l.userDb.RedTeamDb)
	l.gameManager = game_manager.NewGameManager(10, l.userDb, l.matchScore, l.communicateTcpSender, l.listenUdpAddr)
	l.tcpReceiver = internal.NewTcpReceiver(tcpLn, l.userDb.Login, l.userDb.BlueTeamDb, l.userDb.BlueTeamDb, l.communicateTcpSender, l.matchScore, l.listenUdpAddr, l.gameManager.StartGameManager, &l.gameManager.GameStatus, l.gameManager.SendGameInitPacket)
	return l.tcpReceiver, l.tcpSender
}
