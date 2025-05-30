package tserver

import (
	"WeNeedGameServer/external/db"
	"encoding/json"
)

type RoundEndPacket struct {
	PKind      uint8  `json:"-"`
	WinnerTeam uint8  `json:"winnerTeam"`
	RedScore   uint16 `json:"redScore"`
	BlueScore  uint16 `json:"blueScore"`
}

func NewRoundEndPacket(winnerTeam db.Team, blueScore uint16, redScore uint16) *RoundEndPacket {
	if winnerTeam == db.RedTeam {
		return &RoundEndPacket{
			PKind:      'E',
			WinnerTeam: 'R',
			BlueScore:  blueScore,
			RedScore:   redScore,
		}
	} else {
		return &RoundEndPacket{
			PKind:      'E',
			WinnerTeam: 'B',
			BlueScore:  blueScore,
			RedScore:   redScore,
		}
	}
}

func (p *RoundEndPacket) Serialize() []byte {
	data, err := json.Marshal(p)
	if err != nil {
		return []byte{}
	}
	result := make([]byte, 1+len(data))
	result[0] = p.PKind
	copy(result[1:], data)

	return result
}
