package models

import (
	"time"
)

type User struct {
	Id             int64     `xorm:"pk autoincr BIGINT"`
	Uuid           string    `xorm:"VARCHAR"`
	Username       string    `xorm:"VARCHAR"`
	PasswordDigest string    `xorm:"VARCHAR"`
	Mobile         string    `xorm:"VARCHAR"`
	Email          string    `xorm:"VARCHAR"`
	Avatar         string    `xorm:"VARCHAR"`
	Note           string    `xorm:"VARCHAR"`
	FirstLoginAt   time.Time `xorm:"DATETIME"`
	LastActiveAt   time.Time `xorm:"DATETIME"`
	DeletedAt      time.Time `xorm:"deleted"`
	CreatedAt      time.Time `xorm:"created"`
	UpdatedAt      time.Time `xorm:"updated"`
}

// TableName ...
func (User) TableName() string {
	return "users"
}
