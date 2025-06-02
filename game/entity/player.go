package entity

import "WeNeedGameServer/external/db"

type Player struct {
	RespawnPoint int
	team         db.Team
	isAlive      bool
	damage       int16
	positionX    float32
	xDelta       float32
	positionZ    float32
	zDelta       float32
	yawAngle     float32
	yawDelta     float32
	ptAngle      float32
	ptDelta      float32
	jp           bool

	isShoot             bool
	isReload            bool
	ShootHitInformation map[uint32]int16
}

func NewPlayer(respawnPoint int, team db.Team) *Player {
	return &Player{
		RespawnPoint:        respawnPoint,
		team:                team,
		damage:              0,
		ShootHitInformation: make(map[uint32]int16),
	}
}

func (p *Player) GetPlayerState() *PlayerState {
	return NewPlayerState(p.RespawnPoint, p.team, p.isAlive, p.damage, p.positionX, p.positionZ, p.yawAngle, p.ptAngle, p.jp, p.isShoot, p.isReload)
}

func (p *Player) ReflectPlayer(playerPosition *PlayerState) {
	p.positionX += playerPosition.PositionX
	p.positionZ += playerPosition.PositionZ
	p.damage += playerPosition.Damage
	p.jp = playerPosition.Jp
	p.isShoot = playerPosition.IsShoot
	p.isReload = playerPosition.IsReload
	p.ptAngle += playerPosition.PtAngle
	p.yawAngle += playerPosition.YawAngle
}
