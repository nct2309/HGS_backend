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
	GetDashboardData(house_id int) (float64, float64, float64, float64, error)
	UpdateSet(deviceID int, data float64, state bool, house_id int, setting string) error
	UpdateManySets([]entity.Set) error
	GetAllNotifications(userID int) ([]entity.Notification, error)
	GetUnreadNotifications(userID int) ([]entity.Notification, error)
	CreateNotification(userID int, houseID int, notification *entity.Notification) error
	CreateActivityLog(activityLog *entity.ActivityLog) error
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

func (userRepo *userRepository) GetDashboardData(house_id int) (float64, float64, float64, float64, error) {
	// Using Device_id to ordering the output but finally take only Current_data, if select only Current_data, it order by Current_data
	var temp float64
	var humid float64
	var light float64
	var fan_speed float64
	// use LIMIT 1 to get the first row
	err := userRepo.db.Table("Iot_device").Where("House_id = ? and Device_type = ?", house_id, "Temperature").Select("Current_data").Order("Device_id").Limit(1).Scan(&temp).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}
	err = userRepo.db.Table("Iot_device").Where("House_id = ? and Device_type = ?", house_id, "Humidity").Select("Current_data").Order("Device_id").Limit(1).Scan(&humid).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}
	err = userRepo.db.Table("Iot_device").Where("House_id = ? and Device_type = ?", house_id, "Light").Select("Current_data").Order("Device_id").Limit(1).Scan(&light).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}
	err = userRepo.db.Table("Iot_device").Where("House_id = ? and Device_type = ?", house_id, "Fan").Select("Current_data").Order("Device_id").Limit(1).Scan(&fan_speed).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return temp, humid, light, fan_speed, nil
}

func (userRepo *userRepository) GetHouseID(userID int) ([]int, error) {
	var houseIDs []int
	err := userRepo.db.Table("Own").Where("User_id = ?", userID).Select("House_id").Scan(&houseIDs).Error
	if err != nil {
		return nil, err
	}
	return houseIDs, nil
}

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
	err := userRepo.db.Table("Set").Where("\"Set\".House_id = ? and \"Set\".Name = ?", house_id, settingName).Joins("JOIN Iot_device ON \"Set\".Device_id = Iot_device.Device_id").Select("\"Set\".*, Iot_device.Name as Device_name").Find(&sets).Error
	if err != nil {
		return nil, err
	}
	return sets, nil
}

func (userRepo *userRepository) GetActivityLogByHouseID(house_id int) ([]entity.ActivityLog, error) {
	var activityLogs []entity.ActivityLog
	err := userRepo.db.Table("Activity_log").Where("House_id = ?", house_id).Find(&activityLogs).Error
	if err != nil {
		return nil, err
	}
	return activityLogs, nil
}

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

// a func combine two above
func (userRepo *userRepository) UpdateSet(deviceID int, data float64, state bool, house_id int, setting string) error {
	// transaction
	tx := userRepo.db.Begin()
	err := tx.Table("Set").Where("House_id = ? and Name = ? and Device_id = ?", house_id, setting, deviceID).Update("Device_data", data).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Table("Set").Where("House_id = ? and Name = ? and Device_id = ?", house_id, setting, deviceID).Update("Device_state", state).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	if err != nil {
		return err
	}
	return nil
}

func (userRepo *userRepository) UpdateManySets(sets []entity.Set) error {
	// transaction
	tx := userRepo.db.Begin()
	for _, set := range sets {

		// Update data and state in one sql statement
		err := tx.Table("Set").Where("House_id = ? and Name = ? and Device_id = ?", set.House_id, set.Name, set.Device_id).Updates(map[string]interface{}{"Device_data": set.Device_data, "Device_state": set.Device_state}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err := tx.Commit().Error
	if err != nil {
		return err
	}
	return nil
}

func (userRepo *userRepository) GetAllNotifications(userID int) ([]entity.Notification, error) {
	var notifications []entity.Notification
	err := userRepo.db.Table("Send").Where("User_id = ?", userID).Joins("JOIN Notification ON \"Send\".Notification_id = Notification.Notification_id").Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (userRepo *userRepository) GetUnreadNotifications(userID int) ([]entity.Notification, error) {
	var notifications []entity.Notification
	err := userRepo.db.Table("Send").Where("User_id = ? and Read = ?", userID, false).Joins("JOIN Notification ON \"Send\".Notification_id = Notification.Notification_id").Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (userRepo *userRepository) CreateNotification(userID int, houseID int, notification *entity.Notification) error {
	// transaction
	tx := userRepo.db.Begin()
	// Create notification in Notification table
	err := tx.Table("Notification").Create(notification).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// get the id of the notification just created to insert into Send
	err = tx.Table("Send").Create(&entity.Send{Notification_id: notification.ID, User_id: userID, House_id: houseID}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	if err != nil {
		return err
	}
	return nil
}

func (userRepo *userRepository) CreateActivityLog(activityLog *entity.ActivityLog) error {
	err := userRepo.db.Table("Activity_log").Create(activityLog).Error
	if err != nil {
		return err
	}
	return nil
}
