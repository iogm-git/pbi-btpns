package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/iogm-git/task-5-pbi-btpns-Ilham_Rahmat_Akbar/app"
	"github.com/iogm-git/task-5-pbi-btpns-Ilham_Rahmat_Akbar/database"
	"github.com/iogm-git/task-5-pbi-btpns-Ilham_Rahmat_Akbar/helpers"
	"gorm.io/gorm"
)

func Update(c *gin.Context) {
	// identify
	var identify app.User
	auth := helpers.Auth(strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1])
	database.DB.Where("email = ?", auth).First(&identify)

	// cek apakah parameter sama
	userId, _ := strconv.Atoi(c.Param("userId"))
	if userId != identify.Id {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "masukan parameter sesuai dengan id"})
		return
	}

	// cek apakah sudah daftar atau login
	var user app.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"msg": "login / register dahulu"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}
	}

	created_at := user.Created_at
	user.Created_at = created_at
	user.Updated_at = time.Now()

	result, err := govalidator.ValidateStruct(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	if database.DB.Model(&user).Where("id = ?", userId).Updates(&user).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "tidak dapat mengupdate"})
		return
	}

	if result {
		c.JSON(http.StatusOK, gin.H{"msg": user})
	}

}

func Delete(c *gin.Context) {
	// identify
	var identify app.User
	auth := helpers.Auth(strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1])
	database.DB.Where("email = ?", auth).First(&identify)

	// cek apakah parameter sama
	userId, _ := strconv.Atoi(c.Param("userId"))
	if userId != identify.Id {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "masukan parameter sesuai dengan id"})
		return
	}

	var user app.User

	if database.DB.Delete(&user, userId).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "tidak dapat menghapus"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Data berhasil dihapus"})
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
}
