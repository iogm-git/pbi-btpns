package controllers

import (
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iogm-git/task-5-pbi-btpns-Ilham_Rahmat_Akbar/app"
	"github.com/iogm-git/task-5-pbi-btpns-Ilham_Rahmat_Akbar/database"
	"github.com/iogm-git/task-5-pbi-btpns-Ilham_Rahmat_Akbar/helpers"
)

func Upload(c *gin.Context) {
	// identify
	var identify app.User
	auth := helpers.Auth(strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1])
	database.DB.Where("email = ?", auth).First(&identify)

	// cek apakah parameter sama
	photoId := 0
	if c.Param("photoId") == "" {
		photoId = identify.Id
	} else {
		paramId, _ := strconv.Atoi(c.Param("photoId"))
		photoId = paramId
	}
	if photoId != identify.Id {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "masukan parameter sesuai dengan id"})
		return
	}

	// tangkap form
	form, _ := c.MultipartForm()

	// handle file
	urlfile := form.File["photo_url"][0]
	filename := helpers.GenerateRandomString(7) + path.Ext(urlfile.Filename) // rename file
	form.File["photo_url"][0].Filename = filename

	// wajib ada title dan caption
	title, no_title := form.Value["title"]
	caption, no_caption := form.Value["caption"]
	if !no_title || !no_caption {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "sertakan parameter"})
		return
	}
	if title[0] == "" || caption[0] == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "kolom harus diisi"})
		return
	}

	// get id
	var user app.User
	database.DB.Where("email = ?", auth).First(&user)

	// set upload
	var photo app.Photos
	photo.Title = title[0]
	photo.Caption = caption[0]
	photo.PhotoUrl = filename

	// set path folder
	path := "img/" + filename

	if database.DB.Where("id = ?", user.Id).First(&photo).RowsAffected == 0 {
		photo.Id = user.Id
		photo.UserId = user.Id
		// create photos
		if err := database.DB.Create(&photo).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}
		// simpan gambar
		if err := c.SaveUploadedFile(urlfile, path); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "gagal simpan gambar"})
			return
		}
	} else {
		// simpan gambar
		if err := c.SaveUploadedFile(urlfile, path); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "gagal simpan gambar"})
			return
		}

		var oldPhoto app.Photos
		database.DB.First(&oldPhoto, "user_id = ?", user.Id)

		// Delete the file from the server
		oldImg := "img/" + oldPhoto.PhotoUrl
		if _, err := os.Stat(oldImg); err == nil {
			err := os.Remove(oldImg)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from upload folder"})
				return
			}
		}

		photo.PhotoUrl = filename

		// update photos
		if database.DB.Model(&photo).Where("id = ?", user.Id).Updates(&photo).RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "tidak dapat mengupdate"})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"berhasil": photo})
}

func Show(c *gin.Context) {
	// identify
	var identify app.User
	auth := helpers.Auth(strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1])
	database.DB.Where("email = ?", auth).First(&identify)

	// cek photo apakah ada
	var photos app.Photos
	if database.DB.Where("id = ?", identify.Id).First(&photos).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "anda belum menambahkan foto"})
		return
	}

	// kalo ada kembalikan data
	var user app.User
	database.DB.Where("email = ?", auth).First(&user)
	database.DB.Preload("Photos").Find(&user)

	c.AbortWithStatusJSON(http.StatusOK, gin.H{"data": user.Photos})
}

func Remove(c *gin.Context) {
	// identify
	var identify app.User
	auth := helpers.Auth(strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1])
	database.DB.Where("email = ?", auth).First(&identify)

	// cek apakah parameter sama
	photoId, _ := strconv.Atoi(c.Param("photoId"))
	if photoId != identify.Id {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"harusnya masukan parameter id": identify.Id})
		return
	}

	var photo app.Photos
	database.DB.First(&photo, identify.Id)

	// Delete the file from the server
	rmv := "img/" + photo.PhotoUrl
	if _, err := os.Stat(rmv); err == nil {
		err := os.Remove(rmv)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from upload folder"})
			return
		}
	}

	var deletePhotos app.Photos

	// update database
	if database.DB.Delete(&deletePhotos, identify.Id).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "tidak dapat mengupdate"})
		return
	}

}
