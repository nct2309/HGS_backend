package entity

import "time"

type House struct {
	ID       int    `gorm:"primaryKey;column:House_id" json:"house_id"`
	Name     string `gorm:"name" json:"name"`
	Password string `gorm:"password" json:"password"`
}

type ActivityLog struct {
	ID          int       `gorm:"primaryKey;column:Activity_id" json:"activity_id"`
	House_id    int       `gorm:"foreignKey:house_id" json:"house_id"`
	Time        time.Time `gorm:"time" json:"time"`
	Device      string    `gorm:"device" json:"device"`
	TypeOfEvent string    `gorm:"type_of_event" json:"type_of_event"`
}

type FaceEncoding struct {
	Face_encoding string `gorm:"primaryKey;column:face_encoding" json:"face_encoding"`
	House_id      int    `gorm:"foreignKey:house_id" json:"house_id"`
}

//Set and house setting???

type HouseSetting struct {
  // combination of Name and House_id is primary key
  Name string `gorm:"primaryKey;column:Name" json:"name"`
  House_id int `gorm:"primaryKey;foreignKey:House_id" json:"house_id"`
  Selected bool `gorm:"selected" json:"selected"`
}

type Set struct {
  Device_id int `gorm:"primaryKey;foreignKey:Device_id" json:"device_id"`
  House_id int `gorm:"primaryKey;foreignKey:House_id" json:"house_id"`
  Name string `gorm:"primaryKey;foreignKey:Name" json:"name"`
  Device_data int `gorm:"Device_data" json:"device_data"`
  Device_state string `gorm:"Device_state" json:"device_state"`
}
