package entity

import "time"

type Device struct {
	ID       int    `gorm:"primaryKey;column:Device_id" json:"device_id"`
	Type     string `gorm:"Device_type" json:"device_type"`
	Name     string `gorm:"Name" json:"name"`
	Data     int    `gorm:"Device_data" json:"device_data"`
	House_id int    `gorm:"foreignKey:House_id" json:"house_id"`
}

type DataRecord struct {
	Device_id    int       `gorm:"primaryKey;foreignKey:Device_id" json:"device_id"`
	Time         time.Time `gorm:"primaryKey;column:Date_and_time" json:"time"`
	Device_data  int       `gorm:"Device_data" json:"device_data"`
	Device_state string    `gorm:"Device_state" json:"device_state"`
}
