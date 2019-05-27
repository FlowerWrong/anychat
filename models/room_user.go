package models

import (
	"time"
)

type RoomUser struct {
	Id        int64     `xorm:"pk autoincr BIGINT"`
	Uuid      string    `xorm:"VARCHAR"`
	UserId    int64     `xorm:"INTEGER"`
	RoomId    int64     `xorm:"INTEGER"`
	Nickname  string    `xorm:"VARCHAR"`
	DeletedAt time.Time `xorm:"deleted"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

// TableName ...
func (RoomUser) TableName() string {
	return "room_users"
}
