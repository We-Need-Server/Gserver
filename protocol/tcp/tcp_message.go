package tcp

type SendingType uint8

const (
	SendByBroadCast SendingType = iota
	SendByUniCast
	SendByMultiCast
)

type Message struct {
	// 브로드 캐스팅 할지 여부
	SenderType SendingType
	// 어떤 유저한테 보낼지 결정
	UserId uint32
	// 데이터
	Data PacketI
}

func NewBroadCastMessage(data PacketI) *Message {
	return &Message{
		SenderType: SendByBroadCast,
		UserId:     0,
		Data:       data,
	}
}

func NewUniCastMessage(userId uint32, data PacketI) *Message {
	return &Message{
		SenderType: SendByUniCast,
		UserId:     userId,
		Data:       data,
	}
}
