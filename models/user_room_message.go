package models

import (
	"time"
)

type UserRoomMessage struct {
	Id            int64     `xorm:"pk autoincr BIGINT" json:"-"`
	Uuid          string    `xorm:"VARCHAR"`
	UserId        int64     `xorm:"INTEGER"`
	RoomMessageId int64     `xorm:"INTEGER"`
	ReadAt        time.Time `xorm:"DATETIME"`
	DeletedAt     time.Time `xorm:"deleted" json:"-"`
	CreatedAt     time.Time `xorm:"created"`
	UpdatedAt     time.Time `xorm:"updated" json:"-"`
}

// TableName ...
func (UserRoomMessage) TableName() string {
	return "user_room_messages"
}
