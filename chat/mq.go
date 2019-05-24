package chat

import (
	"log"
	"os"
	"strconv"

	"github.com/FlowerWrong/anychat/db"
	"github.com/nats-io/go-nats"
)

func mq(c *Client) {
	subj := "topic/" + strconv.FormatInt(c.userID, 10)
	db.MQClient().Subscribe(subj, func(msg *nats.Msg) {
		log.Printf("Received on [%s] Pid[%d]: '%s'", msg.Subject, os.Getpid(), string(msg.Data))

		// 逻辑
		err := c.logical(msg.Data)
		if err != nil {
			log.Println(err)
		}
	})
	db.MQClient().Flush()
}
