package usecase

import (
	entity "go-jwt/internal/entity"
	repository "go-jwt/internal/infrastructure/repository"
	"go-jwt/internal/middleware/token"
	"strconv"
	"strings"
)

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

type UserUsecase interface {
	// CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUser(id int) (*entity.User, error)
	// UpdateUser(ctx context.Context, id string, data *entity.User) (*entity.User, error)
	// DeleteUser(ctx context.Context, id string) error
	AuthenticateUser(username string, password string) (*entity.User, string, []int, error)
	GetTempAndHumid(house_id int) (float64, float64, error)
	GetDashboardData(house_id int) (float64, float64, float64, float64, error)
	GetHouseSettingByHouseID(house_id int) ([]entity.HouseSetting, error)
	GetSetOfHouseSetting(house_id int, settingName string) ([]entity.Set, error)
	GetActivityLogByHouseID(house_id int) ([]entity.ActivityLog, error)
	UpdateDeviceData(deviceID int, data float64, house_id int, setting string) error
	UpdataDeviceState(deviceID int, state bool, house_id int, setting string) error
	UpdateManySets([]entity.Set) error
	GetAllNotifications(userID int) ([]entity.Notification, error)
	GetUnreadNotifications(userID int) ([]entity.Notification, error)
	// CreateNotification(userID int, houseId int, notification *entity.Notification) error
	CreateActivityLog(*entity.ActivityLog) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}

// func (s *userUsecase) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
// 	return s.userRepo.CreateUser(ctx, user)
// }

func (s *userUsecase) GetUser(id int) (*entity.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *userUsecase) GetTempAndHumid(house_id int) (float64, float64, error) {
	return s.userRepo.GetTempAndHumid(house_id)
}

// func (s *userUsecase) UpdateUser(ctx context.Context, id string, data *entity.User) (*entity.User, error) {
// 	return s.userRepo.UpdateUser(ctx, id, data)
// }

func (s *userUsecase) AuthenticateUser(username string, password string) (*entity.User, string, []int, error) {
	user, err := s.userRepo.GetUserByUsername(username)

	if err != nil {
		return nil, "", nil, err
	}

	if user == nil {
		return nil, "", nil, entity.ErrUserNotFound
	}
	//skip Verify package: we can use hashed for pw
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password),bcrypt.DefaultCost)

	// func VerifyPassword(password,hashedPassword string) error {
	// 	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	// }

	if user.Password != password {
		return nil, "", nil, entity.ErrUserPasswordNotMatch
	}

	token, err := token.GenerateToken(user.Username + strconv.Itoa(user.ID) + user.Password)

	if err != nil {
		return nil, "", nil, err
	}

	house_ids, err := s.userRepo.GetHouseID(user.ID)

	if err != nil {
		return nil, "", nil, err
	}

	// Change all the character in password of user to * at the same length
	user.Password = strings.Repeat("*", len(user.Password))

	return user, token, house_ids, nil
}

func (s *userUsecase) GetHouseSettingByHouseID(house_id int) ([]entity.HouseSetting, error) {
	return s.userRepo.GetHouseSettingByHouseID(house_id)
}

func (s *userUsecase) GetSetOfHouseSetting(house_id int, settingName string) ([]entity.Set, error) {
	return s.userRepo.GetSetOfHouseSetting(house_id, settingName)
}

func (s *userUsecase) GetActivityLogByHouseID(house_id int) ([]entity.ActivityLog, error) {
	return s.userRepo.GetActivityLogByHouseID(house_id)
}

func (s *userUsecase) UpdateDeviceData(deviceID int, data float64, house_id int, setting string) error {
	return s.userRepo.UpdateDeviceData(deviceID, data, house_id, setting)
}

func (s *userUsecase) UpdataDeviceState(deviceID int, state bool, house_id int, setting string) error {
	return s.userRepo.UpdataDeviceState(deviceID, state, house_id, setting)
}

func (s *userUsecase) GetDashboardData(house_id int) (float64, float64, float64, float64, error) {
	return s.userRepo.GetDashboardData(house_id)
}

func (s *userUsecase) UpdateManySets(sets []entity.Set) error {
	// loop through the sets and update each set
	s.userRepo.UpdateManySets(sets)
	return nil
}

func (s *userUsecase) GetAllNotifications(userID int) ([]entity.Notification, error) {
	return s.userRepo.GetAllNotifications(userID)
}

func (s *userUsecase) GetUnreadNotifications(userID int) ([]entity.Notification, error) {
	return s.userRepo.GetUnreadNotifications(userID)
}

// func (s *userUsecase) CreateNotification(userID int, houseId int, notification *entity.Notification) error {
// 	return s.userRepo.CreateNotification(userID, houseId, notification)
// }

func (s *userUsecase) CreateActivityLog(activityLog *entity.ActivityLog) error {
	return s.userRepo.CreateActivityLog(activityLog)
}
