package controllers

import (
	"net/http"
	"net/mail"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/iogm-git/task-5-pbi-btpns-Ilham_Rahmat_Akbar/app"
	"github.com/iogm-git/task-5-pbi-btpns-Ilham_Rahmat_Akbar/database"
	"github.com/iogm-git/task-5-pbi-btpns-Ilham_Rahmat_Akbar/helpers"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	// struct user
	var user app.User

	// json
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	// set time
	user.Created_at = time.Now()
	user.Updated_at = time.Now()

	// validasi
	result, err := govalidator.ValidateStruct(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	// hash
	hashPwd, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashPwd)

	// insert tabel users
	if err := database.DB.Create(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"berhasil": result})
}

func Login(c *gin.Context) {
	// json
	var input app.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	// cek
	if govalidator.IsNull(input.Email) || govalidator.IsNull(input.Password) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "kolom harus diisi"})
		return
	}

	if _, err := mail.ParseAddress(input.Email); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "harus berupa email"})
		return
	}

	// cek apakah sudah daftar
	var user app.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "regitrasi dahulu"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}
	}

	// jika password salah
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "password salah"})
		return
	}

	// jwt
	expTime := time.Now().Add(time.Minute * 1)
	claims := &helpers.JWTClaim{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "jwt",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// declaration
	initial := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign
	token, err := initial.SignedString(helpers.JWT_KEY)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	// set token into cookie
	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"berhasil": true})
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"berhasil": true})
}
