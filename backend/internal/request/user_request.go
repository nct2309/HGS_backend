package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-jwt/internal/entity"
	"io"
	"net/http"
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
	GetLightLevel(ctx *gin.Context) (float64, error)
	GetFanSpeed(ctx *gin.Context) (float64, error)
	SendDataToAdafruit(baseURL string, jsonData map[string]string) error
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

func (r *userRequest) GetLightLevel(ctx *gin.Context) (float64, error) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return 0, errors.New("failed to read request body")
	}

	// Define struct to unmarshal JSON into
	var data map[string]interface{}

	// Unmarshal JSON into the struct
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, errors.New("failed to parse JSON")
	}

	// Extract "light_level" values from the JSON
	light_level, ok := data["light_level"].(float64)

	if !ok {
		return 0, errors.New("invalid or missing 'light_level' value")
	}

	// light_level has to be in 1, 2, 3 and 4.
	if !(light_level == 1 || light_level == 2 || light_level == 3 || light_level == 4) {
		return 0, errors.New("invalid 'light_level' value. Must be 1, 2, 3 or 4")
	}

	return light_level, nil
}

func (r *userRequest) GetFanSpeed(ctx *gin.Context) (float64, error) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return 0, errors.New("failed to read request body")
	}

	// Define struct to unmarshal JSON into
	var data map[string]interface{}

	// Unmarshal JSON into the struct
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, errors.New("failed to parse JSON")
	}
	// Extract "fan_speed" values from the JSON
	fan_speed, ok := data["fan_speed"].(float64)
	if !ok {
		return 0, errors.New("invalid or missing 'fan_speed' value")
	}

	// fan_speed has to be between 0 and 100
	if fan_speed < 0 || fan_speed > 100 {
		return 0, errors.New("invalid 'fan_speed' value. Must be between 0 and 100")
	}

	return fan_speed, nil
}

func (r *userRequest) SendDataToAdafruit(baseURL string, jsonData map[string]string) error {
	// Convert JSON data to bytes
	jsonDataBytes, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return err
	}

	// Create a new HTTP client
	client := &http.Client{}

	// Create a POST request with the JSON data
	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(jsonDataBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}

	// Print the response status code and body
	fmt.Println("Response body:", string(body))

	return nil
}
