package game_manager

import (
	"WeNeedGameServer/external/db"
)

type GameManager struct {
	userDb               *db.Db
	matchScore           uint16 // 게임이 총 몇 판 몇 선제일때의 몇 판
	blueScore            uint16 // 라운드 승리 횟수
	redScore             uint16
	blueAlivePlayerCount uint16
	redAlivePlayerCount  uint16
	finalWinnerTeam      uint8
}

func NewGameManager(userDb *db.Db, matchScore uint16) *GameManager {
	return &GameManager{
		userDb:               userDb,
		matchScore:           matchScore,
		blueScore:            0,
		redScore:             0,
		blueAlivePlayerCount: 0,
		redAlivePlayerCount:  0,
		finalWinnerTeam:      0,
	}
}

func (gm *GameManager) InitAlivePlayer() {
	gm.blueAlivePlayerCount = gm.userDb.GetTeamAlivePlayerCount(db.BlueTeam)
	gm.redAlivePlayerCount = gm.userDb.GetTeamAlivePlayerCount(db.RedTeam)
}

func (gm *GameManager) IncreaseScore(winnerTeam uint8) {
	gm.matchScore -= 1
	if winnerTeam == 'B' {
		gm.blueScore += 1
	} else {
		gm.redScore += 1
	}
}

func (gm *GameManager) DecreasePlayer(deadPlayerTeam uint8) (uint16, uint16) {
	if deadPlayerTeam == 'B' {
		gm.blueAlivePlayerCount -= 1
	} else {
		gm.redAlivePlayerCount -= 1
	}
	return gm.blueAlivePlayerCount, gm.redAlivePlayerCount
}
