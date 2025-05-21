package legacy

//
//import (
//	"WeNeedGameServer/game/player"
//	"fmt"
//)
//
//types Game struct {
//	Players map[uint32]*player.Player
//}
//
//func NewGame() *Game {
//	return &Game{Players: make(map[uint32]*player.Player)}
//}
//
//func (g *Game) AddPlayer(QPort uint32) *player.Player {
//	g.Players[QPort] = player.NewPlayer()
//	return g.Players[QPort]
//}
//
//func (g *Game) GetGameDeltaState() map[uint32]player.PlayerPosition {
//	gameDeltaState := make(map[uint32]player.PlayerPosition)
//	for qPort, p := range g.Players {
//		for key, val := range p.ShootHitInformation {
//			g.Players[key].DamageHP(val)
//		}
//		gameDeltaState[qPort] = p.GetPlayerDeltaState()
//		fmt.Println(*gameDeltaState[qPort].Hp)
//		p.ReflectDeltaValues()
//	}
//	return gameDeltaState
//}
//
//func (g *Game) GetGameState() map[uint32]player.PlayerPosition {
//	gameState := make(map[uint32]player.PlayerPosition)
//	for qPort, p := range g.Players {
//		gameState[qPort] = p.GetPlayerState()
//	}
//	return gameState
//}
//
//func (g *Game) ResetHPDelta() {
//	for _, p := range g.Players {
//		p.ReflectDamageHP()
//	}
//}
