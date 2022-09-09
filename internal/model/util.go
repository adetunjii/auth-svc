// This file contains all the utility/helper functions for the model package

package model

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmailAddress     = errors.New("invalid email address")
	ErrInvalidUserId           = errors.New("invalid user id")
	ErrInvalidPasswordTooShort = errors.New("password is too short. Minimun length is 8")
	ErrInvalidPasswordTooLong  = errors.New("password is too long. Maximum length is 72")
	ErrInvalidPasswordPattern  = errors.New("invalid password. Password must contain at least one uppercase letter, one lowercase letter, one digit and one symbol")
)

const (
	DefaultCost = 10
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	return string(hash), nil
}

func ComparePassword(hash string, password string) error {
	if password == "" || hash == "" {
		return errors.New("empty password or hash")
	}

	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// SanitizeUnicode will remove undesirable Unicode characters from a string.
func SanitizeUnicode(s string) string {
	return strings.Map(filterBlocklist, s)
}

// filterBlocklist returns `r` if it is not in the blocklist, otherwise drop (-1).
// Blocklist is taken from https://www.w3.org/TR/unicode-xml/#Charlist
func filterBlocklist(r rune) rune {
	const drop = -1
	switch r {
	case '\u0340', '\u0341': // clones of grave and acute; deprecated in Unicode
		return drop
	case '\u17A3', '\u17D3': // obsolete characters for Khmer; deprecated in Unicode
		return drop
	case '\u2028', '\u2029': // line and paragraph separator
		return drop
	case '\u202A', '\u202B', '\u202C', '\u202D', '\u202E': // BIDI embedding controls
		return drop
	case '\u206A', '\u206B': // activate/inhibit symmetric swapping; deprecated in Unicode
		return drop
	case '\u206C', '\u206D': // activate/inhibit Arabic form shaping; deprecated in Unicode
		return drop
	case '\u206E', '\u206F': // activate/inhibit national digit shapes; deprecated in Unicode
		return drop
	case '\uFFF9', '\uFFFA', '\uFFFB': // interlinear annotation characters
		return drop
	case '\uFEFF': // byte order mark
		return drop
	case '\uFFFC': // object replacement character
		return drop
	}

	// Scoping for musical notation
	if r >= 0x0001D173 && r <= 0x0001D17A {
		return drop
	}

	// Language tag code points
	if r >= 0x000E0000 && r <= 0x000E007F {
		return drop
	}

	return r
}
