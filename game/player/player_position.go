package player

type PlayerPosition struct {
	Hp        int16
	PositionX float32
	PositionZ float32
	YawAngle  float32
	PtAngle   float32
	Jp        bool
	IsShoot   bool
	IsReload  bool
}

func NewPlayerPositionD() *PlayerPosition {
	return &PlayerPosition{
		Hp:        0,     // 기본 체력
		PositionX: 0.0,   // 기본 X 위치
		PositionZ: 0.0,   // 기본 Z 위치
		YawAngle:  0.0,   // 기본 요 각도
		PtAngle:   0.0,   // 기본 피치 각도
		Jp:        false, // 기본 점프 상태
		IsShoot:   false, // 기본 발사 상태
		IsReload:  false,
	}
}

func (p *PlayerPosition) CalculatePlayerPosition(calP *PlayerPosition) {
	p.PositionX += calP.PositionX
	p.PositionZ += calP.PositionZ
	p.Hp += calP.Hp
	p.PtAngle += calP.PtAngle
	p.YawAngle += calP.YawAngle
	p.IsShoot = p.IsShoot || calP.IsShoot
	p.IsShoot = p.IsShoot || calP.IsShoot
	p.IsReload = p.IsReload || calP.IsReload
}

func NewPlayerPosition(hp int16, positionX float32, positionZ float32, yawAngle float32, ptAngle float32, jp bool, isShoot bool, isReload bool) *PlayerPosition {
	return &PlayerPosition{
		Hp:        hp,
		PositionX: positionX,
		PositionZ: positionZ,
		YawAngle:  yawAngle,
		PtAngle:   ptAngle,
		Jp:        jp,
		IsShoot:   isShoot,
		IsReload:  isReload,
	}
}
