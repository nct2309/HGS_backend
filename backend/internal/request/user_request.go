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
	GetIDFromURL(c *gin.Context) string
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

func (r *userRequest) GetIDFromURL(c *gin.Context) string {
	return c.Param("id")
}
