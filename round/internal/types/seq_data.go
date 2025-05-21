package types

type SEQData struct {
	QPort uint32
	SEQ   uint32
}

func NewSEQData(qPort uint32, seq uint32) *SEQData {
	return &SEQData{qPort, seq}
}
