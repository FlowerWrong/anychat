package models

import (
	"time"
)

type RoomMessage struct {
	Id        int64     `xorm:"pk autoincr BIGINT" json:"-"`
	Uuid      string    `xorm:"VARCHAR"`
	From      int64     `xorm:"INTEGER"`
	RoomId    int64     `xorm:"INTEGER"`
	Content   string    `xorm:"TEXT"`
	DeletedAt time.Time `xorm:"deleted" json:"-"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated" json:"-"`
}

// TableName ...
func (RoomMessage) TableName() string {
	return "room_messages"
}
