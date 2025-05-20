package round

type Round struct {
	totalRound      uint16 // 게임이 총 몇 판 몇 선제일때의 몇 판
	redScore        uint16
	blueScore       uint16 // 라운드 승리 횟수
	finalWinnerTeam uint8  // 9판 최종 우승팀
}

func NewRound(totalRound uint16) *Round {
	return &Round{
		redScore:        0,
		blueScore:       0,
		finalWinnerTeam: 0,
		totalRound:      totalRound,
	}
}

func (r *Round) IncreaseScore(winnerTeam uint8) {
	r.totalRound -= 1
	if winnerTeam == 'B' {
		r.blueScore += 1
	} else {
		r.redScore += 1
	}
}

//func (r *Round) DecreasePlayer(deadPlayerTeam uint8) (uint16, uint16) {
//	if deadPlayerTeam == 'B' {
//		r.BlueAlivePlayerCount -= 1
//	} else {
//		r.RedAlivePlayerCount -= 1
//	}
//	return r.BlueAlivePlayerCount, r.RedAlivePlayerCount
//}
