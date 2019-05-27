package chat

import (
	"encoding/json"
	"log"

	"github.com/FlowerWrong/anychat/utils"
)

const (
	WS_LOGIN = iota
	WS_LOGOUT
	WS_SINGLE_CHAT
	WS_GEO
	WS_LAN_IP
	// WS_RE_CONN 掉线重连
	WS_RE_CONN
	// WS_SERVER_PING 服务端发送到客户端的ping消息
	WS_SERVER_PING
	WS_ACK
)

// Base ...
type Base struct {
	Cmd int32  `json:"cmd"`
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

func buildRes(cmd int32, ack string, rawMsg interface{}) ([]byte, error) {
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
