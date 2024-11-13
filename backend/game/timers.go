package bomberman

import (
	"fmt"
	"time"

	utils "bomberman/utils"
)

func Waiting_room(timer int, timerType string) {
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	for {
		select {
			case <- ticker.C:
				//broadcast
				jsonOtherData := map[string]interface{}{
					"datatype": timerType,
					"data":  timer,
					"nbPlayer": len(Player_list),
				}
				jsonOtherString, err := utils.EncodeJSON(jsonOtherData)
				if err != nil {
					fmt.Println("error encoding JSON:", err)
					return
				}
				utils.BroadcastMessage(jsonOtherString, nil)
				timer--
				if timer == 0 || (RoomLocked && timerType != "game_starting") {
					RoomLocked = true
					close(quit)
				}
				if len(Player_list) < 2 {
					close(quit)
				}
			case <- quit:
				ticker.Stop()
				if RoomLocked && timerType == "game_starting" {
					startGame()
				} else if len(Player_list) >= 2 {
					Waiting_room(10, "game_starting")
				}
				return
		}
	}
}
