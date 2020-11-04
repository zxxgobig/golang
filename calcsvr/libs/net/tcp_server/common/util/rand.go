package util

import (
	"math/rand"
	"time"
)

var r * rand.Rand


func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func GetRandUint64() uint64{
	return r.Uint64()
}

func GetRandInt() int {
	return r.Int()
}
