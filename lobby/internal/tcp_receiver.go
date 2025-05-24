package internal

import (
	"WeNeedGameServer/external/db"
	"WeNeedGameServer/game_manager"
	"WeNeedGameServer/protocol/tcp"
	"WeNeedGameServer/protocol/tcp/tclient"
	"WeNeedGameServer/protocol/tcp/tserver"
	"WeNeedGameServer/util"
	"fmt"
	"net"
	"sync"
	"time"
)

type TcpReceiver struct {
	tcpListener            *net.TCPListener
	loginFunc              func(uint32, net.Conn) (uint32, db.Team, error)
	blueTeamDb             map[uint32]*db.User
	redTeamDb              map[uint32]*db.User
	communicateSenderFunc  func(*tcp.Message)
	matchScore             uint16
	listenUdpAddr          string
	startGameFunc          func()
	gameStatus             *game_manager.GameStatus
	sendGameInitPacketFunc func(uint32)

	// 타이머 관련 필드 추가
	gameTimer     *time.Timer
	timerMutex    sync.Mutex
	isTimerActive bool
}

func NewTcpReceiver(tcpListener *net.TCPListener, loginFunc func(uint32, net.Conn) (uint32, db.Team, error), blueTeamDb map[uint32]*db.User, redTeamDb map[uint32]*db.User, communicateSenderFunc func(*tcp.Message), matchScore uint16, listenUdpAddr string, startGameFunc func(), gameStatus *game_manager.GameStatus, sendGameInitPacketFunc func(uint32)) *TcpReceiver {
	return &TcpReceiver{
		tcpListener:            tcpListener,
		loginFunc:              loginFunc,
		blueTeamDb:             blueTeamDb,
		redTeamDb:              redTeamDb,
		communicateSenderFunc:  communicateSenderFunc,
		matchScore:             matchScore,
		listenUdpAddr:          listenUdpAddr,
		startGameFunc:          startGameFunc,
		gameStatus:             gameStatus,
		sendGameInitPacketFunc: sendGameInitPacketFunc,
		gameTimer:              nil,
		isTimerActive:          false,
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

// 30초 타이머 시작 함수
func (r *TcpReceiver) startGameTimer() {
	r.timerMutex.Lock()
	defer r.timerMutex.Unlock()

	// 기존 타이머가 있다면 정지
	if r.gameTimer != nil {
		r.gameTimer.Stop()
	}

	// 새로운 30초 타이머 시작
	r.gameTimer = time.NewTimer(30 * time.Second)
	r.isTimerActive = true

	go func() {
		<-r.gameTimer.C
		r.timerMutex.Lock()
		r.isTimerActive = false
		r.timerMutex.Unlock()

		// 30초 후 게임 시작
		r.startGameFunc()
	}()
}

// 타이머 취소 함수
func (r *TcpReceiver) cancelGameTimer() {
	r.timerMutex.Lock()
	defer r.timerMutex.Unlock()

	if r.gameTimer != nil && r.isTimerActive {
		r.gameTimer.Stop()
		r.isTimerActive = false
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
			fmt.Println("로그인 성공", connectionRequestPacket.UserId)
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
			switch *r.gameStatus {
			case game_manager.GameReady:
				if len(r.redTeamDb) > 0 && len(r.blueTeamDb) > 0 {
					r.startGameTimer()
				} else {
					r.cancelGameTimer()
				}
				break
			case game_manager.RoundStart:
				r.sendGameInitPacketFunc(connectionRequestPacket.UserId)
				break
			default:
				fmt.Println("잠시 기다려주십시오")
			}
		}
		break
	case 'T':
		// 유저 끊는건 일단 나중에
		break
	}
}
