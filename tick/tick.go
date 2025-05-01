package tick

import (
	"WeNeedGameServer/game"
	"WeNeedGameServer/network"
	"fmt"
	"log"
	"time"
)

type GameTick struct {
	TickTime        int
	Ticker          *time.Ticker
	Game            *game.Game
	networkInstance *network.Network
	stopChan        chan struct{}
}

func NewGameTick(tickTime int, game *game.Game, networkInstance *network.Network) *GameTick {
	return &GameTick{
		TickTime:        tickTime,
		Ticker:          time.NewTicker(time.Second / time.Duration(tickTime)),
		Game:            game,
		networkInstance: networkInstance,
		stopChan:        make(chan struct{}),
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
	// 클라이언트 연결 맵을 유지

	// 모든 클라이언트에 대한 연결 초기화
	//for _, val := range *gt.IPTable {
	//	clientAddr, err := net.ResolveUDPAddr("udp", val)
	//	if err != nil {
	//		log.Println("ClientAddr Error for", val, ":", err)
	//		continue // 오류 발생 시 패닉하지 않고 다음 클라이언트로 진행
	//	}
	//
	//	clientConn, clientConnErr := net.DialUDP("udp", nil, clientAddr)
	//	if clientConnErr != nil {
	//		log.Println("Client Connection Error for", val, ":", clientConnErr)
	//		continue // 오류 발생 시 패닉하지 않고 다음 클라이언트로 진행
	//	}
	//
	//	// 연결을 맵에 저장
	//	clientConns[val] = clientConn
	//}

	gameState := gt.Game.GetGameState()
	//fmt.Println(*gt.ConnTable)
	for key, userAddr := range gt.networkInstance.ConnTable {
		fmt.Println(key, userAddr)
		_, err := gt.networkInstance.Ln.WriteToUDP([]byte(gameState), userAddr)
		if err != nil {
			log.Println("Failed to send message:", err)
		}
	}

	fmt.Println("Game state sent to", len(gt.networkInstance.ConnTable), "clients")
}
