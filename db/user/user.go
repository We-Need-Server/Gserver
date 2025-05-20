package user

import "net"

type User struct {
	QPort   uint32
	TcpAddr *net.TCPAddr
	Team    uint8
}

func NewUser(team uint8) *User {
	return &User{
		QPort:   0,
		TcpAddr: nil,
		Team:    team,
	}
}
