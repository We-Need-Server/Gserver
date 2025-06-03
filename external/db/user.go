package db

import (
	"WeNeedGameServer/game_type"
	"net"
)

type User struct {
	QPort   uint32
	TcpConn net.Conn
	Team    game_type.Team
}

func NewUser(team game_type.Team) *User {
	return &User{
		QPort:   0,
		TcpConn: nil,
		Team:    team,
	}
}
