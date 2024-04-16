package request

import (
	"go-jwt/internal/entity"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewUserRequest() UserRequest {
	return &userRequest{}
}

type UserRequest interface {
	Bind(c *gin.Context) error
	GetIDFromURL(c *gin.Context) string
	GetUsername() string
	GetPassword() string
	GetHouseIDFromURL(ctx *gin.Context) int
	GetHouseSettingNameFromURL(ctx *gin.Context) string
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

func (r *userRequest) GetUsername() string {
	// please write some code to escape the username
	// escape ' or " or ; or --
	// return the escaped username
	if r.user.Username == "" {
		return ""
	}

	// Escape ' or " or ; or --
	username := strings.ReplaceAll(r.user.Username, "'", "\\'")
	username = strings.ReplaceAll(username, "\"", "\\\"")
	username = strings.ReplaceAll(username, ";", "\\;")
	username = strings.ReplaceAll(username, "--", "\\--")

	return username
}

func (r *userRequest) GetPassword() string {
	// please write some code to escape the password
	// escape ' or " or ; or --
	// return the escaped password
	if r.user.Password == "" {
		return ""
	}

	// Escape ' or " or ; or --
	password := strings.ReplaceAll(r.user.Password, "'", "\\'")
	password = strings.ReplaceAll(password, "\"", "\\\"")
	password = strings.ReplaceAll(password, ";", "\\;")
	password = strings.ReplaceAll(password, "--", "\\--")

	return password
}

// /users/getHouseSettingByHouseID?house_id=1
func (r *userRequest) GetHouseIDFromURL(ctx *gin.Context) int {
	houseID, _ := ctx.GetQuery("house_id")
	house_id, _ := strconv.Atoi(houseID)
	return house_id
}

func (r *userRequest) GetHouseSettingNameFromURL(ctx *gin.Context) string {
	return ctx.Query("name")
}
