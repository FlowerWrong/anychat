package models

import (
	"time"
)

type Room struct {
	Id        int64     `xorm:"pk autoincr BIGINT" json:"-"`
	Uuid      string    `xorm:"VARCHAR"`
	Name      string    `xorm:"VARCHAR"`
	Intro     string    `xorm:"TEXT"`
	CreatorId int64     `xorm:"INTEGER"`
	Logo      string    `xorm:"VARCHAR"`
	DeletedAt time.Time `xorm:"deleted" json:"-"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated" json:"-"`
}

// TableName ...
func (Room) TableName() string {
	return "rooms"
}
