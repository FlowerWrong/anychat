package models

import (
	"time"
)

type Room struct {
	Id        int64     `xorm:"pk autoincr BIGINT"`
	Uuid      string    `xorm:"VARCHAR"`
	Name      string    `xorm:"VARCHAR"`
	Intro     string    `xorm:"TEXT"`
	Logo      string    `xorm:"VARCHAR"`
	DeletedAt time.Time `xorm:"deleted"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

// TableName ...
func (Room) TableName() string {
	return "rooms"
}
