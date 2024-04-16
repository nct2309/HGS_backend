package repository

import (
	"fmt"
	entity "go-jwt/internal/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(id int) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)
	GetTempAndHumid(house_id int) (float64, float64, error)
	GetHouseID(userID int) ([]int, error)
	GetHouseSettingByHouseID(house_id int) ([]entity.HouseSetting, error)
	GetSetOfHouseSetting(house_id int, settingName string) ([]entity.Set, error)
	GetActivityLogByHouseID(house_id int) ([]entity.ActivityLog, error)
	UpdateDeviceData(deviceID int, data float64, house_id int, setting string) error
	UpdataDeviceState(deviceID int, state bool, house_id int, setting string) error
}
type userRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (userRepo *userRepository) GetUserByID(id int) (*entity.User, error) {
	user := entity.User{}

	err := userRepo.db.Table("Users").Where("User_id = ?", id).First(&user).Error

	if err != nil {
		fmt.Print("Error", err)
		return nil, err
	}
	return &user, nil
}

func (userRepo *userRepository) GetUserByUsername(username string) (*entity.User, error) {
	user := entity.User{}

	err := userRepo.db.Table("Users").Where("Username = ?", username).First(&user).Error

	if err != nil {
		fmt.Print("Error", err)
		return nil, err
	}

	return &user, nil
}

func (userRepo *userRepository) GetTempAndHumid(house_id int) (float64, float64, error) {
	var temp float64
	var humid float64
	err := userRepo.db.Table("Iot_device").Where("House_id = ? and Device_type = ?", house_id, "Temperature").Select("Current_data").Scan(&temp).Error
	if err != nil {
		return 0, 0, err
	}
	err = userRepo.db.Table("Iot_device").Where("House_id = ? and Device_type = ?", house_id, "Humidity").Select("Current_data").Scan(&humid).Error
	if err != nil {
		return 0, 0, err
	}
	return temp, humid, nil
}

//	type User struct {
//		ID       int    `gorm:"primaryKey;column:User_id" json:"user_id"`
//		Username string `gorm:"username" json:"username"`
//		Password string `gorm:"password" json:"password"`
//	}
//
// type Own struct { // More descriptive name
//
//		UserID  int `gorm:"primary_key;foreignKey:User_id"`
//		HouseID int `gorm:"primary_key;foreignKey:House_id"`
//	}
//
//	type House struct {
//		ID       int    `gorm:"primaryKey;column:House_id" json:"house_id"`
//		Name     string `gorm:"name" json:"name"`
//		Password string `gorm:"password" json:"password"`
//	}

func (userRepo *userRepository) GetHouseID(userID int) ([]int, error) {
	var houseIDs []int
	err := userRepo.db.Table("Own").Where("User_id = ?", userID).Select("House_id").Scan(&houseIDs).Error
	if err != nil {
		return nil, err
	}
	return houseIDs, nil
}

// type HouseSetting struct {
// 	// combination of Name and House_id is primary key
// 	Name string `gorm:"primaryKey;column:Name" json:"name"`
// 	House_id int `gorm:"primaryKey;foreignKey:House_id" json:"house_id"`
// 	Selected bool `gorm:"selected" json:"selected"`
//   }

// type Set struct {
// 	Device_id    int    `gorm:"primaryKey;foreignKey:Device_id" json:"device_id"`
// 	Device_name  string `gorm:"Device_name" json:"device_name"` // this not in the Set table of the database but is needed for frontend
// 	House_id     int    `gorm:"primaryKey;foreignKey:House_id" json:"house_id"`
// 	Name         string `gorm:"primaryKey;foreignKey:Name" json:"name"`
// 	Device_data  int    `gorm:"Device_data" json:"device_data"`
// 	Device_state bool `gorm:"Device_state" json:"device_state"`
// }

func (userRepo *userRepository) GetHouseSettingByHouseID(house_id int) ([]entity.HouseSetting, error) {
	var houseSettings []entity.HouseSetting
	err := userRepo.db.Table("House_setting").Where("House_id = ?", house_id).Find(&houseSettings).Error
	if err != nil {
		return nil, err
	}
	return houseSettings, nil
}

func (userRepo *userRepository) GetSetOfHouseSetting(house_id int, settingName string) ([]entity.Set, error) {
	var sets []entity.Set
	// get all the set of a house setting which join table set to device
	// err := userRepo.db.Table("Set").Where("House_id = ? and Name = ?", house_id, settingName).Find(&sets).Error
	// err := userRepo.db.Table("Set").Where("House_id = ? and Name = ?", house_id, settingName).Joins("JOIN Iot_device ON Set.Device_id = Iot_device.Device_id").Find(&sets).Error
	// Find all the set of a house setting which join table set to device and get the Name in Iot_device as device_name
	err := userRepo.db.Table("Set").Where("\"Set\".House_id = ? and \"Set\".Name = ?", house_id, settingName).Joins("JOIN Iot_device ON \"Set\".Device_id = Iot_device.Device_id").Select("\"Set\".*, Iot_device.Name as Device_name").Find(&sets).Error
	if err != nil {
		return nil, err
	}
	return sets, nil
}

// type Device struct {
// 	ID       int    `gorm:"primaryKey;column:Device_id" json:"device_id"`
// 	Type     string `gorm:"Device_type" json:"device_type"`
// 	Name     string `gorm:"Name" json:"name"`
// 	Data     int    `gorm:"Device_data" json:"device_data"`
// 	House_id int    `gorm:"foreignKey:House_id" json:"house_id"`
// }

// type ActivityLog struct {
// 	ID          int       `gorm:"primaryKey;column:Activity_id" json:"activity_id"`
// 	House_id    int       `gorm:"foreignKey:House_id" json:"house_id"`
// 	Time        time.Time `gorm:"Time" json:"time"`
// 	Device      string    `gorm:"Device" json:"device"`
// 	TypeOfEvent string    `gorm:"Type_of_event" json:"type_of_event"`
// }

func (userRepo *userRepository) GetActivityLogByHouseID(house_id int) ([]entity.ActivityLog, error) {
	var activityLogs []entity.ActivityLog
	err := userRepo.db.Table("Activity_log").Where("House_id = ?", house_id).Find(&activityLogs).Error
	if err != nil {
		return nil, err
	}
	return activityLogs, nil
}

//	type DataRecord struct {
//		Device_id    int       `gorm:"primaryKey;foreignKey:Device_id" json:"device_id"`
//		Time         time.Time `gorm:"primaryKey;column:Date_and_time" json:"time"`
//		Device_data  float64       `gorm:"Device_data" json:"device_data"`
//		Device_state bool    `gorm:"Device_state" json:"device_state"`
//	}
func (userRepo *userRepository) UpdateDeviceData(deviceID int, data float64, house_id int, setting string) error {
	err := userRepo.db.Table("Set").Where("House_id = ? and Name = ? and Device_id = ?", house_id, setting, deviceID).Update("Device_data", data).Error
	if err != nil {
		return err
	}
	// update in Iot_device table
	err = userRepo.db.Table("Iot_device").Where("Device_id = ?", deviceID).Update("Current_data", data).Error
	if err != nil {
		return err
	}

	return nil
}

func (userRepo *userRepository) UpdataDeviceState(deviceID int, state bool, house_id int, setting string) error {
	err := userRepo.db.Table("Set").Where("House_id = ? and Name = ? and Device_id = ?", house_id, setting, deviceID).Update("Device_state", state).Error
	if err != nil {
		return err
	}

	return nil
}
