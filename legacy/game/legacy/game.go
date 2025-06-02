package legacy

//
//import (
//	"WeNeedGameServer/game/entity"
//	"fmt"
//)
//
//internal_types Game struct {
//	Players map[uint32]*entity.Player
//}
//
//func NewGame() *Game {
//	return &Game{Players: make(map[uint32]*entity.Player)}
//}
//
//func (g *Game) AddPlayer(QPort uint32) *entity.Player {
//	g.Players[QPort] = entity.NewPlayer()
//	return g.Players[QPort]
//}
//
//func (g *Game) GetGameDeltaState() map[uint32]entity.PlayerPosition {
//	gameDeltaState := make(map[uint32]entity.PlayerPosition)
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
//func (g *Game) GetGameState() map[uint32]entity.PlayerPosition {
//	gameState := make(map[uint32]entity.PlayerPosition)
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
