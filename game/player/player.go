package player

type Player struct {
	PositionX float32
	XDelta    float32
	PositionY float32
	YDelta    float32
	PositionZ float32
	ZDelta    float32
	YawAngle  float32
	YawDelta  float32
	PTAngle   float32
	PTDelta   float32
}

type PlayerPosition struct {
	PositionX float32 `json:"positionX"`
	PositionY float32 `json:"positionY"`
	PositionZ float32 `json:"positionZ"`
	YawAngle  float32 `json:"yawAngle"`
	PTAngle   float32 `json:"ptAngle"`
}

func NewPlayer() *Player {
	return &Player{}
}

func (p Player) GetPlayerInfo() Player {
	return p
}

func (p *Player) MoveFoward(XDelta float32) {
	// 좌표에 값을 넣어두려면 계산 공식을 써야함
	// p.PositionX += xDelta
	p.XDelta += XDelta
}

func (p *Player) MoveSide(ZDelta float32) {
	// 좌표에 값을 넣어두려면 계산 공식을 써야함
	// p.PositionZ += ZDelta
	p.ZDelta += ZDelta
}

func (p *Player) TransferYaw(YawDelta float32) {
	// 좌표에 값을 넣어두려면 계산 공식을 써야함
	// p.YawAngle += yawDelta
	p.YawDelta += YawDelta
}

func (p *Player) TransferPT(PTDelta float32) {
	// 좌표에 값을 넣어두려면 계산 공식을 써야함
	// p.PTAngle += ptDelta
	p.PTDelta += PTDelta
}
