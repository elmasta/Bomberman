package bomberman

import (
	"fmt"
	"net/http"

	// "log"

	h "bomberman/handlers"
)

func Runserver() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../frontend/static"))))
	http.Handle("/framework/", http.StripPrefix("/framework/", http.FileServer(http.Dir("../frontend/framework"))))
	http.Handle("/game/", http.StripPrefix("/game/", http.FileServer(http.Dir("../frontend/game"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../frontend/assets"))))

	http.HandleFunc("/ws", h.HandleWebSocket)
	http.HandleFunc("/", h.ServeHTML)

	port := 8080
	fmt.Printf("Server listening on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println(err)
	}
}
