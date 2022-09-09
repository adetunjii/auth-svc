package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	LowercaseAlphabets = "abcdefghijklmnopqrstuvwxyz"
	UppercaseAlphabets = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits             = "0123456789"
	Symbols            = "!@#$%^&*()_{}[]><.-"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomUUID() string {
	return uuid.New().String()
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder

	k := len(LowercaseAlphabets)

	for i := 0; i < n; i++ {
		char := LowercaseAlphabets[rand.Intn(k)]
		sb.WriteByte(char)
	}

	return sb.String()
}

func RandomName() string {
	return RandomString(6)
}

func RandomPhoneNumber() string {
	num := RandomInt(1000000000, 9999999999)
	return fmt.Sprintf("%d", num)
}

func RandomCountryCode() string {
	codes := []string{"234", "233", "1"}

	return codes[rand.Intn(len(codes))]
}

func RandomPassword() string {
	var sb strings.Builder

	k := len(LowercaseAlphabets)

	for i := 0; i < 6; i++ {
		char := LowercaseAlphabets[rand.Intn(k)]
		sb.WriteByte(char)
	}

	sb.WriteByte(UppercaseAlphabets[rand.Intn(2)])
	sb.WriteByte(Digits[rand.Intn(2)])
	sb.WriteByte(Symbols[rand.Intn(2)])

	return sb.String()
}

func RandomOtp() int {
	return int(RandomInt(111111, 999999))
}
