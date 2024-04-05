package request

import (
	"go-jwt/internal/entity"

	"github.com/gin-gonic/gin"
)

func NewDeviceRequest() DeviceRequest {
	return &deviceRequest{}
}

type DeviceRequest interface {
	Bind(c *gin.Context) error
	GetIDFromURL(c *gin.Context) string
}

type deviceRequest struct {
	user entity.User
}

func (r *deviceRequest) Bind(c *gin.Context) error {
	return c.ShouldBindJSON(&r.user)
}

func (r *deviceRequest) GetIDFromURL(c *gin.Context) string {
	return c.Param("id")
}
