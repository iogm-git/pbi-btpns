package helpers

import "github.com/golang-jwt/jwt"

func Auth(code string) string {
	token, _ := jwt.ParseWithClaims(code, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return JWT_KEY, nil
	})

	email := ""

	if claims, ok := token.Claims.(*JWTClaim); ok && token.Valid {
		email = claims.Email
	}

	return email
}
