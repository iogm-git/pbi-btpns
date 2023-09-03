package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/iogm-git/task-5-pbi-btpns-Ilham_Rahmat_Akbar/helpers"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := c.Cookie("token")

		if err != nil {
			if err == http.ErrNoCookie {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "login dahulu"})
				return
			}
		}

		auth := c.Request.Header.Get("Authorization")

		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Authorization Header Not Found"})
			return
		}
		splitToken := strings.Split(auth, "Bearer ")
		auth = splitToken[1]

		// token
		token, err := jwt.ParseWithClaims(auth, &helpers.JWTClaim{}, func(t *jwt.Token) (interface{}, error) {
			return helpers.JWT_KEY, nil
		})

		if err != nil {
			validation, _ := err.(*jwt.ValidationError)
			switch validation.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "login dahulu"})
				return
			case jwt.ValidationErrorExpired:
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "login dahulu"})
				return
			default:
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "login dahulu"})
				return
			}
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "login dahulu"})
			return
		}

		c.Next()
	}
}
