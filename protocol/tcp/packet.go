package tcp

type PacketI interface {
	Serialize() []byte
}
