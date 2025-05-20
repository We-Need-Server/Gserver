package receiver

type SendingType uint8

const (
	SendByBroadCast SendingType = iota
	SendByUniCast
)

type TcpReceiverMessage struct {
	// 브로드 캐스팅 할지 여부
	SenderType SendingType
	// 어떤 유저한테 보낼지 결정
	UserId uint32
	// 어떤 메시지를 보내야 할지 결정
	PKind uint8
}

func NewBroadCastMessage(pKind uint8) *TcpReceiverMessage {
	return &TcpReceiverMessage{
		SenderType: SendByBroadCast,
		UserId:     0,
		PKind:      pKind,
	}
}

func NewUniCastMessage(userId uint32, pKind uint8) *TcpReceiverMessage {
	return &TcpReceiverMessage{
		SenderType: SendByUniCast,
		UserId:     userId,
		PKind:      pKind,
	}
}
