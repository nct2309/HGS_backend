package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-jwt/internal/entity"
	"go-jwt/internal/middleware"
	request "go-jwt/internal/request"
	usecase "go-jwt/internal/usecase"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService    usecase.UserUsecase
	NewUserRequest func() request.UserRequest
}

func SetupUserRoutes(router *gin.Engine, userService usecase.UserUsecase) {
	userController := UserController{
		userService:    userService,
		NewUserRequest: request.NewUserRequest,
	}

	publicRoutes := router.Group("/public")
	{
		publicRoutes.Use(middleware.CORS())
		publicRoutes.POST("/login", userController.login)
		// publicRoutes.POST("/", userController.create)
	}

	userRoutes := router.Group("/users").Use(middleware.JwtAuthMiddleware())
	{
		userRoutes.Use(middleware.CORS())
		userRoutes.GET("/:id", userController.get)
		userRoutes.POST("/turnOnLight", userController.turnOnLight)
		userRoutes.POST("/turnOffLight", userController.turnOffLight)
		userRoutes.POST("/updateLightLevel", userController.updateLightLevel)
		userRoutes.POST("/turnOnFan", userController.turnOnFan)
		userRoutes.POST("/turnOffFan", userController.turnOffFan)
		userRoutes.POST("/updateFanSpeed", userController.updateFanSpeed)
		userRoutes.GET("/getDashboardData", userController.getDashboardData)
		userRoutes.GET("/getHouseSetting", userController.getHouseSettingByHouseID)
		userRoutes.GET("/getSetOfHouseSetting", userController.getSetOfHouseSetting)
		userRoutes.GET("/getActivityLog", userController.getActivityLogByHouseID)
		userRoutes.POST("/updateSets", userController.updateSets)
	}
}

func (h UserController) login(ctx *gin.Context) {
	request := h.NewUserRequest()

	if err := request.Bind(ctx); err != nil {
		fmt.Println("bind user failed:", err.Error())
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, token, house_ids, err := h.userService.AuthenticateUser(request.GetUsername(), request.GetPassword())

	if err != nil {
		fmt.Println("login user failed:", err.Error())
		// 404 not found http status code
		ctx.JSON(http.StatusNotFound, gin.H{"message": "login failed", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token, "user": user, "house_ids": house_ids})
}

func (h UserController) get(ctx *gin.Context) {

	request := h.NewUserRequest()
	id, err := strconv.Atoi(request.GetIDFromURL(ctx))
	if err != nil {
		fmt.Println("get user failed:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "get failed", "error": err.Error()})
		return
	}
	user, err := h.userService.GetUser(id)

	if err != nil {
		fmt.Println("get user failed:", err.Error())
		// ctx.AbortWithError(http.StatusBadRequest, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "get failed", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h UserController) turnOnLight(ctx *gin.Context) {

	// Build the API endpoint URL
	baseURL := "https://io.adafruit.com/api/v2/webhooks/feed/Ye9oEbz9VvPgzjLYzjz7dDC8R1dL"

	jsonData := map[string]string{
		"value": "Alarm On",
	}

	// Convert JSON data to bytes
	jsonDataBytes, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Create a new HTTP client
	client := &http.Client{}

	// Create a POST request with the JSON data
	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(jsonDataBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the response status code and body
	fmt.Println("Response body:", string(body))

	ctx.JSON(http.StatusOK, gin.H{"message": "Light turned on successfully"})
}

func (h UserController) turnOffLight(ctx *gin.Context) {

	// Build the API endpoint URL
	baseURL := "https://io.adafruit.com/api/v2/webhooks/feed/Ye9oEbz9VvPgzjLYzjz7dDC8R1dL"

	// Create the data you want to send
	jsonData := map[string]string{
		"value": "Alarm Off",
	}

	// Convert JSON data to bytes
	jsonDataBytes, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Create a new HTTP client
	client := &http.Client{}

	// Create a POST request with the JSON data
	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(jsonDataBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the response body
	fmt.Println("Response body:", string(body))

	ctx.JSON(http.StatusOK, gin.H{"message": "Light turned off successfully"})
}

func (h UserController) updateLightLevel(ctx *gin.Context) {

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Define struct to unmarshal JSON into
	var data map[string]interface{}

	// Unmarshal JSON into the struct
	if err := json.Unmarshal(body, &data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON"})
		return
	}

	// Extract "light_level" values from the JSON
	light_level, ok := data["light_level"].(float64)

	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing 'light_level' value"})
		return
	}

	// light_level has to be in 1, 2, 3 and 4.
	if !(light_level == 1 || light_level == 2 || light_level == 3 || light_level == 4) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'light_level' value. Must be 1, 2, 3 or 4"})
		return
	}

	baseURL := "https://io.adafruit.com/api/v2/webhooks/feed/YUgssBNR6j1J24jF6RDYG71QqH4c"

	// Create the data you want to send
	jsonData := map[string]string{
		"value": strconv.FormatFloat(light_level, 'f', -1, 64),
	}

	// Convert JSON data to bytes
	jsonDataBytes, err := json.Marshal(jsonData)

	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Create a new HTTP client
	client := &http.Client{}

	// Create a POST request with the JSON data
	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(jsonDataBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()

	// Read the response body
	bodyResponse, err1 := io.ReadAll(resp.Body)
	if err1 != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the response body
	fmt.Println("Response body:", string(bodyResponse))

	ctx.JSON(http.StatusOK, gin.H{"message": "Light level updated successfully"})
}

func (h UserController) updateFanSpeed(ctx *gin.Context) {

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Define struct to unmarshal JSON into
	var data map[string]interface{}

	// Unmarshal JSON into the struct
	if err := json.Unmarshal(body, &data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON"})
		return
	}
	// Extract "fan_speed" values from the JSON
	fan_speed, ok := data["fan_speed"].(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing 'fan_speed' value"})
		return
	}

	// fan_speed has to be between 0 and 100
	if fan_speed < 0 || fan_speed > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'fan_speed' value. Must be between 0 and 100"})
		return
	}

	// Build the API endpoint URL
	baseURL := "https://io.adafruit.com/api/v2/webhooks/feed/GDfmkBYDyWBUV6A6M17stLHytSEM"

	// Create the data you want to send
	jsonData := map[string]string{
		"value": strconv.FormatFloat(fan_speed, 'f', -1, 64),
	}

	// Convert JSON data to bytes
	jsonDataBytes, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Create a new HTTP client
	client := &http.Client{}

	// Create a POST request with the JSON data
	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(jsonDataBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()

	// Read the response body
	bodyResponse, err1 := io.ReadAll(resp.Body)
	if err1 != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the response body
	fmt.Println("Response body:", string(bodyResponse))

	ctx.JSON(http.StatusOK, gin.H{"message": "Fan speed updated successfully"})
}

func (h UserController) turnOnFan(ctx *gin.Context) {

	// Build the API endpoint URL
	baseURL := "https://io.adafruit.com/api/v2/webhooks/feed/9xJ4R9ZM7A9tKEeJcaJh9rS7t6L5"

	// Create the data you want to send
	jsonData := map[string]string{
		"value": "Fan On",
	}

	// Convert JSON data to bytes
	jsonDataBytes, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Create a new HTTP client
	client := &http.Client{}

	// Create a POST request with the JSON data
	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(jsonDataBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()

	// Read the response body
	bodyResponse, err1 := io.ReadAll(resp.Body)
	if err1 != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the response body
	fmt.Println("Response body:", string(bodyResponse))

	ctx.JSON(http.StatusOK, gin.H{"message": "Fan turned on successfully"})
}

func (h UserController) turnOffFan(ctx *gin.Context) {

	// Build the API endpoint URL
	baseURL := "https://io.adafruit.com/api/v2/webhooks/feed/9xJ4R9ZM7A9tKEeJcaJh9rS7t6L5"

	// Create the data you want to send
	jsonData := map[string]string{
		"value": "Fan Off",
	}

	// Convert JSON data to bytes
	jsonDataBytes, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Create a new HTTP client
	client := &http.Client{}

	// Create a POST request with the JSON data
	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(jsonDataBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()

	// Read the response body
	bodyResponse, err1 := io.ReadAll(resp.Body)
	if err1 != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the response body
	fmt.Println("Response body:", string(bodyResponse))

	ctx.JSON(http.StatusOK, gin.H{"message": "Fan turned off successfully"})
}

func (h UserController) getDashboardData(ctx *gin.Context) {

	// temperature, humid, light, fan_speed, err := h.userService.GetDashboardData(1)
	lightURL := "https://io.adafruit.com/api/v2/QuangThien15/feeds/iot-alarm/data"
	lightLevelURL := "https://io.adafruit.com/api/v2/QuangThien15/feeds/iot-state/data"
	fanURL := "https://io.adafruit.com/api/v2/QuangThien15/feeds/iot-fan/data"
	fanSpeedURL := "https://io.adafruit.com/api/v2/QuangThien15/feeds/iot-fanspeed/data"
	doorURL := "https://io.adafruit.com/api/v2/QuangThien15/feeds/iot-door/data"
	temperatureURL := "https://io.adafruit.com/api/v2/QuangThien15/feeds/iot-temperature/data"
	humidityURL := "https://io.adafruit.com/api/v2/QuangThien15/feeds/iot-humidity/data"

	URLmaps := map[string]string{
		"light":       lightURL,
		"light_level": lightLevelURL,
		"fan":         fanURL,
		"fan_speed":   fanSpeedURL,
		"door":        doorURL,
		"temperature": temperatureURL,
		"humidity":    humidityURL,
	}

	res := make(map[string]string)

	// Create a new HTTP client
	client := &http.Client{}

	for key, URL := range URLmaps {
		// Create a GET request
		req, err := http.NewRequest(http.MethodGet, URL, nil)

		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		// Send the request
		resp, err := client.Do(req)

		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}

		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		// take the "value" from the JSON response

		var data []map[string]interface{}

		// Unmarshal JSON into the struct know that data is an array of JSON objects
		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		// Extract "value" values from the first JSON object
		value, ok := data[0]["value"].(string)
		if !ok {
			fmt.Println("Error getting value from JSON")
			return
		}

		res[key] = value
	}

	temperature, _ := strconv.ParseFloat(res["temperature"], 64)
	humidity, _ := strconv.ParseFloat(res["humidity"], 64)

	var light bool
	if res["light"] == "Alarm On" {
		light = true
	} else {
		light = false
	}

	var fan bool
	if res["fan"] == "Fan On" {
		fan = true
	} else {
		fan = false
	}

	var door bool
	if res["door"] == "Open Door" {
		door = true
	} else {
		door = false
	}

	light_level, _ := strconv.ParseFloat(res["light_level"], 64)
	fan_speed, _ := strconv.ParseFloat(res["fan_speed"], 64)

	ctx.JSON(http.StatusOK, gin.H{
		"temperature": temperature,
		"humidity":    humidity,
		"light":       light,
		"fan":         fan,
		"door":        door,
		"light_level": light_level,
		"fan_speed":   fan_speed,
	})
}

func (h UserController) getHouseSettingByHouseID(ctx *gin.Context) {
	request := h.NewUserRequest()
	house_id := request.GetHouseIDFromURL(ctx)

	houseSetting, err := h.userService.GetHouseSettingByHouseID(house_id)

	if err != nil {
		fmt.Println("get house setting failed:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "get house setting failed", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, houseSetting)
}

func (h UserController) getSetOfHouseSetting(ctx *gin.Context) {
	request := h.NewUserRequest()
	house_id := request.GetHouseIDFromURL(ctx)
	settingName := request.GetHouseSettingNameFromURL(ctx)

	sets, err := h.userService.GetSetOfHouseSetting(house_id, settingName)

	if err != nil {
		fmt.Println("get set of house setting failed:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "get set of house setting failed", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, sets)
}

// /users/getActivityLog?house_id=1
func (h UserController) getActivityLogByHouseID(ctx *gin.Context) {
	request := h.NewUserRequest()
	house_id := request.GetHouseIDFromURL(ctx)

	activityLog, err := h.userService.GetActivityLogByHouseID(house_id)

	if err != nil {
		fmt.Println("get activity log failed:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "get activity log failed", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, activityLog)
}

func (h UserController) updateSets(ctx *gin.Context) {

	// map to store the data to Set
	var Sets []entity.Set
	// use ShouldBindJSON
	if err := ctx.ShouldBindJSON(&Sets); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// call the usecase
	err := h.userService.UpdateManySets(Sets)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Sets updated successfully"})
}
