package usecase

import (
	entity "go-jwt/internal/entity"
	repository "go-jwt/internal/infrastructure/repository"
	"go-jwt/internal/token"
	"strconv"
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
	AuthenticateUser(username string, password string) (*entity.User, string, error)
	GetTempAndHumid(house_id int) (float64, float64, error)
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

// func (s *userUsecase) DeleteUser(ctx context.Context, id string) error {
// 	return s.userRepo.DeleteUser(ctx, id)
// }

func (s *userUsecase) AuthenticateUser(username string, password string) (*entity.User, string, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, "", err
	}

	if user == nil {
		return nil, "", entity.ERR_USER_NOT_FOUND
	}
	//skip Verify package: we can use hashed for pw
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password),bcrypt.DefaultCost)

	// func VerifyPassword(password,hashedPassword string) error {
	// 	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	// }

	if user.Password != password {
		return nil, "", entity.ERR_USER_PASSWORD_NOT_MATCH
	}

	token, err := token.GenerateToken(user.Username + strconv.Itoa(user.ID) + user.Password)

	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
