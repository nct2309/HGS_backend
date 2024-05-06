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
	UpdateFaceEncodings(houseID int, faceEncode string) error
	GetFaceEncoding(houseID int) ([]string, error)
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

func (r *deviceRepository) UpdateFaceEncodings(houseID int, faceEncode string) error {
	// the combination of both face_encode and house_id is primary key
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if err := tx.Table("Face_encoding").Create(map[string]interface{}{
		"Face_encoding": faceEncode,
		"House_id":      houseID,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *deviceRepository) GetFaceEncoding(houseID int) ([]string, error) {
	// get the face_encodings of the house (know that there is one it is NULL convert to empty string)
	var faceEncodings []string
	if err := r.db.Table("Face_encoding").Where("House_id = ? AND Face_encoding IS NOT NULL", houseID).Pluck("Face_encoding", &faceEncodings).Error; err != nil {
		return nil, err
	}
	return faceEncodings, nil
}
