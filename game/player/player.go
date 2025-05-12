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
	YawAngle  float32 `json:"yaw"`
	PTAngle   float32 `json:"pitch"`
}

func NewPlayer() *Player {
	return &Player{}
}

func NewPlayerPosition(positionX float32, positionY float32, positionZ float32, yawAngle float32, ptAngle float32) PlayerPosition {
	return PlayerPosition{
		PositionX: positionX,
		PositionY: positionY,
		PositionZ: positionZ,
		YawAngle:  yawAngle,
		PTAngle:   ptAngle,
	}
}

func (p *Player) GetPlayerInfo() *Player {
	return p
}

func (p *Player) MoveForward(XDelta float32) {
	p.XDelta += XDelta
}

func (p *Player) ReflectMoveForward() {
	p.PositionX += p.XDelta
	p.XDelta = 0
}

func (p *Player) MoveSide(ZDelta float32) {
	p.ZDelta += ZDelta
}

func (p *Player) ReflectMoveSide() {
	p.PositionZ += p.ZDelta
	p.ZDelta = 0
}

func (p *Player) TransferYaw(YawDelta float32) {
	p.YawDelta += YawDelta
}

func (p *Player) ReflectTransferYaw() {
	p.YawAngle += p.YawDelta
	p.YawDelta = 0
}

func (p *Player) TransferPT(PTDelta float32) {
	p.PTDelta += PTDelta
}

func (p *Player) ReflectTransferPT() {
	p.PTAngle += p.PTDelta
	p.PTDelta = 0
}
