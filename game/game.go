package game

import (
	"WeNeedGameServer/game/entity"
	"WeNeedGameServer/game_type"
	"fmt"
)

type Game struct {
	decreasePlayerFunc func(team game_type.Team)
	playerList         map[uint32]*entity.Player
}

func NewGame(decreasePlayerFunc func(team game_type.Team)) *Game {
	return &Game{
		decreasePlayerFunc: decreasePlayerFunc,
		playerList:         make(map[uint32]*entity.Player),
	}
}

func (g *Game) ReadyGame(blueTeamUserIdList []uint32, redTeamUserIdList []uint32, userSpawnPositionList []int) *Game {
	// 블루팀 스폰
	for index, userId := range blueTeamUserIdList {
		fmt.Println("add player", userId)
		g.addPlayer(userId, userSpawnPositionList[index], game_type.BlueTeam)
	}
	// 레드팀 스폰
	for index, userId := range redTeamUserIdList {
		fmt.Println("add player", userId)
		g.addPlayer(userId, userSpawnPositionList[index], game_type.RedTeam)
	}
	return g
}

func (g *Game) GetGameState() map[uint32]*game_type.PlayerState {
	gameState := make(map[uint32]*game_type.PlayerState)
	for userId, p := range g.playerList {
		gameState[userId] = p.GetPlayerState()
		if !gameState[userId].IsAlive {
			g.DeletePlayer(userId)
			g.decreasePlayerFunc(gameState[userId].Team)
		}
	}
	return gameState
}

func (g *Game) GetPlayerSpawnStatusList() []*game_type.UserSpawnStatus {
	var userSpawnStatusArr []*game_type.UserSpawnStatus
	for key, val := range g.playerList {
		userSpawnStatusArr = append(userSpawnStatusArr, game_type.NewUserSpawnStatus(key, int16(val.RespawnPoint)))
	}
	return userSpawnStatusArr
}

func (g *Game) addPlayer(userId uint32, respawnPosition int, team game_type.Team) {
	g.playerList[userId] = entity.NewPlayer(respawnPosition, team)
	fmt.Println(g.playerList[userId])
}

func (g *Game) DeletePlayer(userId uint32) {
	fmt.Println("플레이어 제거")
	delete(g.playerList, userId)
}

func (g *Game) ReflectPlayers(playerStateMap map[uint32]*game_type.PlayerState) {
	fmt.Printf("자세한 맵2: %+v\n", playerStateMap)
	for key, val := range playerStateMap {
		if _, exists := g.playerList[key]; exists {
			fmt.Println("reflect player", key, val.Damage)
			g.playerList[key].ReflectPlayer(val)
		} else {
			fmt.Println("not reflect player", key)
		}
	}
}
