package db

import (
	"WeNeedGameServer/util"
	"fmt"
	"net"
	"sync/atomic"
)

type Team bool

const (
	BlueTeam Team = false
	RedTeam       = true
)

type Db struct {
	qPortArr           []uint32
	userList           map[uint32]*User
	RedTeamDb          map[uint32]*User
	BlueTeamDb         map[uint32]*User
	redTeamAliveCount  int64
	blueTeamAliveCount int64
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

func (db *Db) GetTeamAliveCount(team Team) int64 {
	if team == RedTeam {
		return db.redTeamAliveCount
	} else {
		return db.blueTeamAliveCount
	}
}

func (db *Db) DecreaseTeamAliveCount(team Team) {
	if team == RedTeam {
		atomic.AddInt64(&db.redTeamAliveCount, -1)
	} else {
		atomic.AddInt64(&db.blueTeamAliveCount, -1)
	}
}

func (db *Db) IncreaseTeamAliveCount(team Team) {
	if team == RedTeam {
		atomic.AddInt64(&db.redTeamAliveCount, 1)
	} else {
		atomic.AddInt64(&db.blueTeamAliveCount, 1)
	}
}

func (db *Db) Login(userId uint32, userConn net.Conn) error {
	if u, exists := db.userList[userId]; exists {
		u.TcpConn = userConn
		u.QPort = db.qPortArr[len(db.qPortArr)-1]
		db.qPortArr = db.qPortArr[:len(db.qPortArr)-1]
		if u.Team == RedTeam {
			db.RedTeamDb[userId] = u
		} else {
			db.BlueTeamDb[userId] = u
		}
		db.IncreaseTeamAliveCount(u.Team)
		return nil
	} else {
		return fmt.Errorf("login failed")
	}
}

func (db *Db) FindUserByQPort(qPort uint32) uint32 {
	for key, val := range db.userList {
		if val.QPort == qPort {
			return key
		}
	}
	return 0
}

func (db *Db) CheckLogin(userId uint32) bool {
	if db.userList[userId].QPort == 0 || db.userList[userId].TcpConn == nil {
		return false
	} else {
		return true
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

// 더 이상 사용하지 않음(그리고 로직도 명확하지 않음)
//func (db *Db) GetTeamAlivePlayerCount(team Team) uint16 {
//	if team {
//		return uint16(len(db.BlueTeamDb))
//	} else {
//		return uint16(len(db.RedTeamDb))
//	}
//}
