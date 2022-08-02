package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Vector2d struct {
	X float32
	Y float32
}

var (
	ScreenWidth  int = 2 * 320
	ScreenHeight int = 2 * 240
)

var RandSign = func() float32 {
	if rand.Float32() < 0.5 {
		return -1
	} else {
		return 1
	}
}
