package repository

import (
	"context"
	"fmt"
	entity "go-jwt/internal/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUserByID(id int) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)
	UpdateUser(ctx context.Context, id string, data *entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, id string) error
	GetTempAndHumid(house_id int) (float64, float64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (userRepo *userRepository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	// err := userRepo.db.Create(user).Error
	// if err != nil {
	// 	return nil, err
	// }
	// return user, nil
	return nil, nil
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

func (userRepo *userRepository) UpdateUser(ctx context.Context, id string, data *entity.User) (*entity.User, error) {
	// err := userRepo.db.Model(&entity.User{}).Where("id = ?", id).Updates(data).Error
	// if err != nil {
	// 	return nil, err
	// }
	// return data, nil
	return nil, nil
}

func (userRepo *userRepository) DeleteUser(ctx context.Context, id string) error {
	// err := userRepo.db.Where("id = ?", id).Delete(&entity.User{}).Error
	// if err != nil {
	// 	return err
	// }
	// return nil
	return nil
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
