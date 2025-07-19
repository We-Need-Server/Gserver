package tserver

import (
	"WeNeedGameServer/game_type"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type RoundEndPacket struct {
	ContentLength uint32 `json:"-"`
	PKind         uint8  `json:"-"`
	WinnerTeam    uint8  `json:"winnerTeam"`
	RedScore      uint16 `json:"redScore"`
	BlueScore     uint16 `json:"blueScore"`
}

func NewRoundEndPacket(winnerTeam game_type.Team, blueScore uint16, redScore uint16) *RoundEndPacket {
	if winnerTeam == game_type.RedTeam {
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
	p.ContentLength = uint32(len(data) + 1)
	result := make([]byte, 5+len(data))
	binary.LittleEndian.PutUint32(result[0:4], p.ContentLength)
	result[4] = p.PKind
	copy(result[5:], data)
	fmt.Println("E")
	fmt.Println(p.ContentLength)
	fmt.Println(result)
	return result
}
