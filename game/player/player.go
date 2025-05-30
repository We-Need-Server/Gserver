package player

import (
	"WeNeedGameServer/external/db"
	"fmt"
)

type Player struct {
	RespawnPoint        int
	Team                db.Team
	isAlive             bool
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
	isShoot             bool
	isReload            bool
	ShootHitInformation map[uint32]int16
	decreasePlayerFunc  func(team db.Team)
}

func NewPlayer(respawnPoint int, team db.Team, decreasePlayerFunc func(team db.Team)) *Player {
	return &Player{
		RespawnPoint:        respawnPoint,
		Team:                team,
		hp:                  0,
		isAlive:             true,
		ShootHitInformation: make(map[uint32]int16),
		decreasePlayerFunc:  decreasePlayerFunc,
	}
}

func (p *Player) GetPlayerState() *PlayerPosition {
	return NewPlayerPosition(p.RespawnPoint, p.Team, p.isAlive, p.hp, p.positionX, p.positionZ, p.yawAngle, p.ptAngle, p.jp, p.isShoot, p.isReload)
}

func (p *Player) ReflectDamageHP(hpDelta int16) {
	p.hp += hpDelta
	if p.hp >= 100 {
		p.isAlive = false
		p.hp = 0
	}
}

func (p *Player) ReflectPlayerPosition(playerPosition *PlayerPosition) {
	p.positionX += playerPosition.PositionX
	p.positionZ += playerPosition.PositionZ

	fmt.Println("reflect hp", playerPosition.Hp)
	p.ReflectDamageHP(playerPosition.Hp)
	p.jp = playerPosition.Jp
	p.isShoot = playerPosition.IsShoot
	p.isReload = playerPosition.IsReload
	p.ptAngle += playerPosition.PtAngle
	p.yawAngle += playerPosition.YawAngle
}

// 이게 그러면 game_tick 패킷이 만들어질 때 다 락킹이 걸린다.
// 락킹을 하고 싶지 않아
