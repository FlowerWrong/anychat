package chat

import "encoding/json"

const (
	WS_LOGIN = iota
	WS_LOGOUT
	WS_SINGLE_CHAT
	WS_GEO
	WS_LAN_IP
	WS_RE_CONN // 掉线重连
)

type Base struct {
	Cmd int32  `json:"cmd"`
	Ack string `json:"ack"`
}

// Req is common use request body
type Req struct {
	Base
	Body json.RawMessage `json:"body"`
}

// Res is common use response body
type Res struct {
	Base
	Body json.RawMessage `json:"body"`
}

type Ack struct {
	Base
}

type Error struct {
	Base
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

// LoginCmd ...
type LoginCmd struct {
	UserAgent string `json:"userAgent"`
	Domain    string `json:"domain"`
	Token     string `json:"token"`
	Username  string `json:"username"`
	Mobile    string `json:"mobile"`
	Email     string `json:"email"`
}

// LoginRes ...
type LoginRes struct {
	UserID string `json:"userID"`
	ChatID string `json:"chatID"`
}

// SingleChatCmd ...
type SingleChatCmd struct {
	From string `json:"from"`
	To   string `json:"to"`
	Msg  string `json:"msg"`
}

// SingleChatRes ...
type SingleChatRes struct {
	Cmd  int32  `json:"cmd"`
	Ack  string `json:"ack"`
	From string `json:"from"`
	To   string `json:"to"`
	Msg  string `json:"msg"`
	UUID string `json:"uuid"`
}

// LanIPCmd 上报 LAN ip
type LanIPCmd struct {
	LanIP string `json:"lanIP"`
}

// GeoCmd 上报地理位置经纬度
type GeoCmd struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}