package chat

const (
	WS_LOGIN = iota
	WS_LOGOUT
	WS_SINGLE_CHAT
)

type Base struct {
	Cmd int32  `json:"cmd"`
	Ack string `json:"ack"`
}

type Req struct {
	Base
	Body []byte `json:"body"`
}

type Res struct {
	Base
	Body []byte `json:"body"`
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
// navigator.userAgent
// local ip https://jsfiddle.net/ourcodeworld/cks0v68q/
type LoginCmd struct {
	UserAgent string  `json:"userAgent"`
	LanIP     string  `json:"lanIP"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type LoginRes struct {
	Uuid   string `json:"uuid"`
	ChatID string `json:"chatID"`
}

type SingleChatCmd struct {
	Base
	From string `json:"from"`
	To   string `json:"to"`
	Msg  string `json:"msg"`
}

type SingleChatRes struct {
	Cmd  int32  `json:"cmd"`
	Ack  string `json:"ack"`
	From string `json:"from"`
	To   string `json:"to"`
	Msg  string `json:"msg"`
	Uuid string `json:"uuid"`
}
