package tserver

import (
	"WeNeedGameServer/external/db"
	"encoding/json"
)

type UserTeamStatus struct {
	UserId uint32 `json:"userId"`
	Team   uint8  `json:"team"`
}

func NewUserTeamStatus(userId uint32, team db.Team) UserTeamStatus {
	if team == db.RedTeam {
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
	PKind    uint8            `json:"-"`
	UserList []UserTeamStatus `json:"userList"`
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
	result := make([]byte, 1+len(data))
	result[0] = p.PKind
	copy(result[1:], data)

	return result
}
