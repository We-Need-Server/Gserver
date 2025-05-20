package tcp

type PacketI interface {
	DeSerialize()
	Serialize() []byte
}
