package repository

import (
	"gorm.io/gorm"
)

type DeviceRepository interface {
	UpdateTemperature(id int, temperature float64) error
	UpdateHumidity(id int, humid float64) error
	UpdateFanSpeed(id int, speed int) error
}

type deviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepo(db *gorm.DB) DeviceRepository {
	return &deviceRepository{
		db: db,
	}
}

func (r *deviceRepository) UpdateTemperature(id int, temperature float64) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if err := tx.Table("Iot_device").Where("House_id = ? and Device_type = ?", id, "Temperature").Update("Current_data", temperature).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *deviceRepository) UpdateHumidity(id int, humid float64) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if err := tx.Table("Iot_device").Where("House_id = ? and Device_type = ?", id, "Humidity").Update("Current_data", humid).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *deviceRepository) UpdateFanSpeed(id int, speed int) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if err := tx.Table("Iot_device").Where("House_id = ? and Device_type = ?", id, "Fan").Update("Current_data", speed).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
