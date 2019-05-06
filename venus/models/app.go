package models

import (
	"time"

	"github.com/FlowerWrong/new_chat/venus/models/ext"
)

type App struct {
	Id        int64        `xorm:"pk autoincr BIGINT"`
	Uuid      string       `xorm:"VARCHAR"`
	Name      string       `xorm:"VARCHAR"`
	CompanyId int64        `xorm:"index INTEGER"`
	Intro     string       `xorm:"TEXT"`
	Domains   ext.StrArray `xorm:"VARCHAR" json:"domains" form:"domains[]"`
	DeletedAt time.Time    `xorm:"deleted"`
	CreatedAt time.Time    `xorm:"created"`
	UpdatedAt time.Time    `xorm:"updated"`
}

// TableName ...
func (App) TableName() string {
	return "apps"
}
