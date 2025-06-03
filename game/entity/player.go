package entity

import "WeNeedGameServer/game_type"

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
	}
}

func (p *Player) GetPlayerState() *game_type.PlayerState {
	return game_type.NewPlayerState(p.RespawnPoint, p.team, p.isAlive, p.damage, p.positionX, p.positionZ, p.yawAngle, p.ptAngle, p.jp, p.isShoot, p.isReload)
}

func (p *Player) ReflectPlayer(playerPosition *game_type.PlayerState) {
	p.positionX += playerPosition.PositionX
	p.positionZ += playerPosition.PositionZ
	p.damage += playerPosition.Damage
	p.jp = playerPosition.Jp
	p.isShoot = playerPosition.IsShoot
	p.isReload = playerPosition.IsReload
	p.ptAngle += playerPosition.PtAngle
	p.yawAngle += playerPosition.YawAngle
}
