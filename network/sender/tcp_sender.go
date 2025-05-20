package sender

import (
	"WeNeedGameServer/db/user"
	"WeNeedGameServer/network/receiver"
)

type TcpSender struct {
	listenUdpAddr string
	blueTeamDb    map[uint32]*user.User
	redTeamDb     map[uint32]*user.User
}

func NewTcpSender(listenUdpAddr string, blueTeamDb map[uint32]*user.User, redTeamDb map[uint32]*user.User) *TcpSender {
	return &TcpSender{
		listenUdpAddr: listenUdpAddr,
		blueTeamDb:    blueTeamDb,
		redTeamDb:     redTeamDb,
	}
}

func (s *TcpSender) ProcessMessage(message receiver.TcpReceiverMessage) {

}
