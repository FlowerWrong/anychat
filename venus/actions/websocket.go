package actions

import (
	"net/http"

	"github.com/FlowerWrong/new_chat/venus/chat"
)

// WsHandler ...
func WsHandler(hub *chat.Hub, w http.ResponseWriter, r *http.Request) {
	chat.HandleWs(hub, w, r)
}
