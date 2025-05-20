package db

import (
	"WeNeedGameServer/db/user"
	"WeNeedGameServer/util"
	"fmt"
	"net"
)

type UserDb struct {
	qPortArr   []uint32
	userList   map[uint32]*user.User
	RedTeamDb  map[uint32]*user.User
	BlueTeamDb map[uint32]*user.User
}

func NewUserDb() *UserDb {
	qPortArr := make([]uint32, 128)
	for i := 1; i <= 128; i++ {
		qPortArr[i-1] = uint32(i)
	}
	util.ShuffleUint32Arr(qPortArr)
	return &UserDb{
		qPortArr:   qPortArr,
		userList:   make(map[uint32]*user.User),
		RedTeamDb:  make(map[uint32]*user.User),
		BlueTeamDb: make(map[uint32]*user.User),
	}
}

func (db *UserDb) Init() {
	db.AddUser(16, 'B')
	db.AddUser(32, 'B')
	db.AddUser(64, 'B')
	db.AddUser(128, 'B')
	db.AddUser(256, 'B')
	db.AddUser(8, 'R')
	db.AddUser(24, 'R')
	db.AddUser(48, 'R')
	db.AddUser(96, 'R')
	db.AddUser(192, 'R')
}

func (db *UserDb) AddUser(userId uint32, team uint8) {
	db.userList[userId] = user.NewUser(team)
}

func (db *UserDb) Login(userId uint32, userAddr *net.TCPAddr) (uint32, error) {
	if user, exists := db.userList[userId]; exists {
		user.TcpAddr = userAddr
		user.QPort = db.qPortArr[len(db.qPortArr)-1]
		db.qPortArr = db.qPortArr[:len(db.qPortArr)-1]
		return user.QPort, nil
	} else {
		return 0, fmt.Errorf("login failed")
	}
}
