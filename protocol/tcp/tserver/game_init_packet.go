package tserver

import (
	"WeNeedGameServer/game_type"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type GameInitPacket struct {
	ContentLength      uint32                       `json:"-"`
	PKind              uint8                        `json:"-"`
	TickNumber         uint32                       `json:"tickNumber"`
	BlueScore          uint16                       `json:"blueScore"`
	RedScore           uint16                       `json:"redScore"`
	UserSpawnStatusArr []*game_type.UserSpawnStatus `json:"userSpawnStatusArr"`
}

func NewGameInitPacket(tickNumber uint32, blueScore uint16, redScore uint16, userSpawnStatusArr []*game_type.UserSpawnStatus) *GameInitPacket {
	return &GameInitPacket{
		PKind:              'R',
		TickNumber:         tickNumber,
		BlueScore:          blueScore,
		RedScore:           redScore,
		UserSpawnStatusArr: userSpawnStatusArr,
	}
}

func (p *GameInitPacket) Serialize() []byte {
	data, err := json.Marshal(p)
	if err != nil {
		return []byte{}
	}
	p.ContentLength = uint32(len(data) + 1)
	result := make([]byte, 5+len(data))
	binary.LittleEndian.PutUint32(result[0:4], p.ContentLength)
	result[4] = p.PKind
	copy(result[5:], data)
	fmt.Println('R')
	fmt.Println(len(result))
	fmt.Println(p.ContentLength)
	fmt.Print(len(data))
	return result
}
