package bomberman

import (
	"fmt"
	"net/http"
	"reflect"

	game "bomberman/game"
	utils "bomberman/utils"

	"github.com/gorilla/websocket"
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// connect client to socket
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	clientID := utils.GenerateUniqueID()
	client := &utils.Client{
		ID:   clientID,
		Conn: conn,
	}
	//var player *game.Player
	utils.AddClient(client)
	fmt.Printf("Connected client - ID: %s\n", clientID)
	for {
		var data map[string]interface{}
		err = conn.ReadJSON(&data)
		if err != nil {
			// where I detect when someone disconect
			if client.Nickname != "" {
				jsonOtherData := map[string]interface{}{
					"Datatype": "anotherDisconnection",
					"id":       client.Nickname,
				}
				jsonOtherString, err := utils.EncodeJSON(jsonOtherData)
				if err != nil {
					fmt.Println("error encoding JSON:", err)
					return
				}
				utils.BroadcastMessage(jsonOtherString, client)
			}
			utils.RemoveClient(client)
			fmt.Println(err)
			return
		}
		// log.Printf("Recieved data from %s: %v\n", clientID, data)
		switch data["datatype"] {
		case "chat_message":
			jsonOtherData := map[string]interface{}{
				"datatype": "chat_message",
				"data":     data["message"],
				"origin": client.Nickname,
			}
			jsonOtherString, err := utils.EncodeJSON(jsonOtherData)
			if err != nil {
				fmt.Println("error encoding JSON:", err)
				return
			}
			utils.BroadcastMessage(jsonOtherString, client)
		case "nick_setting":
			client.Nickname = fmt.Sprintf("%v", data["message"])
		case "ready_to_play":
			if len(game.Player_list) < 5 && !game.RoomLocked {
				player := &game.Player{
					Number: len(game.Player_list) + 1,
					Lives:  3,
					Range: 2,
					NMBBomb: 1,
					User:   client,
				}
				game.Player_list = append(game.Player_list, player) // add to player list
				jsonOtherData := map[string]interface{}{
					"datatype": "player_number",
					"data":  len(game.Player_list),
				}
				jsonOtherString, err := utils.EncodeJSON(jsonOtherData)
				if err != nil {
					fmt.Println("error encoding JSON:", err)
					return
				}
				utils.BroadcastMessage(jsonOtherString, nil)
				switch len(game.Player_list) {
				case 1:
					player.Position = []int{0, 0}
				case 2:
					player.Position = []int{10, 16}
					go game.Waiting_room(20, "waiting_room")
				case 3:
					if !game.RoomLocked {
						player.Position = []int{0, 16}
					}
				case 4:
					if !game.RoomLocked {
						player.Position = []int{10, 0}
					}
					game.RoomLocked = true
				}
			}
		case "key_press":
			//could also use json unmarshal for this one but I like being disgusting
			m := make(map[string]string)
			tmp := reflect.ValueOf(data["data"])
			for _, k := range tmp.MapKeys() {
				m[k.Interface().(string)] = tmp.MapIndex(k).Interface().(string)
			}
			for i := range game.Player_list {
				if game.Player_list[i].User == client && game.GameStarted && !reflect.DeepEqual(game.Player_list[i].Position, []int{-1, -1}) {
					if m["status"] == "on" {
						if m["key"] == "bomb" {
							game.BombPlanting(game.Player_list[i])
						} else {
							switch m["key"] {
							case "up":
								game.Player_list[i].Directions.Up = true
							case "down":
								game.Player_list[i].Directions.Down = true
							case "left":
								game.Player_list[i].Directions.Left = true
							case "right":
								game.Player_list[i].Directions.Right = true
							}
							if !game.Player_list[i].Moving {
								go game.MoveLoop(game.Player_list[i])
							}
						}
					} else {
						switch m["key"] {
						case "up":
							game.Player_list[i].Directions.Up = false
						case "down":
							game.Player_list[i].Directions.Down = false
						case "left":
							game.Player_list[i].Directions.Left = false
						case "right":
							game.Player_list[i].Directions.Right = false
						}
					}
				}
			}
		}
	}
}

func ServeHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/index.html")
}
