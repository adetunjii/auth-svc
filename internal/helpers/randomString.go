package helpers

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(RandInt(65, 90))
	}
	return string(bytes)
}
func RandInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

func RandomOtp() int {
	return RandInt(100000, 999999)
}

func TrimPhoneNumber(phoneNumber string, phoneCode string) string {

	if phoneCode == "234" && phoneNumber[0] == '0' {
		return phoneNumber[1:]
	}

	return phoneNumber
}
