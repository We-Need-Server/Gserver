package db

import (
	"WeNeedGameServer/util"
	"fmt"
	"net"
)

type Team bool

const (
	BlueTeam Team = false
	RedTeam       = true
)

type Db struct {
	qPortArr   []uint32
	userList   map[uint32]*User
	RedTeamDb  map[uint32]*User
	BlueTeamDb map[uint32]*User
}

func NewUserDb() *Db {
	qPortArr := make([]uint32, 128)
	for i := 1; i <= 128; i++ {
		qPortArr[i-1] = uint32(i)
	}
	util.ShuffleUint32Arr(qPortArr)
	return &Db{

		qPortArr:   qPortArr,
		userList:   make(map[uint32]*User), // 회원 가입 한 유저들
		RedTeamDb:  make(map[uint32]*User), // 현재 레드팀에서 활성화된 유저들 => 2명
		BlueTeamDb: make(map[uint32]*User), // 현재 블루팀에서 활성화된 유저들 => 1명
	}
}

func (db *Db) Init() {
	db.AddUser(16, BlueTeam)
	db.AddUser(32, BlueTeam)
	db.AddUser(64, BlueTeam)
	db.AddUser(128, BlueTeam)
	db.AddUser(256, BlueTeam)
	db.AddUser(8, RedTeam)
	db.AddUser(24, RedTeam)
	db.AddUser(48, RedTeam)
	db.AddUser(96, RedTeam)
	db.AddUser(192, RedTeam)
}

func (db *Db) AddUser(userId uint32, team Team) {
	db.userList[userId] = NewUser(team)
}

func (db *Db) Login(userId uint32, userConn net.Conn) error {
	if u, exists := db.userList[userId]; exists {
		u.TcpConn = userConn
		u.QPort = db.qPortArr[len(db.qPortArr)-1]
		db.qPortArr = db.qPortArr[:len(db.qPortArr)-1]
		if u.Team {
			db.RedTeamDb[userId] = u
		} else {
			db.BlueTeamDb[userId] = u
		}
		return nil
	} else {
		return fmt.Errorf("login failed")
	}
}

func (db *Db) ResetUser(userId uint32, team Team) {
	db.userList[userId].QPort = 0
	db.userList[userId].TcpConn = nil
	if team {
		delete(db.RedTeamDb, userId)
	} else {
		delete(db.BlueTeamDb, userId)
	}
}

func (db *Db) GetTeamAlivePlayerCount(team Team) uint16 {
	if team {
		return uint16(len(db.BlueTeamDb))
	} else {
		return uint16(len(db.RedTeamDb))
	}
}
