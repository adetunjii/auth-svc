package jwt

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	jwt.StandardClaims
}

func validateToken(token string, key *rsa.PublicKey) (*jwt.Token, *Claims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("valid proto token required")
		}
		return key, nil
	})
	if claims, ok := jwtToken.Claims.(*Claims); ok && jwtToken.Valid {
		return jwtToken, claims, nil
	}
	return nil, nil, err
}
