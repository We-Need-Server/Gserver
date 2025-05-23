package player

type Player struct {
	respawnPoint        int
	hp                  int16
	hpDelta             int16
	positionX           float32
	xDelta              float32
	positionZ           float32
	zDelta              float32
	yawAngle            float32
	yawDelta            float32
	ptAngle             float32
	ptDelta             float32
	jp                  bool
	isAlive             bool
	isShoot             bool
	isReload            bool
	ShootHitInformation map[uint32]int16
}

func NewPlayer(respawnPoint int) *Player {
	return &Player{
		respawnPoint:        respawnPoint,
		hp:                  100,
		ShootHitInformation: make(map[uint32]int16),
	}
}

//	func (p *Player) GetPlayerDeltaState() PlayerPosition {
//		return NewPlayerPosition(p.hpDelta, p.xDelta, p.zDelta, p.yawDelta, p.ptDelta, p.jp, p.isShoot)
//	}
func (p *Player) GetPlayerState() *PlayerPosition {
	return NewPlayerPosition(p.hp, p.positionX, p.positionZ, p.yawAngle, p.ptAngle, p.jp, p.isShoot, p.isReload)
}

func (p *Player) ReflectDeltaValues() {
	p.ReflectMoveForward()
	p.ReflectMoveSide()
	p.ReflectTransferPT()
	p.ReflectTransferYaw()
	p.ReflectHitInformation()
	p.ReflectIsShoot()
}

func (p *Player) MoveForward(zDelta float32) {
	p.zDelta += zDelta
}

func (p *Player) ReflectMoveForward() {
	p.positionZ += p.zDelta
	p.zDelta = 0
}

func (p *Player) MoveSide(xDelta float32) {
	p.xDelta += xDelta
}

func (p *Player) ReflectMoveSide() {
	p.positionX += p.xDelta
	p.xDelta = 0
}

func (p *Player) TransferYaw(yawDelta float32) {
	p.yawDelta += yawDelta
}

func (p *Player) ReflectTransferYaw() {
	p.yawAngle += p.yawDelta
	p.yawDelta = 0
}

func (p *Player) TransferPT(ptDelta float32) {
	p.ptDelta += ptDelta
}

func (p *Player) ReflectTransferPT() {
	p.ptAngle += p.ptDelta
	p.ptDelta = 0
}

func (p *Player) TurnJP(jp bool) {
	p.jp = jp
}

func (p *Player) DamageHP(hpDelta int16) {
	p.hpDelta += hpDelta
}

func (p *Player) ReflectDamageHP() {
	p.hp -= p.hpDelta
	if p.hp <= 0 {
		p.isAlive = false
	}
	p.hpDelta = 0
}

// false
// true

func (p *Player) TurnIsShoot() {
	p.isShoot = true
}

func (p *Player) ReflectIsShoot() {
	p.isShoot = false
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

func (p *Player) ReflectPlayerPosition(playerPosition *PlayerPosition) {
	p.positionX += (*playerPosition).PositionX
	p.positionZ += (*playerPosition).PositionZ
	p.hp -= (*playerPosition).Hp
	p.jp = (*playerPosition).Jp
	p.isShoot = (*playerPosition).IsShoot
	p.isReload = playerPosition.IsReload
	p.ptAngle += (*playerPosition).PtAngle
	p.yawAngle += (*playerPosition).YawAngle
}

// 이게 그러면 game_tick 패킷이 만들어질 때 다 락킹이 걸린다.
// 락킹을 하고 싶지 않아
