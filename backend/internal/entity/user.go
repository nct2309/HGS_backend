package entity

import (
	"errors"
	"time"
)

var (
	ERR_USER_NOT_FOUND          = errors.New("User not found")
	ERR_USER_PASSWORD_NOT_MATCH = errors.New("Password not match")
)

type User struct {
	ID       int    `gorm:"primaryKey;column:User_id" json:"user_id"`
	Username string `gorm:"username" json:"username"`
	Password string `gorm:"password" json:"password"`
}

// Combination of User and House id is primary key for Own table and the foreign key for House table
type Own struct { // More descriptive name
	UserID  int `gorm:"primary_key;foreignKey:UserID"`
	HouseID int `gorm:"primary_key;foreignKey:HouseID"`
}

type Notification struct {
	ID          int       `gorm:"primaryKey;column:Notification_id" json:"notification_id"`
	Description string    `gorm:"description" json:"description"`
	Time        time.Time `gorm:"time" json:"time"` //3/25/2024 5:06:00 PM
	Title       string    `gorm:"title" json:"title"`
}

type Send struct {
	Notification_id int `gorm:"primaryKey;foreignKey:Notification_id" json:"notification_id"`
	User_id         int `gorm:"primaryKey;foreignKey:User_id" json:"user_id"`
	House_id        int `gorm:"foreignkey:House_id" json:"house_id"`
}
