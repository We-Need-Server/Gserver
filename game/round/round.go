package round

type Round struct {
	redScore        uint16
	blueScore       uint16
	finalWinnerTeam uint8
	totalRound      uint16
}

func NewRound(totalRound uint16) *Round {
	return &Round{
		redScore:        0,
		blueScore:       0,
		finalWinnerTeam: 0,
		totalRound:      totalRound,
	}
}

func (r *Round) IncreaseRedTeam(winnerTeam uint8) uint8 {
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
