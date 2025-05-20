package user

import "net"

type User struct {
	QPort   uint32
	TcpConn net.Conn
	Team    uint8
}

func NewUser(team uint8) *User {
	return &User{
		QPort:   0,
		TcpConn: nil,
		Team:    team,
	}
}
