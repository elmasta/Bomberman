package bomberman

import (
	"fmt"
	"math/rand"
	"time"
	"strconv"

	utils "bomberman/utils"
)

type Player struct {
	Number			int
	Lives			int
	Range			int
	NMBBomb			int
	Directions		Direction
	Position		[]int
	Moving			bool
	Invulnerable	int
	User			*utils.Client
}

type Direction struct {
	Up 		bool
	Left	bool
	Right	bool
	Down	bool
}

type Bomb struct {
	Player		*Player
	Position	[]int
	Timer		int
}

type ToBroadcast struct {
	PlayerNMB		int
	PlayerPosition	[]int
	PlayerLives		int
}

type ChangingCoor struct {
	Y 			int
	X 			int
	SquareType	string
	Extra		string
}

var Player_list []*Player
var bomb_list []Bomb
var flame_list [][]int
var Terrain [][]string
var RoomLocked bool
var GameStarted bool

func GenerateTerrain() {

	//puting the bricks
	for y := 1; y < 12; y++ {
		var tmp []string
		if (y >= 1 && y <= 2) || (y >= 10 && y <= 11) {
			for x := 1; x < 18; x++ {
				if (x >= 1 && x <= 2) || (x >= 16 && x <= 17) {
					tmp = append(tmp, ".")
				} else {
					if rand.Intn(2) == 0 {
						tmp = append(tmp, ".")
					} else {
						tmp = append(tmp, "B")
					}
				}
			}
		} else {
			for x := 1; x < 18; x++ {
				if rand.Intn(2) == 0 {
					tmp = append(tmp, ".")
				} else {
					tmp = append(tmp, "B")
				}
			}
		}
		Terrain = append(Terrain, tmp)
	}

	//puting the indestructible briques
	for y := range Terrain {
		if y % 2 != 0 {
			for x := range Terrain[y] {
				if x % 2 != 0 {
					Terrain[y][x] = "X"
				}
			}
		}
	}

	//puting the power ups
	generatePowerUps("r", 3)
	generatePowerUps("m", 3)
	generatePowerUps("l", 1)

	//debug
	for _, v := range Terrain {
		fmt.Println(v)
	}
}

func generatePowerUps(powerUpType string, amount int) {

	for amount > 0 {
		for {
			x := rand.Intn(17)
			y := rand.Intn(11)
			if Terrain[y][x] == "B" {
				Terrain[y][x] = powerUpType
				break
			}
		}
		amount--
	}
}

func startGame() {

	GenerateTerrain()

	// remove spoilers for players
	tmp := [][]string{}
	for _, v := range Terrain {
		stmp := make([]string, len(Terrain[0]))
		copy(stmp, v)
		tmp = append(tmp, stmp)
	}
	for y := range tmp {
		for x := range tmp[y] {
			if tmp[y][x] == "r" || tmp[y][x] == "l" || tmp[y][x] == "m" {
				tmp[y][x] = "B"
			}
		}
	}

	jsonOtherData := map[string]interface{}{
		"datatype": "game_grid",
		"data":     tmp,
		"players":	Player_list,
	}
	jsonOtherString, err := utils.EncodeJSON(jsonOtherData)
	if err != nil {
		fmt.Println("error encoding JSON:", err)
		return
	}
	utils.BroadcastMessage(jsonOtherString, nil)

	GameStarted = true
	gameLoop()
}

func gameLoop() {
	var dead int
	for {
		time.Sleep(16 * time.Millisecond)
		var NewToBroadcast []ChangingCoor
		var flameToKeep [][]int
		for i := range flame_list { //remove old flames
			if flame_list[i][2] == 0 {
				Terrain[flame_list[i][0]][flame_list[i][1]] = "."
				NewToBroadcast = append(NewToBroadcast, ChangingCoor{flame_list[i][0], flame_list[i][1], ".", ""})
			} else {
				flame_list[i][2]--
				flameToKeep = append(flameToKeep, flame_list[i])
			}
		}
		flame_list = flameToKeep
		var bombToRemove []int
		for i := range bomb_list {
			if bomb_list[i].Timer == 0 {
				alreadyDone := false
				for _, b := range bombToRemove {
					if b == i {
						alreadyDone = true
					}
				} 
				if !alreadyDone {
					tmp1, tmp2 := BombExplode(NewToBroadcast, bombToRemove, i, "")
					NewToBroadcast = append(NewToBroadcast, tmp1...)
					bombToRemove = append(bombToRemove, tmp2...)
				}				
			} else {
				bomb_list[i].Timer--
			}
		}
		if len(bombToRemove) > 0 {
			var tmp []Bomb
			for b := range bomb_list {
				for i, v := range bombToRemove { //remove exploded bombs
					if b == v {
						break
					} else if i == len(bombToRemove)-1 {
						tmp = append(tmp, bomb_list[b])
					}
				}
			}
			bomb_list = tmp
		}
		for i := range Player_list {
			if Player_list[i].Position[0] != -1 {
				if Player_list[i].Invulnerable > 0 {
					Player_list[i].Invulnerable--
				}
				if Terrain[Player_list[i].Position[0]][Player_list[i].Position[1]] == "F" && Player_list[i].Invulnerable == 0 { //check player position to make them get hit
					Player_list[i].Lives--
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{Player_list[i].Lives, 0, "H", Player_list[i].User.Nickname})
					Player_list[i].Invulnerable = 110
				}
				if Player_list[i].Lives == 0 {
					dead++
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{Player_list[i].Position[0], Player_list[i].Position[1], "D", strconv.Itoa(Player_list[i].Number)})
					Player_list[i].Position = []int{-1, -1} //put the player out of the game
				}
			}
		}
		if len(NewToBroadcast) > 0 || dead >= len(Player_list)-1 {
			jsonOtherData := make(map[string]interface{})
			jsonOtherData["datatype"] = "loop_result"
			jsonOtherData["data"] = NewToBroadcast
			jsonOtherString, err := utils.EncodeJSON(jsonOtherData)
			if err != nil {
				fmt.Println("error encoding JSON:", err)
				return
			}
			utils.BroadcastMessage(jsonOtherString, nil)
			if dead >= len(Player_list)-1 {
				jsonOtherData["datatype"] = "loop_result"
				jsonOtherData["data"] = "endgame"
				jsonOtherString, err := utils.EncodeJSON(jsonOtherData)
				if err != nil {
					fmt.Println("error encoding JSON:", err)
					return
				}
				utils.BroadcastMessage(jsonOtherString, nil)
				break
			}
		}
	}
	GameStarted = false
}

func MoveLoop(player *Player) {
	player.Moving = true
	for player.Moving{
		if (player.Lives == 0){
			return
		}
		switch {
		case player.Directions.Up:
			if player.Position[0]-1 >= 0 {
				//check if there's no other player
				switch Terrain[player.Position[0]-1][player.Position[1]] {
					//check if damage taken
				case ".":
					player.Position[0] = player.Position[0]-1
					toBroadcast := ToBroadcast{player.Number, player.Position, player.Lives}
					broadcastMoves(toBroadcast, "player_move")
					time.Sleep(500 * time.Millisecond)
				case "F":
					player.Position[0] = player.Position[0]-1
					if player.Invulnerable == 0 {
						player.Lives--
						player.Invulnerable = 110
						broadcastMoves(ChangingCoor{player.Lives, 0, "H", player.User.Nickname}, "loop_result")
					}
					toBroadcast := ToBroadcast{player.Number, player.Position, player.Lives}
					broadcastMoves(toBroadcast, "player_move")
					time.Sleep(500 * time.Millisecond)
				case "L":
					player.Position[0] = player.Position[0]-1
					toBroadcast := ToBroadcast{player.Number, player.Position, player.Lives}
					broadcastMoves(toBroadcast, "player_move")
					broadcastMoves([]int{player.Position[0], player.Position[1]}, "remove_powerup")
										
					player.Lives++
					broadcastMoves(ChangingCoor{player.Lives, 1, "H", player.User.Nickname}, "loop_result")
					time.Sleep(500 * time.Millisecond)
				case "R":
					player.Position[0] = player.Position[0]-1
					toBroadcast := ToBroadcast{player.Number, player.Position, player.Lives}
					broadcastMoves(toBroadcast, "player_move")
					broadcastMoves([]int{player.Position[0], player.Position[1]}, "remove_powerup")
					player.Range++
					time.Sleep(500 * time.Millisecond)
				case "M":
					player.Position[0] = player.Position[0]-1
					toBroadcast := ToBroadcast{player.Number, player.Position, player.Lives}
					broadcastMoves(toBroadcast, "player_move")
					broadcastMoves([]int{player.Position[0], player.Position[1]}, "remove_powerup")
					player.NMBBomb++
					time.Sleep(500 * time.Millisecond)
				}
			}
		case player.Directions.Down:
			if player.Position[0]+1 <= 10 {
				switch Terrain[player.Position[0]+1][player.Position[1]] {
				case ".":
					player.Position[0] = player.Position[0]+1
					toBroadcast := ToBroadcast{player.Number, player.Position, player.Lives}
					broadcastMoves(toBroadcast, "player_move")
					time.Sleep(500 * time.Millisecond)
				case "F":
					player.Position[0] = player.Position[0]+1
					if player.Invulnerable == 0 {
						player.Lives--
						player.Invulnerable = 110
						broadcastMoves(ChangingCoor{player.Lives, 0, "H", player.User.Nickname}, "loop_result")
					}
					toBroadcast := ToBroadcast{player.Number, player.Position, player.Lives}
					broadcastMoves(toBroadcast, "player_move")
					time.Sleep(500 * time.Millisecond)
				case "L":
					player.Position[0] = player.Position[0]+1
					toBroadcast := ToBroadcast{player.Number, player.Position, player.Lives}
					broadcastMoves(toBroadcast, "player_move")
					broadcastMoves([]int{player.Position[0], player.Position[1]}, "remove_powerup")
					player.Lives++
					broadcastMoves(ChangingCoor{player.Lives, 1, "H", player.User.Nickname}, "loop_result")
					time.Sleep(500 * time.Millisecond)
				case "R":
					player.Position[0] = player.Position[0]+1
					toBroadcast := ToBroadcast{player.Number, player.Position, player.Lives}
					broadcastMoves(toBroadcast, "player_move")
					broadcastMoves([]int{player.Position[0], player.Position[1]}, "remove_powerup")
					player.Range++
					time.Sleep(500 * time.Millisecond)
				case "M":
					player.Position[0] = player.Position[0]+1
					toBroadcast := ToBroadcast{player.Number, player.Position, player.Lives}
					broadcastMoves(toBroadcast, "player_move")
					broadcastMoves([]int{player.Position[0], player.Position[1]}, "remove_powerup")
					player.NMBBomb++
					time.Sleep(500 * time.Millisecond)
				}
			}
		case player.Directions.Left:
			if player.Position[1]-1 >= 0 {
				switch Terrain[player.Position[0]][player.Position[1]-1] {
				case ".":
					player.Position[1] = player.Position[1]-1
					toBroadcast := ToBroadcast{player.Number, player.Position, player.Lives}
					broadcastMoves(toBroadcast, "player_move")
					time.Sleep(500 * time.Millisecond)
				case "F":
					player.Position[1] = player.Position[1]-1
					if player.Invulnerable == 0 {
						player.Lives--
						player.Invulnerable = 110
						broadcastMoves(ChangingCoor{player.Lives, 0, "H", player.User.Nickname}, "loop_result")
					}
					toBroadcast := ToBroadcast{player.Number, player.Position, player.Lives}
					broadcastMoves(toBroadcast, "player_move")
					time.Sleep(500 * time.Millisecond)
				case "L":
					player.Position[1] = player.Position[1]-1
					toBroadcast := ToBroadcast{player.Number, player.Position, player.Lives}
					broadcastMoves(toBroadcast, "player_move")
					broadcastMoves([]int{player.Position[0], player.Position[1]}, "remove_powerup")
					player.Lives++
					broadcastMoves(ChangingCoor{player.Lives, 1, "H", player.User.Nickname}, "loop_result")
					time.Sleep(500 * time.Millisecond)
				case "R":
					player.Position[1] = player.Position[1]-1
					toBroadcast := ToBroadcast{player.Number, player.Position, player.Lives}
					broadcastMoves(toBroadcast, "player_move")
					broadcastMoves([]int{player.Position[0], player.Position[1]}, "remove_powerup")
					player.Range++
					time.Sleep(500 * time.Millisecond)
				case "M":
					player.Position[1] = player.Position[1]-1
					toBroadcast := ToBroadcast{player.Number, player.Position, player.Lives}
					broadcastMoves(toBroadcast, "player_move")
					broadcastMoves([]int{player.Position[0], player.Position[1]}, "remove_powerup")
					player.NMBBomb++
					time.Sleep(500 * time.Millisecond)
				}
			}
		case player.Directions.Right:
			if player.Position[1]+1 <= 16 {
				switch Terrain[player.Position[0]][player.Position[1]+1] {
				case ".":
					player.Position[1] = player.Position[1]+1
					toBroadcast := ToBroadcast{player.Number, player.Position, player.Lives}
					broadcastMoves(toBroadcast, "player_move")
					time.Sleep(500 * time.Millisecond)
				case "F":
					player.Position[1] = player.Position[1]+1
					if player.Invulnerable == 0 {
						player.Lives--
						player.Invulnerable = 110
						broadcastMoves(ChangingCoor{player.Lives, 0, "H", player.User.Nickname}, "loop_result")
					}
					toBroadcast := ToBroadcast{player.Number, player.Position, player.Lives}
					broadcastMoves(toBroadcast, "player_move")
					time.Sleep(500 * time.Millisecond)
				case "L":
					player.Position[1] = player.Position[1]+1
					toBroadcast := ToBroadcast{player.Number, player.Position, player.Lives}
					broadcastMoves(toBroadcast, "player_move")
					broadcastMoves([]int{player.Position[0], player.Position[1]}, "remove_powerup")
					player.Lives++
					broadcastMoves(ChangingCoor{player.Lives, 1, "H", player.User.Nickname}, "loop_result")
					time.Sleep(500 * time.Millisecond)
				case "R":
					player.Position[1] = player.Position[1]+1
					toBroadcast := ToBroadcast{player.Number, player.Position, player.Lives}
					broadcastMoves(toBroadcast, "player_move")
					broadcastMoves([]int{player.Position[0], player.Position[1]}, "remove_powerup")
					player.Range++
					time.Sleep(500 * time.Millisecond)
				case "M":
					player.Position[1] = player.Position[1]+1
					toBroadcast := ToBroadcast{player.Number, player.Position, player.Lives}
					broadcastMoves(toBroadcast, "player_move")
					broadcastMoves([]int{player.Position[0], player.Position[1]}, "remove_powerup")
					player.NMBBomb++
					time.Sleep(500 * time.Millisecond)
				}
			}
		default:
			player.Moving = false		
		}
	}
}

func BombExplode(NewToBroadcast []ChangingCoor, bombToRemove []int, i int, directionToIgnore string) ([]ChangingCoor, []int) {
	bombToRemove = append(bombToRemove, i)
	flame_list = append(flame_list, []int{bomb_list[i].Position[0], bomb_list[i].Position[1], 100})
	Terrain[bomb_list[i].Position[0]][bomb_list[i].Position[1]] = "F"
	NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0], bomb_list[i].Position[1], "F", ""})
	
	if directionToIgnore != "-y" {
		//flame -y
		for r := 1; r <= bomb_list[i].Player.Range; r++ {
			if bomb_list[i].Position[0]-r < 0 {
				break
			}
			var done bool
			switch Terrain[bomb_list[i].Position[0]-r][bomb_list[i].Position[1]] {
				case "X":
					done = true
				case "B":
					Terrain[bomb_list[i].Position[0]-r][bomb_list[i].Position[1]] = "."
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0]-r, bomb_list[i].Position[1], ".", ""})
					done = true
				case "r":
					Terrain[bomb_list[i].Position[0]-r][bomb_list[i].Position[1]] = "R"
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0]-r, bomb_list[i].Position[1], "R", ""})
					done = true
				case "l":
					Terrain[bomb_list[i].Position[0]-r][bomb_list[i].Position[1]] = "L"
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0]-r, bomb_list[i].Position[1], "L", ""})
					done = true
				case "m":
					Terrain[bomb_list[i].Position[0]-r][bomb_list[i].Position[1]] = "M"
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0]-r, bomb_list[i].Position[1], "M", ""})
					done = true
				case "E":
					var index int
					for ind, v := range bomb_list {
						if v.Position[0] == bomb_list[i].Position[0]-r && v.Position[1] == bomb_list[i].Position[1] {
							index = ind
							break
						}
					}
					tmp1, tmp2 := BombExplode(NewToBroadcast, bombToRemove, index, "+y")
					NewToBroadcast = append(NewToBroadcast, tmp1...)
					bombToRemove = append(bombToRemove, tmp2...)
					done = true
				default:
					Terrain[bomb_list[i].Position[0]-r][bomb_list[i].Position[1]] = "F"
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0]-r, bomb_list[i].Position[1], "F", ""})
					flame_list = append(flame_list, []int{bomb_list[i].Position[0]-r, bomb_list[i].Position[1], 100})
			}
			if done {
				break
			}
		}
	}
	if directionToIgnore != "+y" {
		//flame +y
		for r := 1; r <= bomb_list[i].Player.Range; r++ {
			if bomb_list[i].Position[0]+r > 10 {
				break
			}
			var done bool
			switch Terrain[bomb_list[i].Position[0]+r][bomb_list[i].Position[1]] {
				case "X":
					done = true
				case "B":
					Terrain[bomb_list[i].Position[0]+r][bomb_list[i].Position[1]] = "."
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0]+r, bomb_list[i].Position[1], ".", ""})
					done = true
				case "r":
					Terrain[bomb_list[i].Position[0]+r][bomb_list[i].Position[1]] = "R"
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0]+r, bomb_list[i].Position[1], "R", ""})
					done = true
				case "l":
					Terrain[bomb_list[i].Position[0]+r][bomb_list[i].Position[1]] = "L"
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0]+r, bomb_list[i].Position[1], "L", ""})
					done = true
				case "m":
					Terrain[bomb_list[i].Position[0]+r][bomb_list[i].Position[1]] = "M"
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0]+r, bomb_list[i].Position[1], "M", ""})
					done = true
				case "E":
					var index int
					for ind, v := range bomb_list {
						if v.Position[0] == bomb_list[i].Position[0]+r && v.Position[1] == bomb_list[i].Position[1] {
							index = ind
							break
						}
					}
					tmp1, tmp2 := BombExplode(NewToBroadcast, bombToRemove, index, "-y")
					NewToBroadcast = append(NewToBroadcast, tmp1...)
					bombToRemove = append(bombToRemove, tmp2...)
					done = true
				default:
					Terrain[bomb_list[i].Position[0]+r][bomb_list[i].Position[1]] = "F"
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0]+r, bomb_list[i].Position[1], "F", ""})
					flame_list = append(flame_list, []int{bomb_list[i].Position[0]+r, bomb_list[i].Position[1], 100})
			}
			if done {
				break
			}
		}
	}
	if directionToIgnore != "-x" {
		//flame -x
		for r := 1; r <= bomb_list[i].Player.Range; r++ {
			if bomb_list[i].Position[1]-r < 0 {
				break
			}
			var done bool
			switch Terrain[bomb_list[i].Position[0]][bomb_list[i].Position[1]-r] {
				case "X":
					done = true
				case "B":
					Terrain[bomb_list[i].Position[0]][bomb_list[i].Position[1]-r] = "."
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0], bomb_list[i].Position[1]-r, ".", ""})
					done = true
				case "r":
					Terrain[bomb_list[i].Position[0]][bomb_list[i].Position[1]-r] = "R"
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0], bomb_list[i].Position[1]-r, "R", ""})
					done = true
				case "l":
					Terrain[bomb_list[i].Position[0]][bomb_list[i].Position[1]-r] = "L"
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0], bomb_list[i].Position[1]-r, "L", ""})
					done = true
				case "m":
					Terrain[bomb_list[i].Position[0]][bomb_list[i].Position[1]-r] = "M"
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0], bomb_list[i].Position[1]-r, "M", ""})
					done = true
				case "E":
					var index int
					for ind, v := range bomb_list {
						if v.Position[0] == bomb_list[i].Position[0] && v.Position[1] == bomb_list[i].Position[1]-r {
							index = ind
							break
						}
					}
					tmp1, tmp2 := BombExplode(NewToBroadcast, bombToRemove, index, "+x")
					NewToBroadcast = append(NewToBroadcast, tmp1...)
					bombToRemove = append(bombToRemove, tmp2...)
					done = true
				default:
					Terrain[bomb_list[i].Position[0]][bomb_list[i].Position[1]-r] = "F"
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0], bomb_list[i].Position[1]-r, "F", ""})
					flame_list = append(flame_list, []int{bomb_list[i].Position[0], bomb_list[i].Position[1]-r, 100})
			}
			if done {
				break
			}
		}
	}
	if directionToIgnore != "+x" {
		//flame +x
		for r := 1; r <= bomb_list[i].Player.Range; r++ {
			if bomb_list[i].Position[1]+r > 16 {
				break
			}
			var done bool
			switch Terrain[bomb_list[i].Position[0]][bomb_list[i].Position[1]+r] {
				case "X":
					done = true
				case "B":
					Terrain[bomb_list[i].Position[0]][bomb_list[i].Position[1]+r] = "."
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0], bomb_list[i].Position[1]+r, ".", ""})
					done = true
				case "r":
					Terrain[bomb_list[i].Position[0]][bomb_list[i].Position[1]+r] = "R"
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0], bomb_list[i].Position[1]+r, "R", ""})
					done = true
				case "l":
					Terrain[bomb_list[i].Position[0]][bomb_list[i].Position[1]+r] = "L"
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0], bomb_list[i].Position[1]+r, "L", ""})
					done = true
				case "m":
					Terrain[bomb_list[i].Position[0]][bomb_list[i].Position[1]+r] = "M"
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0], bomb_list[i].Position[1]+r, "M", ""})
					done = true
				case "E":
					var index int
					for ind, v := range bomb_list {
						if v.Position[0] == bomb_list[i].Position[0] && v.Position[1] == bomb_list[i].Position[1]+r {
							index = ind
							break
						}
					}
					tmp1, tmp2 := BombExplode(NewToBroadcast, bombToRemove, index, "-x")
					NewToBroadcast = append(NewToBroadcast, tmp1...)
					bombToRemove = append(bombToRemove, tmp2...)
					done = true
				default:
					Terrain[bomb_list[i].Position[0]][bomb_list[i].Position[1]+r] = "F"
					NewToBroadcast = append(NewToBroadcast, ChangingCoor{bomb_list[i].Position[0], bomb_list[i].Position[1]+r, "F", ""})
					flame_list = append(flame_list, []int{bomb_list[i].Position[0], bomb_list[i].Position[1]+r, 100})
			}
			if done {
				break
			}
		}
	}
	return NewToBroadcast, bombToRemove
}

func BombPlanting(player *Player) {
	if Terrain[player.Position[0]][player.Position[1]] != "F" {
		//check if player can put bomb
		count := 0
		for _, v := range bomb_list {
			if v.Player == player {
				count++
			}
		}
		if count < player.NMBBomb {
			bomb_list = append(bomb_list, Bomb{Player: player, Position: []int{player.Position[0], player.Position[1]}, Timer: 240})
			Terrain[player.Position[0]][player.Position[1]] = "E"
			jsonOtherData := map[string]interface{}{
				"datatype": "bomb",
				"data": []int{player.Position[0], player.Position[1]} ,
			}
			jsonOtherString, err := utils.EncodeJSON(jsonOtherData)
			if err != nil {
				fmt.Println("error encoding JSON:", err)
				return
			}
			utils.BroadcastMessage(jsonOtherString, nil)
		}
	}
}

func broadcastMoves(toBroadcast interface{}, datatype string) {
	jsonOtherData := map[string]interface{}{
		"datatype": datatype,
		"data":  toBroadcast,
	}
	jsonOtherString, err := utils.EncodeJSON(jsonOtherData)
	if err != nil {
		fmt.Println("error encoding JSON:", err)
		return
	}
	utils.BroadcastMessage(jsonOtherString, nil)
}
