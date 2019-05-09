package chat

import (
	"errors"

	"github.com/FlowerWrong/new_chat/venus/models"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// NewHub ...
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run ...
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// FindClientByUserID ...
func (h *Hub) FindClientByUserID(userID int64) (*Client, error) {
	for client := range h.clients {
		if client.userID == userID {
			return client, nil
		}
	}
	return nil, errors.New("record not found")
}

// FindOnlineUserList ...
func (h *Hub) FindOnlineUserList(users *[]models.User) (onlineUsers []*models.User) {
	for client := range h.clients {
		for _, user := range *users {
			if client.userID == user.Id {
				onlineUsers = append(onlineUsers, &user)
			}
		}
	}
	return onlineUsers
}
