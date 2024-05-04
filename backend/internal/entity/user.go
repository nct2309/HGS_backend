package entity

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserPasswordNotMatch = errors.New("password not match")
)

type User struct {
	ID       int    `gorm:"primaryKey;column:User_id" json:"user_id"`
	Username string `gorm:"username" json:"username"`
	Password string `gorm:"password" json:"password"`
}

// Combination of User and House id is primary key for Own table and the foreign key for House table
type Own struct { // More descriptive name
	UserID  int `gorm:"primary_key;foreignKey:User_id"`
	HouseID int `gorm:"primary_key;foreignKey:House_id"`
}

type Notification struct {
	ID          int       `gorm:"primaryKey;column:Notification_id" json:"notification_id"`
	Description string    `gorm:"Description" json:"description"`
	Time        time.Time `gorm:"Time" json:"time"` //3/25/2024 5:06:00 PM
	Title       string    `gorm:"Title" json:"title"`
	Read        bool      `gorm:"Read" json:"read"`
}

type Send struct {
	Notification_id int `gorm:"primaryKey;foreignKey:Notification_id" json:"notification_id"`
	User_id         int `gorm:"primaryKey;foreignKey:User_id" json:"user_id"`
	House_id        int `gorm:"foreignkey:House_id" json:"house_id"`
}
