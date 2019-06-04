package models

import (
	"time"
)

type LoginLog struct {
	Id        int64     `xorm:"pk autoincr BIGINT" json:"-"`
	UserId    int64     `xorm:"INTEGER"`
	Ua        string    `xorm:"VARCHAR"`
	Ip        string    `xorm:"VARCHAR"`
	LanIp     string    `xorm:"VARCHAR"`
	Os        string    `xorm:"VARCHAR"`
	Browser   string    `xorm:"VARCHAR"`
	Latitude  float64   `xorm:"NUMERIC"`
	Longitude float64   `xorm:"NUMERIC"`
	DeletedAt time.Time `xorm:"deleted" json:"-"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated" json:"-"`
}

// TableName ...
func (LoginLog) TableName() string {
	return "login_logs"
}
