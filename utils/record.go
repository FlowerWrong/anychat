package utils

import (
	"errors"

	"github.com/FlowerWrong/anychat/db"
)

// UpdateRecord ...
func UpdateRecord(id int64, modelPointer interface{}) error {
	affected, err := db.Engine().Id(id).Update(modelPointer)
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.New("affected not 1")
	}
	return nil
}
