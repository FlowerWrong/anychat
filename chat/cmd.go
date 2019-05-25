package chat

import "encoding/json"

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
	// WS_PONG 服务端接收客户端的pong消息
	WS_PONG
)

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
	UserAgent string `json:"userAgent"`
	Domain    string `json:"domain"`
	Token     string `json:"token"`
}

// LoginRes ...
type LoginRes struct {
	UserID string `json:"userID"`
}

// SingleChatCmd ...
type SingleChatCmd struct {
	From string `json:"from"`
	To   string `json:"to"`
	Msg  string `json:"msg"`
}

// SingleChatRes ...
type SingleChatRes struct {
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

// PingCmd ...
type PingCmd struct {
	PingAt interface{} `json:"ping_at"`
}
