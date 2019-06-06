package actions

import (
	"net/http"

	"github.com/FlowerWrong/anychat/chat"
)

// WebsocketHandler ...
func WebsocketHandler(hub *chat.Hub, w http.ResponseWriter, r *http.Request) {
	chat.HandleWs(hub, w, r)
}
