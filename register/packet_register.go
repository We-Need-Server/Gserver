package register

import (
	"WeNeedGameServer/packet/udp"
	"fmt"
)

type PacketRegister map[string]udp.PropertyMap

func (p PacketRegister) Add(s string, packetMap udp.PropertyMap) (string, error) {
	if s == "" {
		return "", fmt.Errorf("packet name cannot be empty")
	}

	if packetMap == nil {
		return s, fmt.Errorf("cannot register nil packet map")
	}

	if _, exists := p[s]; !exists {
		p[s] = packetMap
		return s, nil
	}
	return s, fmt.Errorf("a property map with that name as %s already exists", s)
}

func (p PacketRegister) Delete(s string) (string, error) {
	if _, exists := p[s]; exists {
		delete(p, s)
		return s, nil
	}
	return s, fmt.Errorf("not found property map about %s", s)
}
