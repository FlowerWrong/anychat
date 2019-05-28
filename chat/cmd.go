package chat

import (
	"encoding/json"
	"log"

	"github.com/FlowerWrong/anychat/utils"
	"github.com/FlowerWrong/util"
)

const (
	TypeWelcome    = "welcome"
	TypeLogin      = "login"
	TypeDisconnect = "disconnect"
	TypeReConn     = "re_conn"
	TypePing       = "ping"
	TypeAck        = "ack"
	TypeGeo        = "geo"
	TypeLanIP      = "lan_ip"
	TypeSingleChat = "single_chat"
	TypeRoomChat   = "room_chat"
)

// Base ...
type Base struct {
	Cmd string `json:"cmd"`
	Ack string `json:"ack"`
}

// Req is common use request body
type Req struct {
	Base
	Data json.RawMessage `json:"data"`
}

// Res is common use response body
type Res struct {
	Base
	Data json.RawMessage `json:"data"`
}

// ErrorRes ...
type ErrorRes struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

// LoginCmd ...
type LoginCmd struct {
	UserAgent string `json:"user_agent"`
	Domain    string `json:"domain"`
	Token     string `json:"token"`
}

// LoginRes ...
type LoginRes struct {
	UserID string `json:"user_id"`
}

// DisconnectCmd ...
type DisconnectCmd struct {
	Reason    string `json:"reason"`
	Reconnect bool   `json:"reconnect"`
}

// SingleChatCmd ...
type SingleChatCmd struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Msg       string `json:"msg"`
	CreatedAt int64  `json:"created_at"` // 纳秒
}

// SingleChatRes ...
type SingleChatRes struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Msg       string `json:"msg"`
	UUID      string `json:"uuid"`
	CreatedAt int64  `json:"created_at"` // 纳秒
}

// LanIPCmd 上报 LAN ip
type LanIPCmd struct {
	LanIP string `json:"lan_ip"`
}

// GeoCmd 上报地理位置经纬度
type GeoCmd struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// PingCmd ...
type PingCmd struct {
	PingAt interface{} `json:"ping_at"`
}

// WelcomeCmd ...
type WelcomeCmd struct {
	Message string `json:"message"`
}

// Ack ...
type Ack struct {
	Action string `json:"action"`
}

// RoomChatCmd ...
type RoomChatCmd struct {
	From      string `json:"from"` // user uuid
	To        string `json:"to"`   // room uuid
	Msg       string `json:"msg"`
	CreatedAt int64  `json:"created_at"` // 纳秒
}

// RoomChatRes ...
type RoomChatRes struct {
	From      string `json:"from"` // user uuid
	To        string `json:"to"`   // room uuid
	Msg       string `json:"msg"`
	UUID      string `json:"uuid"`       // user_room_message uuid
	CreatedAt int64  `json:"created_at"` // 纳秒
}

func buildRes(cmd, ack string, rawMsg interface{}) ([]byte, error) {
	raw, err := utils.RawMsg(rawMsg)
	if err != nil {
		return nil, err
	}
	res := Res{Base: Base{Cmd: cmd, Ack: ack}, Data: raw}
	data, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return data, nil
}

func (c *Client) sendDisconnectRes(reason string, reconnect bool) error {
	data, err := buildRes(TypeDisconnect, util.UUID(), DisconnectCmd{Reason: reason, Reconnect: reconnect})
	if err != nil {
		return err
	}
	c.send <- data
	return nil
}

func (c *Client) sendAckRes(ack, action string) error {
	data, err := buildRes(TypeAck, ack, Ack{Action: action})
	if err != nil {
		return err
	}
	c.send <- data
	return nil
}
