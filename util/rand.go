package util

import (
	"math/rand"
	"time"
)

func ShuffleUint32Arr(arr []uint32) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	r.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
}
