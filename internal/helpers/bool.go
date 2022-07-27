package helpers

import (
	"math/rand"
	"time"
)

// RandBool generates a boolean randomly
//
// https://stackoverflow.com/a/61164288/14163928
func RandBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}
