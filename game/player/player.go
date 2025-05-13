package player

type Player struct {
	HP        int16
	HPDelta   int16
	PositionX float32
	XDelta    float32
	PositionZ float32
	ZDelta    float32
	YawAngle  float32
	YawDelta  float32
	PTAngle   float32
	PTDelta   float32
	JP        bool
	IsAlive   bool
}

type PlayerPosition struct {
	HP        int16
	PositionX float32
	PositionZ float32
	YawAngle  float32
	PTAngle   float32
	JP        bool
}

func NewPlayer() *Player {
	return &Player{IsAlive: true}
}

func NewPlayerPosition(hp int16, positionX float32, positionZ float32, yawAngle float32, ptAngle float32, JP bool) PlayerPosition {
	return PlayerPosition{
		HP:        hp,
		PositionX: positionX,
		PositionZ: positionZ,
		YawAngle:  yawAngle,
		PTAngle:   ptAngle,
		JP:        JP,
	}
}

func (p *Player) GetPlayerInfo() *Player {
	return p
}

func (p *Player) ReflectDeltaValues() {
	p.ReflectMoveForward()
	p.ReflectMoveSide()
	p.ReflectTransferPT()
	p.ReflectTransferYaw()
}

func (p *Player) MoveForward(ZDelta float32) {
	p.ZDelta += ZDelta
}

func (p *Player) ReflectMoveForward() {
	p.PositionZ += p.ZDelta
	p.ZDelta = 0
}

func (p *Player) MoveSide(XDelta float32) {
	p.XDelta += XDelta
}

func (p *Player) ReflectMoveSide() {
	p.PositionX += p.XDelta
	p.XDelta = 0
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

func (p *Player) TurnJP(jp bool) {
	p.JP = jp
}

func (p *Player) DamageHP(hpDelta int16) {
	p.HPDelta += hpDelta
}

func (p *Player) ReflectDamageHP() {
	p.HP -= p.HPDelta
	if p.HP <= 0 {
		p.IsAlive = false
	}
	p.HPDelta = 0
}
