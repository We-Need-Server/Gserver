package common

type UserSpawnStatus struct {
	UserId       uint32 `json:"userId"`
	RespawnPoint int16  `json:"respawnPoint"`
}

func NewUserSpawnStatus(userId uint32, respawnPoint int16) *UserSpawnStatus {
	return &UserSpawnStatus{
		UserId:       userId,
		RespawnPoint: respawnPoint,
	}
}
