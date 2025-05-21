package db

import (
	"net"
)

type User struct {
	QPort   uint32
	TcpConn net.Conn
	Team    Team
}

func NewUser(team Team) *User {
	return &User{
		QPort:   0,
		TcpConn: nil,
		Team:    team,
	}
}
