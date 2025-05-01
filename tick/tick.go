package tick

import (
	"WeNeedGameServer/game"
	"fmt"
	"log"
	"net"
	"time"
)

type GameTick struct {
	TickTime int
	Ticker   *time.Ticker
	Game     *game.Game
	IPTable  *map[uint32]string
	stopChan chan struct{}
}

func NewGameTick(tickTime int, game *game.Game, IPTable *map[uint32]string) *GameTick {
	return &GameTick{
		TickTime: tickTime,
		Ticker:   time.NewTicker(time.Second / time.Duration(tickTime)),
		Game:     game,
		IPTable:  IPTable,
		stopChan: make(chan struct{}),
	}
}

func (gt *GameTick) StartGameLoop() {
	// 루프 시작 틱이 될때마다 processTick함수 실행
	for {
		select {
		case <-gt.Ticker.C:
			gt.processTick()
		case <-gt.stopChan:
			gt.Ticker.Stop()
			return
		}
	}
}

func (gt *GameTick) StopGameLoop() {
	// StartGameLoop가 실행되는 중 어떤 문제가 생길 경우 멈추는 역할
	gt.stopChan <- struct{}{}
}

func (gt *GameTick) processTick() {
	for _, val := range *gt.IPTable {
		clientAddr, err := net.ResolveUDPAddr("udp", val)
		if err != nil {
			log.Panicln("ClientAddr Error:", err)
		}
		fmt.Println(clientAddr)
		clientConn, clientConnErr := net.DialUDP("udp", nil, clientAddr)
		if clientConnErr != nil {
			log.Panicln("Client Connection Error:", clientConnErr)
		}
		defer clientConn.Close() // 연결 종료 보장

		clientConn.Write([]byte("hello"))
	}
	// packetTable 반복문을 통해서 UDP로 getGameState된 결과물 전송
	// packetTable은 Key는 IP Value로 Port와 QPort가 있는데 이를 활용해서 브로드캐스팅
	//gameState :=
	//if err != nil {
	//	return
	//}
	//fmt.Println(gt.PacketTable)
	//for ip, packetVal := range *gt.PacketTable {
	//	// UDP 패킷 전송 (실제 UDP 송신 코드는 network 패키지에 있다고 가정)
	//	network.SendUDPPacket(ip, packetVal, gameState)
	//}
	fmt.Println(gt.Game.GetGameState())
}
