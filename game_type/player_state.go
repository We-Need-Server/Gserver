package game_type

type PlayerState struct {
	RespawnPoint int
	Team         Team
	IsAlive      bool
	Damage       int16
	PositionX    float32
	PositionZ    float32
	YawAngle     float32
	PtAngle      float32
	Jp           bool
	IsShoot      bool
	IsReload     bool
}

func NewPlayerStateDefault() *PlayerState {
	return &PlayerState{
		RespawnPoint: 0,
		Team:         false,
		Damage:       0,     // 기본 체력
		PositionX:    0.0,   // 기본 X 위치
		PositionZ:    0.0,   // 기본 Z 위치
		YawAngle:     0.0,   // 기본 요 각도
		PtAngle:      0.0,   // 기본 피치 각도
		Jp:           false, // 기본 점프 상태
		IsShoot:      false, // 기본 발사 상태
		IsReload:     false,
		IsAlive:      false,
	}
}

func (p *PlayerState) CalculatePlayerState(calP *PlayerState) {
	p.RespawnPoint = calP.RespawnPoint
	p.Team = calP.Team
	p.IsAlive = p.IsAlive && calP.IsAlive
	p.PositionX += calP.PositionX
	p.PositionZ += calP.PositionZ
	p.Damage += calP.Damage
	p.PtAngle += calP.PtAngle
	p.YawAngle += calP.YawAngle
	p.Jp = p.Jp || calP.Jp
	p.IsShoot = p.IsShoot || calP.IsShoot
	p.IsReload = p.IsReload || calP.IsReload

}

func NewPlayerState(respawnPoint int, team Team, isAlive bool, hp int16, positionX float32, positionZ float32, yawAngle float32, ptAngle float32, jp bool, isShoot bool, isReload bool) *PlayerState {
	return &PlayerState{
		RespawnPoint: respawnPoint,
		Team:         team,
		IsAlive:      isAlive,
		Damage:       hp,
		PositionX:    positionX,
		PositionZ:    positionZ,
		YawAngle:     yawAngle,
		PtAngle:      ptAngle,
		Jp:           jp,
		IsShoot:      isShoot,
		IsReload:     isReload,
	}
}
