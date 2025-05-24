package internal_types

import "net"

type UdpUserConnStatus struct {
	Conn   *net.UDPAddr
	UserId uint32
}

func NewUdpUserConnStatus(conn *net.UDPAddr, userId uint32) *UdpUserConnStatus {
	return &UdpUserConnStatus{
		Conn:   conn,
		UserId: userId,
	}
}
