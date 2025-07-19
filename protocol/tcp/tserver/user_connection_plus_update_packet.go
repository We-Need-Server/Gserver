package tserver

import (
	"WeNeedGameServer/game_type"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type UserTeamStatus struct {
	UserId uint32 `json:"userId"`
	Team   uint8  `json:"team"`
}

func NewUserTeamStatus(userId uint32, team game_type.Team) UserTeamStatus {
	if team == game_type.RedTeam {
		return UserTeamStatus{
			UserId: userId,
			Team:   'R',
		}
	} else {
		return UserTeamStatus{
			UserId: userId,
			Team:   'B',
		}
	}
}

type UserConnectionPUpdatePacket struct {
	ContentLength uint32           `json:"-"`
	PKind         uint8            `json:"-"`
	UserList      []UserTeamStatus `json:"userList"`
}

func NewUserConnectionPUpdatePacket(userList []UserTeamStatus) *UserConnectionPUpdatePacket {
	return &UserConnectionPUpdatePacket{
		PKind:    'P',
		UserList: userList,
	}
}

func (p *UserConnectionPUpdatePacket) Serialize() []byte {
	data, err := json.Marshal(p)
	if err != nil {
		return []byte{}
	}
	p.ContentLength = uint32(len(data) + 1)
	fmt.Println(p.ContentLength)
	result := make([]byte, 5+len(data))
	binary.LittleEndian.PutUint32(result[0:4], p.ContentLength)
	result[4] = p.PKind
	copy(result[5:], data)
	fmt.Println('P')
	fmt.Println(len(result))
	fmt.Println(p.ContentLength)
	fmt.Print(len(data))
	return result
}
