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

	IsShoot             bool
	ShootHitInformation map[uint32]int16
}

type HitInformationField struct {
	HP      int16
	HPDelta int16
}

type PlayerPosition struct {
	HP        *int16
	PositionX float32
	PositionZ float32
	YawAngle  float32
	PTAngle   float32
	JP        bool
	IsShoot   bool
}

func NewPlayer() *Player {
	return &Player{HP: 100, ShootHitInformation: make(map[uint32]int16)}
}

func NewPlayerPosition(hp *int16, positionX float32, positionZ float32, yawAngle float32, ptAngle float32, JP bool, isShoot bool) PlayerPosition {
	return PlayerPosition{
		HP:        hp,
		PositionX: positionX,
		PositionZ: positionZ,
		YawAngle:  yawAngle,
		PTAngle:   ptAngle,
		JP:        JP,
		IsShoot:   isShoot,
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
	p.ReflectHitInformation()
	p.ReflectIsShoot()
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
	//if p.HP <= 0 {
	//	p.IsAlive = false
	//}
	p.HPDelta = 0
}

func (p *Player) TurnIsShoot() {
	p.IsShoot = true
}

func (p *Player) ReflectIsShoot() {
	p.IsShoot = false
}

func (p *Player) StoreHitInformation(qPort uint32, hpDelta int16) {
	// 키가 존재하는지 확인
	p.ShootHitInformation[qPort] += hpDelta
	//if _, exists := p.ShootHitInformation[qPort]; exists {
	//
	//} else {
	//	p.ShootHitInformation[qPort] = HitInformationField{HPDelta: hpDelta}
	//}
}

func (p *Player) ReflectHitInformation() {
	for key, _ := range p.ShootHitInformation {
		p.ShootHitInformation[key] = 0
	}
}
