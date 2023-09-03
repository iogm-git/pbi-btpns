package app

import "time"

type User struct {
	Id         int       `valid:"required,type(int)" gorm:"primaryKey" json:"id"`
	Username   string    `valid:"required,type(string)" gorm:"type:varchar(100)" json:"username"`
	Email      string    `valid:"required,email,type(string)" gorm:"type:varchar(100);unique;not null;" json:"email"`
	Password   string    `valid:"required,type(string),minstringlength(6)" gorm:"type:varchar(100)" json:"password"`
	Created_at time.Time `valid:"required" gorm:"type:time;default:CURRENT_TIMESTAMP" json:"created_at"`
	Updated_at time.Time `valid:"required" gorm:"type:time;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Photos     *Photos   `gorm:"foreignKey:UserId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"photos"`
}

type Photos struct {
	Id       int    `form:"id" valid:"required,type(int)" gorm:"primaryKey" json:"id"`
	Title    string `form:"title" gorm:"type:varchar(100)" json:"title"`
	Caption  string `form:"caption" gorm:"type:varchar(253)" json:"caption"`
	PhotoUrl string `form:"photo_url" gorm:"type:varchar(177)" json:"photo_url"`
	UserId   int    `form:"user_id" valid:"required,type(int);" gorm:"type:int;primaryKey" json:"user_id"`
}
