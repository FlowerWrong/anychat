package models

import (
	"time"
)

type ChatMessage struct {
	Id        int64     `xorm:"pk autoincr BIGINT"`
	Uuid      string    `xorm:"VARCHAR"`
	From      int64     `xorm:"INTEGER"`
	To        int64     `xorm:"INTEGER"`
	Content   string    `xorm:"TEXT"`
	Ack       string    `xorm:"VARCHAR"`
	ReadAt    time.Time `xorm:"DATETIME"`
	DeletedAt time.Time `xorm:"deleted"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

// TableName ...
func (ChatMessage) TableName() string {
	return "chat_messages"
}
