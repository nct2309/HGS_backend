package usecase

import (
	"context"
	entity "go-jwt/internal/entity"
	repository "go-jwt/internal/infrastructure/repository"
)

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

type UserUsecase interface {
	// CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUser(ctx context.Context, id string) (*entity.User, error)
	// UpdateUser(ctx context.Context, id string, data *entity.User) (*entity.User, error)
	// DeleteUser(ctx context.Context, id string) error
	// AuthenticateUser(ctx context.Context, username string, password string) (*entity.User, string, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

// func (s *userUsecase) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
// 	return s.userRepo.CreateUser(ctx, user)
// }

func (s *userUsecase) GetUser(ctx context.Context, id string) (*entity.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

// func (s *userUsecase) UpdateUser(ctx context.Context, id string, data *entity.User) (*entity.User, error) {
// 	return s.userRepo.UpdateUser(ctx, id, data)
// }

// func (s *userUsecase) DeleteUser(ctx context.Context, id string) error {
// 	return s.userRepo.DeleteUser(ctx, id)
// }

// func (s *userUsecase) AuthenticateUser(ctx context.Context, username string, password string) (*entity.User, string, error) {
// 	user, err := s.userRepo.GetUserByUsername(ctx, username)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	if user == nil {
// 		return nil, "", entity.ERR_USER_NOT_FOUND
// 	}
// 	//skip Verify package: we can use hashed for pw
// 	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password),bcrypt.DefaultCost)

// 	// func VerifyPassword(password,hashedPassword string) error {
// 	// 	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
// 	// }

// 	if user.Password != password {
// 		return nil, "", entity.ERR_USER_PASSWORD_NOT_MATCH
// 	}

// 	token, err := token.GenerateToken(user.ID.String())

// 	if err != nil {
// 		return nil, "", err
// 	}

// 	return user, token, nil
// }
