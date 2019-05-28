# anychat design

* 保证消息送达: ack机制 正在发送/发送成功/发送失败
* 重发: ack timer, 3 times
* 消息时间戳: 返回历史消息
* 避免重复消息: uuid https://chai2010.cn/advanced-go-programming-book/ch6-cloud/ch6-01-dist-id.html
* 消息顺序: 时间戳

## 结构

```json
{
    "cmd": "xxx",
    "ack": "uuid",
    "data: "xxx_cmd data"
}
```

## cmds

```golang
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
```

## ack

* client send uuid(v4) ack to server
* server reply this ack to client
