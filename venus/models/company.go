package models

import (
	"time"
)

// Company ...
type Company struct {
	Id                      int64     `xorm:"pk autoincr BIGINT"`
	Name                    string    `xorm:"VARCHAR"`
	AliasName               string    `xorm:"VARCHAR"`
	Intro                   string    `xorm:"TEXT"`
	LegalPerson             string    `xorm:"VARCHAR"`
	Tel                     string    `xorm:"VARCHAR"`
	Website                 string    `xorm:"VARCHAR"`
	Email                   string    `xorm:"VARCHAR"`
	Address                 string    `xorm:"VARCHAR"`
	EstablishedAt           time.Time `xorm:"DATETIME"`
	UnifiedSocialCreditCode string    `xorm:"VARCHAR"`
	Logo                    string    `xorm:"VARCHAR"`
	DeletedAt               time.Time `xorm:"DATETIME"`
	CreatedAt               time.Time `xorm:"not null DATETIME"`
	UpdatedAt               time.Time `xorm:"not null DATETIME"`
}

// TableName ...
func (Company) TableName() string {
	return "companies"
}
