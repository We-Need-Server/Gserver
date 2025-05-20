package internal_type

type UserSpawnStatus struct {
	userId       uint32
	respawnPoint int16
}

func NewUserSpawnStatus(userId uint32, respawnPoint int16) *UserSpawnStatus {
	return &UserSpawnStatus{
		userId:       userId,
		respawnPoint: respawnPoint,
	}
}
