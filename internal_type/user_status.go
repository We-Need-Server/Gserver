package internal_type

import "WeNeedGameServer/game/player"

type UserSpawnStatus struct {
	UserId            uint32                            `json:"userId"`
	RespawnPoint      int16                             `json:"respawnPoint"`
	PlayerPositionMap map[uint32]*player.PlayerPosition `json:"playerPositionMap"`
}

func NewUserSpawnStatus(userId uint32, respawnPoint int16, playerPositionMap map[uint32]*player.PlayerPosition) *UserSpawnStatus {
	return &UserSpawnStatus{
		UserId:            userId,
		RespawnPoint:      respawnPoint,
		PlayerPositionMap: playerPositionMap,
	}
}
