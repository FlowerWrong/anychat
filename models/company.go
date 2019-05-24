package models

import (
	"time"
)

// Company ...
type Company struct {
	Id                      int64     `xorm:"pk autoincr BIGINT"`
	Uuid                    string    `xorm:"VARCHAR"`
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
	DeletedAt               time.Time `xorm:"deleted"`
	CreatedAt               time.Time `xorm:"created"`
	UpdatedAt               time.Time `xorm:"updated"`
}

// TableName ...
func (Company) TableName() string {
	return "companies"
}
