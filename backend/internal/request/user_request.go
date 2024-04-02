package request

import (
	"go-jwt/internal/entity"

	"github.com/gin-gonic/gin"
)

func NewUserRequest() UserRequest {
	return &userRequest{}
}

type UserRequest interface {
	Bind(c *gin.Context) error
	//Validate() error
	// GetID() string
	GetIDFromURL(c *gin.Context) string
	// GetName() string
	// GetUsername() string
	// GetPassword() string
	// GetPhonenum() string
	// GetAge() int
	// GetSSN() string
	// GetRole() int
	// GetCountFine() int
	// GetBookIDFromURL(c *gin.Context) string
	// GetGender() string
	// GetReservingList() []entity.UserActivity
	// GetBorrowingList() []entity.UserActivity
	// GetBorrowedList() []entity.UserActivity
	// GetToken() string
	// GetRefreshToken() string
	// GetExpireTime() int64
	// GetRefreshExpireTime() int64
}

type userRequest struct {
	user entity.User
}

func (r *userRequest) Bind(c *gin.Context) error {
	return c.ShouldBindJSON(&r.user)
}

// func (r *userRequest) GetName() string {
// 	return r.user.Name
// }

// func (r *userRequest) GetID() string {
// 	return r.user.ID.Hex()
// }

func (r *userRequest) GetIDFromURL(c *gin.Context) string {
	return c.Param("id")
}

// func (r *userRequest) GetUsername() string {
// 	return r.user.Username
// }

// func (r *userRequest) GetPassword() string {
// 	return r.user.Password
// }

// func (r *userRequest) GetPhonenum() string {
// 	return r.user.Phonenum
// }

// func (r *userRequest) GetAge() int {
// 	return r.user.Age
// }

// func (r *userRequest) GetSSN() string {
// 	return r.user.SSN
// }

// func (r *userRequest) GetRole() int {
// 	return r.user.Role
// }

// func (r *userRequest) GetCountFine() int {
// 	return r.user.CountFine
// }

// func (r *userRequest) GetBookIDFromURL(c *gin.Context) string {
// 	return c.Param("bookId")
// }

// func (r *userRequest) GetGender() string {
// 	return r.user.Gender
// }

// func (r *userRequest) GetReservingList() []entity.UserActivity {
// 	return r.user.ReservingList
// }

// func (r *userRequest) GetBorrowingList() []entity.UserActivity {
// 	return r.user.BorrowingList
// }

// func (r *userRequest) GetBorrowedList() []entity.UserActivity {
// 	return r.user.BorrowedList
// }

// func (r *userRequest) GetToken() string {
// 	return r.user.Token
// }
//

// func (r *userRequest) GetRefreshToken() string {
// 	return r.user.RefreshToken
// }

// func (r *userRequest) GetExpireTime() int64 {
// 	return r.user.ExpireTime
// }

// func (r *userRequest) GetRefreshExpireTime() int64 {
// 	return r.user.RefreshExpireTime
// }

// func (r *userRequest) Validate() error {
// 	if r.user.Name == "" {
// 		return entity.ERR_INVALID_NAME
// 	}
// 	if r.user.Username == "" {
// 		return entity.ERR_INVALID_USERNAME
// 	}
// 	if r.user.Password == "" {
// 		return entity.ERR_INVALID_PASSWORD
// 	}
// 	if r.user.Phonenum == "" {
// 		return entity.ERR_INVALID_PHONENUM
// 	}
// 	if r.user.Age == 0 {
// 		return entity.ERR_INVALID_AGE
// 	}
// 	if r.user.SSN == "" {
// 		return entity.ERR_INVALID_SSN
// 	}
// 	if r.user.Role == 0 {
// 		return entity.ERR_INVALID_ROLE
// 	}
// 	if r.user.CountFine == 0 {
// 		return entity.ERR_INVALID_COUNTFINE
// 	}
// 	return nil
// }
