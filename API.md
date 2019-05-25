# anychat design

* 保证消息送达: ack机制 正在发送/发送成功/发送失败
* 重发: ack timer
* 消息时间戳: 返回历史消息
* 避免重复消息: uuid
* 消息顺序: 时间戳

## 结构

```json
{
    "cmd": "xxx",
    "ack": "uuid",
    "body: "xxx_cmd data"
}
```

## cmd

```golang
const (
	WS_LOGIN = iota
	WS_LOGOUT
	WS_SINGLE_CHAT
	WS_GEO
	WS_LAN_IP
	WS_RE_CONN // 掉线重连
	WS_SERVER_PING
)
```

## ack

* client send uuid(v4) ack to server
* server reply this ack to client
