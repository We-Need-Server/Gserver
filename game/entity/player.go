package entity

import (
	"WeNeedGameServer/game_type"
	"fmt"
)

type Player struct {
	RespawnPoint int
	team         game_type.Team
	isAlive      bool
	damage       int16
	positionX    float32
	positionZ    float32
	yawAngle     float32
	ptAngle      float32
	jp           bool
	isShoot      bool
	isReload     bool
}

func NewPlayer(respawnPoint int, team game_type.Team) *Player {
	return &Player{
		RespawnPoint: respawnPoint,
		team:         team,
		isAlive:      true,
	}
}

func (p *Player) GetPlayerState() *game_type.PlayerState {
	return game_type.NewPlayerState(p.RespawnPoint, p.team, p.isAlive, p.damage, p.positionX, p.positionZ, p.yawAngle, p.ptAngle, p.jp, p.isShoot, p.isReload)
}

func (p *Player) ReflectPlayer(playerState *game_type.PlayerState) {
	p.positionX += playerState.PositionX
	p.positionZ += playerState.PositionZ
	p.damage += playerState.Damage
	fmt.Println("player damage", p.damage, playerState.Damage)
	if p.damage >= 100 {
		p.isAlive = false
	}
	p.jp = playerState.Jp
	p.isShoot = playerState.IsShoot
	p.isReload = playerState.IsReload
	p.ptAngle += playerState.PtAngle
	p.yawAngle += playerState.YawAngle
}
