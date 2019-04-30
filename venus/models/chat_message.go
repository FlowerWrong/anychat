package models

import (
	"time"
)

type ChatMessage struct {
	Id        int64     `xorm:"pk autoincr BIGINT"`
	Uuid      string    `xorm:"VARCHAR"`
	From      int       `xorm:"INTEGER"`
	To        int       `xorm:"INTEGER"`
	Content   string    `xorm:"TEXT"`
	ReadAt    time.Time `xorm:"DATETIME"`
	DeletedAt time.Time `xorm:"DATETIME"`
	CreatedAt time.Time `xorm:"not null DATETIME"`
	UpdatedAt time.Time `xorm:"not null DATETIME"`
}

// TableName ...
func (ChatMessage) TableName() string {
	return "chat_messages"
}
