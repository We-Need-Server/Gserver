package internal

import (
	"WeNeedGameServer/external/db"
	"WeNeedGameServer/protocol/tcp"
	"WeNeedGameServer/protocol/tcp/tclient"
	"WeNeedGameServer/protocol/tcp/tserver"
	"WeNeedGameServer/util"
	"fmt"
	"net"
)

type TcpReceiver struct {
	tcpListener           *net.TCPListener
	loginFunc             func(uint32, net.Conn) (uint32, db.Team, error)
	blueTeamDb            map[uint32]*db.User
	redTeamDb             map[uint32]*db.User
	communicateSenderFunc func(*tcp.Message)
	matchScore            uint16
	listenUdpAddr         string
	startGameFunc         func()
}

func NewTcpReceiver(tcpListener *net.TCPListener, loginFunc func(uint32, net.Conn) (uint32, db.Team, error), blueTeamDb map[uint32]*db.User, redTeamDb map[uint32]*db.User, communicateSenderFunc func(*tcp.Message), matchScore uint16, listenUdpAddr string, startGameFunc func()) *TcpReceiver {
	return &TcpReceiver{
		tcpListener:           tcpListener,
		loginFunc:             loginFunc,
		blueTeamDb:            blueTeamDb,
		redTeamDb:             redTeamDb,
		communicateSenderFunc: communicateSenderFunc,
		matchScore:            matchScore,
		listenUdpAddr:         listenUdpAddr,
		startGameFunc:         startGameFunc,
	}
}

func (r *TcpReceiver) StartTcp() {
	for {
		conn, err := r.tcpListener.Accept()
		if err != nil {
			fmt.Println("연결 수락 오류:", err)
			continue
		}
		go r.handleConnection(conn)
	}
}

func (r *TcpReceiver) handleConnection(conn net.Conn) {
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}

		r.processData(conn, buffer[:n])
	}
}

func (r *TcpReceiver) processData(conn net.Conn, b []byte) {
	pKind := b[0]
	switch pKind {
	case 'H':
		connectionRequestPacket := tclient.NewConnectionRequestPacket(util.ConvertBinaryToUint32(b[1:]))
		qPort, team, err := r.loginFunc(connectionRequestPacket.UserId, conn)
		if err != nil {
			fmt.Println("login fail", connectionRequestPacket.UserId)
		} else {
			r.communicateSenderFunc(tcp.NewUniCastMessage(connectionRequestPacket.UserId, tserver.NewConnectionResponsePacket(qPort, r.listenUdpAddr, r.matchScore)))
			r.communicateSenderFunc(tcp.NewMultiCastMessage(connectionRequestPacket.UserId, tserver.NewUserConnectionPUpdatePacket([]tserver.UserTeamStatus{tserver.NewUserTeamStatus(connectionRequestPacket.UserId, team)})))
			var userList []tserver.UserTeamStatus
			for userId, u := range r.blueTeamDb {
				userList = append(userList, tserver.NewUserTeamStatus(userId, u.Team))
			}
			for userId, u := range r.redTeamDb {
				userList = append(userList, tserver.NewUserTeamStatus(userId, u.Team))
			}
			r.communicateSenderFunc(tcp.NewUniCastMessage(connectionRequestPacket.UserId, tserver.NewUserConnectionPUpdatePacket(userList)))
		}
		break
	case 'T':
		// 유저 끊는건 일단 나중에
		break
	}
}
