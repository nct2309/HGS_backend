package repository

import (
	"time"

	"gorm.io/gorm"
)

type DeviceRepository interface {
	UpdateTemperature(id int, temperature float64) error
	UpdateHumidity(id int, humid float64) error
	UpdateFanSpeed(id int, speed int) error
	UpdateDevice(houseID int, deviceID int, deviceType string, data float64, state bool) error
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

func (r *deviceRepository) UpdateDevice(houseID int, deviceID int, deviceType string, data float64, state bool) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if err := tx.Table("Iot_device").Where("House_id = ? and Device_id = ?", houseID, deviceID).Update("Current_data", data).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Table("Data_record").Create(map[string]interface{}{
		"Device_id":     deviceID,
		"Date_and_time": time.Now(),
		"Device_data":   data,
		"Device_state":  state,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
