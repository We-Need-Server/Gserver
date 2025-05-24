package tserver

import (
	"WeNeedGameServer/common"
	"encoding/json"
)

type GameInitPacket struct {
	PKind              uint32                    `json:"packetKind"`
	TickNumber         uint32                    `json:"tickNumber"`
	BlueScore          uint16                    `json:"blueScore"`
	RedScore           uint16                    `json:"redScore"`
	UserSpawnStatusArr []*common.UserSpawnStatus `json:"userSpawnStatusArr"`
}

func NewGameInitPacket(tickNumber uint32, blueScore uint16, redScore uint16, userSpawnStatusArr []*common.UserSpawnStatus) *GameInitPacket {
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
	return data
}
