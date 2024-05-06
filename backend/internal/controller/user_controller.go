package controller

import (
	"encoding/json"
	"fmt"
	"go-jwt/internal/entity"
	"go-jwt/internal/middleware"
	request "go-jwt/internal/request"
	usecase "go-jwt/internal/usecase"
	"io"
	"net/http"
	"strconv"
	"time"

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
		// devices
		userRoutes.POST("/turnOnLight", userController.turnOnLight)
		userRoutes.POST("/turnOffLight", userController.turnOffLight)
		userRoutes.POST("/updateLightLevel", userController.updateLightLevel)
		userRoutes.POST("/turnOnFan", userController.turnOnFan)
		userRoutes.POST("/turnOffFan", userController.turnOffFan)
		userRoutes.POST("/updateFanSpeed", userController.updateFanSpeed)
		userRoutes.GET("/getDashboardData", userController.getDashboardData)
		userRoutes.POST("/openDoor", userController.openDoor)
		userRoutes.POST("/closeDoor", userController.closeDoor)
		// some of the user's house setting
		userRoutes.GET("/getHouseSetting", userController.getHouseSettingByHouseID)
		userRoutes.GET("/getSetOfHouseSetting", userController.getSetOfHouseSetting)
		userRoutes.GET("/getActivityLog", userController.getActivityLogByHouseID)
		userRoutes.POST("/updateSets", userController.updateSets)
		// notifications
		userRoutes.GET("/getAllNotifications", userController.getAlltNotifications)
		userRoutes.GET("/getUnreadNotifications", userController.getUnreadNotifications)
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

	err := h.userService.TurnOnLight(1)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	_ = h.userService.CreateActivityLog(&entity.ActivityLog{
		House_id:      1,
		Device:        "Light",
		Time:          time.Now(),
		Type_of_event: "Turn on the light",
	})

	ctx.JSON(http.StatusOK, gin.H{"message": "Light turned on successfully"})
}

func (h UserController) turnOffLight(ctx *gin.Context) {

	err := h.userService.TurnOffLight(1)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	_ = h.userService.CreateActivityLog(&entity.ActivityLog{
		House_id:      1,
		Device:        "Light",
		Time:          time.Now(),
		Type_of_event: "Turn off the light",
	})

	ctx.JSON(http.StatusOK, gin.H{"message": "Light turned off successfully"})
}

func (h UserController) updateLightLevel(ctx *gin.Context) {
	request := h.NewUserRequest()

	light_level, e := request.GetLightLevel(ctx)

	if e != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return
	}

	err := h.userService.UpdateLightLevel(1, light_level)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	_ = h.userService.CreateActivityLog(&entity.ActivityLog{
		House_id:      1,
		Device:        "Light",
		Time:          time.Now(),
		Type_of_event: "Update the light level to " + strconv.FormatFloat(light_level, 'f', -1, 64),
	})

	ctx.JSON(http.StatusOK, gin.H{"message": "Light level updated successfully"})
}

func (h UserController) updateFanSpeed(ctx *gin.Context) {
	request := h.NewUserRequest()
	fan_speed, e := request.GetFanSpeed(ctx)

	if e != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return
	}

	err := h.userService.UpdateFanSpeed(1, fan_speed)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	_ = h.userService.CreateActivityLog(&entity.ActivityLog{
		House_id:      1,
		Device:        "Fan",
		Time:          time.Now(),
		Type_of_event: "Update the fan speed to " + strconv.FormatFloat(fan_speed, 'f', -1, 64),
	})

	ctx.JSON(http.StatusOK, gin.H{"message": "Fan speed updated successfully"})
}

func (h UserController) turnOnFan(ctx *gin.Context) {

	err := h.userService.TurnOnFan(1)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	_ = h.userService.CreateActivityLog(&entity.ActivityLog{
		House_id:      1,
		Device:        "Fan",
		Time:          time.Now(),
		Type_of_event: "Turn on the fan",
	})

	ctx.JSON(http.StatusOK, gin.H{"message": "Fan turned on successfully"})
}

func (h UserController) turnOffFan(ctx *gin.Context) {

	err := h.userService.TurnOffFan(1)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	_ = h.userService.CreateActivityLog(&entity.ActivityLog{
		House_id:      1,
		Device:        "Fan",
		Time:          time.Now(),
		Type_of_event: "Turn off the fan",
	})

	ctx.JSON(http.StatusOK, gin.H{"message": "Fan turned off successfully"})
}

func (h UserController) openDoor(ctx *gin.Context) {

	err := h.userService.OpenDoor(1)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	_ = h.userService.CreateActivityLog(&entity.ActivityLog{
		House_id:      1,
		Device:        "Door",
		Time:          time.Now(),
		Type_of_event: "Open the door",
	})

	ctx.JSON(http.StatusOK, gin.H{"message": "Door opened successfully"})
}

func (h UserController) closeDoor(ctx *gin.Context) {

	err := h.userService.CloseDoor(1)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	_ = h.userService.CreateActivityLog(&entity.ActivityLog{
		House_id:      1,
		Device:        "Door",
		Time:          time.Now(),
		Type_of_event: "Close the door",
	})

	ctx.JSON(http.StatusOK, gin.H{"message": "Door closed successfully"})
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
	//temp >= 40 và humid <= 15 thì notify
	if temperature >= 40 && humidity <= 15 {
		err := h.userService.CreateNotification(1, 1, &entity.Notification{
			Time:        time.Now(),
			Title:       "Fire Warning!",
			Description: "Temperature: " + res["temperature"] + "°C, Humidity: " + res["humidity"] + "%",
			Read:        false,
		})
		if err != nil {
			fmt.Println("create notification failed:", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "create notification failed", "error": err.Error()})
			return
		}
	}
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

func (h UserController) getAlltNotifications(ctx *gin.Context) {
	request := h.NewUserRequest()
	userID := request.GetUserIDFromURL(ctx)

	notifications, err := h.userService.GetAllNotifications(userID)

	if err != nil {
		fmt.Println("get all notifications failed:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "get all notifications failed", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, notifications)
}

func (h UserController) getUnreadNotifications(ctx *gin.Context) {
	request := h.NewUserRequest()
	userID := request.GetUserIDFromURL(ctx)

	notifications, err := h.userService.GetUnreadNotifications(userID)

	if err != nil {
		fmt.Println("get unread notifications failed:", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "get unread notifications failed", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, notifications)
}
