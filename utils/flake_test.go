package utils

import (
	"log"
	"testing"
)

func TestFlake(t *testing.T) {
	id, err := NextFlakeID()
	log.Println(id, err)
}
