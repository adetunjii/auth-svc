package helpers

import (
	"math"
	"math/rand"
)

func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(RandInt(65, 90))
	}
	return string(bytes)
}
func RandInt(min int, max int) int {
	rand.Seed(math.MaxInt)
	return min + rand.Intn(max-min)
}
