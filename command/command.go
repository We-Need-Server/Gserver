package command

type Command string

type CommandValType interface {
	uint32 | bool | int16
}

const (
	FB Command = "FB" // uint32
	LR         = "LR" // uint32
	YW         = "YW" // uint32
	PT         = "PT" // uint32
	JP         = "JP" // bool
	SH         = "SH" // bool
	HT         = "HT" // uint32, int16이 합쳐져서 총 6바이트를 나눠서 쪼갬
)

type CommandMap map[Command]interface{}

// 커맨드에 따라 기본값으로 초기화하는 함수
func (cm *CommandMap) SetCommand(c Command) {
	// 맵이 nil인 경우 초기화
	if *cm == nil {
		*cm = make(CommandMap)
	}

	// 이미 존재하는 경우 아무 작업도 하지 않음
	if _, exists := (*cm)[c]; exists {
		return
	}

	// 커맨드별 기본값 설정
	switch c {
	case FB, LR, YW, PT:
		// uint32 타입의 기본값은 0
		(*cm)[c] = uint32(0)
	case JP, SH:
		// bool 타입의 기본값은 false
		(*cm)[c] = true
	case HT:
		(*cm)[c] = make(map[uint32]int16)
	}
}

func (cm *CommandMap) AddCommand() {
	knownCommands := []Command{FB, LR, YW, PT, JP, SH, HT}

	for _, cmd := range knownCommands {
		cm.SetCommand(cmd)
	}

}
