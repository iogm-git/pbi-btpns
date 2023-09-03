package database

import (
	"fmt"

	"github.com/iogm-git/task-5-pbi-btpns-Ilham_Rahmat_Akbar/app"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var DB *gorm.DB

var server = "LAPTOP-AF7TD45M\\SQLEXPRESS"
var port = 1433
var user = "sa"
var password = "ira"
var database = "magang"

func ConDB() {
	dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;", server, user, password, port, database)

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&app.User{})
	db.AutoMigrate(&app.Photos{})
	DB = db
}
