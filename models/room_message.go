package models

import (
	"time"
)

type RoomMessage struct {
	Id        int64     `xorm:"pk autoincr BIGINT"`
	Uuid      string    `xorm:"VARCHAR"`
	From      int64     `xorm:"INTEGER"`
	RoomId    int64     `xorm:"INTEGER"`
	Content   string    `xorm:"TEXT"`
	DeletedAt time.Time `xorm:"deleted"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

// TableName ...
func (RoomMessage) TableName() string {
	return "room_messages"
}
