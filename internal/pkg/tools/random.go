package tools

import (
	"math/rand"
	"time"
)

func GenerateRandomInt(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min+1) + min
}
