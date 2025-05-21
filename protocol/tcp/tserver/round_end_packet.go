package tserver

type RoundEndPacket struct {
	PKind      uint8  `json:"packetKind"`
	WinnerTeam uint8  `json:"winnerTeam"`
	RedScore   uint16 `json:"redScore"`
	BlueScore  uint16 `json:"blueScore"`
}

func NewRoundEndPacket(winnerTeam uint8, blueScore uint16, redScore uint16) *RoundEndPacket {
	return &RoundEndPacket{
		PKind:      'E',
		WinnerTeam: winnerTeam,
		BlueScore:  blueScore,
		RedScore:   redScore,
	}
}

func (p *RoundEndPacket) DeSerialize() {

}

func (p *RoundEndPacket) Serialize() []byte {
	return []byte{}
}
