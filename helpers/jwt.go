package helpers

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("abcd1234")

type JWTClaim struct {
	Email string
	jwt.RegisteredClaims
}
