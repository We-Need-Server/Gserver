package round

type Round struct {
	RedAlivePlayerCount  uint16
	BlueAlivePlayerCount uint16
	totalRound           uint16
	redScore             uint16
	blueScore            uint16
	finalWinnerTeam      uint8
}

func NewRound(redAlivePlayerCount uint16, blueAlivePlayerCount uint16, totalRound uint16) *Round {
	return &Round{
		RedAlivePlayerCount:  redAlivePlayerCount,
		BlueAlivePlayerCount: blueAlivePlayerCount,
		redScore:             0,
		blueScore:            0,
		finalWinnerTeam:      0,
		totalRound:           totalRound,
	}
}

func (r *Round) IncreaseScore(winnerTeam uint8) uint8 {
	if winnerTeam == 'B' {
		r.blueScore += 1
	} else {
		r.redScore += 1
	}
	r.totalRound -= 1
	if r.totalRound == 0 {
		if r.blueScore > r.redScore {
			r.finalWinnerTeam = 'B'
		} else {
			r.finalWinnerTeam = 'R'
		}
	}
	return r.finalWinnerTeam
}

func (r *Round) DecreasePlayer(deadPlayerTeam uint8) (uint16, uint16) {
	if deadPlayerTeam == 'B' {
		r.BlueAlivePlayerCount -= 1
	} else {
		r.RedAlivePlayerCount -= 1
	}
	return r.BlueAlivePlayerCount, r.RedAlivePlayerCount
}
