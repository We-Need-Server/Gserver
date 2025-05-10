package register

import (
	"WeNeedGameServer/packet"
	"fmt"
)

type PacketRegister struct {
	packetRegistry map[string]packet.PropertyMap
}

func NewPacketRegister() *PacketRegister {
	return &PacketRegister{packetRegistry: make(map[string]packet.PropertyMap)}
}

func (p *PacketRegister) Register(s string, packetMap packet.PropertyMap) (string, error) {
	if s == "" {
		return "", fmt.Errorf("packet name cannot be empty")
	}

	if packetMap == nil {
		return s, fmt.Errorf("cannot register nil packet map")
	}

	if _, exists := p.packetRegistry[s]; !exists {
		p.packetRegistry[s] = packetMap
		return s, nil
	}
	return s, fmt.Errorf("a property map with that name as %s already exists", s)
}

func (p *PacketRegister) Delete(s string) (string, error) {
	if _, exists := p.packetRegistry[s]; exists {
		delete(p.packetRegistry, s)
		return s, nil
	}
	return s, fmt.Errorf("not found property map about %s", s)
}
