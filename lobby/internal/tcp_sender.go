package internal

import (
	"WeNeedGameServer/external/db"
	"WeNeedGameServer/protocol/tcp"
	"fmt"
)

type TcpSender struct {
	listenUdpAddr string
	blueTeamDb    map[uint32]*db.User
	redTeamDb     map[uint32]*db.User
}

func NewTcpSender(listenUdpAddr string, blueTeamDb map[uint32]*db.User, redTeamDb map[uint32]*db.User) *TcpSender {
	return &TcpSender{
		listenUdpAddr: listenUdpAddr,
		blueTeamDb:    blueTeamDb,
		redTeamDb:     redTeamDb,
	}
}

func (s *TcpSender) ProcessMessage(message *tcp.Message) {
	switch message.SenderType {
	case tcp.SendByBroadCast:
		s.sendByBroadCast(message.Data)
		break
	case tcp.SendByUniCast:
		s.sendByUniCast(message.UserId, message.Data)
		break
	case tcp.SendByMultiCast:
		s.sendByMultiCast(message.UserId, message.Data)
		break
	}
}

//func (s *TcpSender) makePacket(pKind uint8, userId uint32) tcp.PacketI {
//	switch pKind {
//	case 'I':
//		if u := s.getUser(userId); u != nil {
//			return tserver.NewConnectionResponsePacket(u.QPort, s.listenUdpAddr)
//		}
//	case 'U':
//		var userList []uint32
//		for key, _ := range s.blueTeamDb {
//			userList = append(userList, key)
//		}
//		for key, _ := range s.redTeamDb {
//			userList = append(userList, key)
//		}
//		return tserver.NewUserConnectionUpdatePacket(userList)
//	case 'S':
//		break
//	}
//	return nil
//}

func (s *TcpSender) sendByBroadCast(packet tcp.PacketI) {
	if packet != nil {
		for _, u := range s.blueTeamDb {
			if u.QPort != 0 && u.TcpConn != nil {
				if _, err := u.TcpConn.Write(packet.Serialize()); err != nil {
					fmt.Println("send broadcast error", err)
				}
			}
		}
		for _, u := range s.redTeamDb {
			if u.QPort != 0 && u.TcpConn != nil {
				if _, err := u.TcpConn.Write(packet.Serialize()); err != nil {
					fmt.Println("send broadcast error", err)
				}
			}
		}
	}
}

func (s *TcpSender) sendByUniCast(userId uint32, packet tcp.PacketI) {
	if packet != nil {
		if u := s.getUser(userId); u.QPort != 0 && u.TcpConn != nil {
			if _, err := u.TcpConn.Write(packet.Serialize()); err != nil {
				fmt.Println("send uni cast error", err)
			}
		}
	}
}

func (s *TcpSender) sendByMultiCast(excludedUserId uint32, packet tcp.PacketI) {
	if packet != nil {
		for userId, u := range s.blueTeamDb {
			if userId == excludedUserId {
				continue
			}
			if u.QPort != 0 && u.TcpConn != nil {
				if _, err := u.TcpConn.Write(packet.Serialize()); err != nil {
					fmt.Println("send multicast error", err)
				}
			}
		}
		for userId, u := range s.redTeamDb {
			if userId == excludedUserId {
				continue
			}
			if u.QPort != 0 && u.TcpConn != nil {
				if _, err := u.TcpConn.Write(packet.Serialize()); err != nil {
					fmt.Println("send multicast error", err)
				}
			}
		}
	}
}

func (s *TcpSender) getUser(userId uint32) *db.User {
	if blueUser, bExists := s.blueTeamDb[userId]; bExists {
		return blueUser
	} else if redUser, rExists := s.redTeamDb[userId]; rExists {
		return redUser
	} else {
		return nil
	}
}
