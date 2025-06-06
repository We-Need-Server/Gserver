package db

import (
	"WeNeedGameServer/game_type"
	"WeNeedGameServer/util"
	"fmt"
	"net"
	"sync/atomic"
)

type UserDb struct {
	qPortArr           []uint32
	userList           map[uint32]*User
	RedTeamDb          map[uint32]*User
	BlueTeamDb         map[uint32]*User
	redTeamAliveCount  int64
	blueTeamAliveCount int64
}

func NewUserDb() *UserDb {
	qPortArr := make([]uint32, 128)
	for i := 1; i <= 128; i++ {
		qPortArr[i-1] = uint32(i)
	}
	util.ShuffleUint32Arr(qPortArr)
	return &UserDb{
		qPortArr:   qPortArr,
		userList:   make(map[uint32]*User), // 회원 가입 한 유저들
		RedTeamDb:  make(map[uint32]*User), // 현재 레드팀에서 활성화된 유저들 => 2명
		BlueTeamDb: make(map[uint32]*User), // 현재 블루팀에서 활성화된 유저들 => 1명
	}
}

func (db *UserDb) Init() {
	db.AddUser(16, game_type.BlueTeam)
	db.AddUser(32, game_type.BlueTeam)
	db.AddUser(64, game_type.BlueTeam)
	db.AddUser(128, game_type.BlueTeam)
	db.AddUser(256, game_type.BlueTeam)
	db.AddUser(8, game_type.RedTeam)
	db.AddUser(24, game_type.RedTeam)
	db.AddUser(48, game_type.RedTeam)
	db.AddUser(96, game_type.RedTeam)
	db.AddUser(192, game_type.RedTeam)
}

func (db *UserDb) AddUser(userId uint32, team game_type.Team) {
	db.userList[userId] = NewUser(team)
}

func (db *UserDb) GetTeamAliveCount(team game_type.Team) int64 {
	if team == game_type.RedTeam {
		return db.redTeamAliveCount
	} else {
		return db.blueTeamAliveCount
	}
}

func (db *UserDb) ResetTeamAliveCount() {
	db.blueTeamAliveCount = int64(len(db.BlueTeamDb))
	db.redTeamAliveCount = int64(len(db.RedTeamDb))
}

func (db *UserDb) DecreaseTeamAliveCount(team game_type.Team) {
	if team == game_type.RedTeam {
		atomic.AddInt64(&db.redTeamAliveCount, -1)
	} else {
		atomic.AddInt64(&db.blueTeamAliveCount, -1)
	}
}

func (db *UserDb) IncreaseTeamAliveCount(team game_type.Team) {
	if team == game_type.RedTeam {
		atomic.AddInt64(&db.redTeamAliveCount, 1)
	} else {
		atomic.AddInt64(&db.blueTeamAliveCount, 1)
	}
}

func (db *UserDb) Login(userId uint32, userConn net.Conn) (uint32, game_type.Team, error) {
	if u, exists := db.userList[userId]; exists {
		u.TcpConn = userConn
		u.QPort = db.qPortArr[len(db.qPortArr)-1]
		db.qPortArr = db.qPortArr[:len(db.qPortArr)-1]
		if u.Team == game_type.RedTeam {
			db.RedTeamDb[userId] = u
		} else {
			db.BlueTeamDb[userId] = u
		}
		db.IncreaseTeamAliveCount(u.Team)
		return u.QPort, u.Team, nil
	} else {
		return 0, false, fmt.Errorf("login failed")
	}
}

func (db *UserDb) FindUserByQPort(qPort uint32) uint32 {
	for key, val := range db.userList {
		if val.QPort == qPort {
			return key
		}
	}
	return 0
}

func (db *UserDb) CheckLogin(userId uint32) bool {
	if db.userList[userId].QPort == 0 || db.userList[userId].TcpConn == nil {
		return false
	} else {
		return true
	}
}

func (db *UserDb) ResetUser(userId uint32, team game_type.Team) {
	db.userList[userId].QPort = 0
	db.userList[userId].TcpConn = nil
	if team {
		delete(db.RedTeamDb, userId)
	} else {
		delete(db.BlueTeamDb, userId)
	}
}

// 더 이상 사용하지 않음(그리고 로직도 명확하지 않음)
//func (db *UserDb) GetTeamAlivePlayerCount(team Team) uint16 {
//	if team {
//		return uint16(len(db.BlueTeamDb))
//	} else {
//		return uint16(len(db.RedTeamDb))
//	}
//}
