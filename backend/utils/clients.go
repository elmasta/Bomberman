package bomberman

import (
	"sync"
	"fmt"

	"github.com/gorilla/websocket"
)

var (
	clients     = make(map[string]*Client)
	clientsLock sync.Mutex
)

type Client struct {
	ID   		string
	Nickname 	string
	Conn 		*websocket.Conn
}

func BroadcastMessage(message string, brodcaster_client *Client) {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	for _, client := range clients {
		err := client.Conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Printf("Error sending message to - ID: %s, error: %v\n", client.ID, err)
		}
	}
}

func AddClient(client *Client) {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	clients[client.ID] = client
}

func RemoveClient(client *Client) {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	delete(clients, client.ID)
}
