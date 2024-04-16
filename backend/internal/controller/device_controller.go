package controller

import (
	"encoding/json"
	"go-jwt/internal/middleware"
	request "go-jwt/internal/request"
	usecase "go-jwt/internal/usecase"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeviceController struct {
	deviceService    usecase.DeviceUsecase
	NewDeviceRequest func() request.DeviceRequest
}

func SetupDeviceRoutes(router *gin.Engine, deviceService usecase.DeviceUsecase) {
	deviceController := DeviceController{
		deviceService:    deviceService,
		NewDeviceRequest: request.NewDeviceRequest,
	}

	deviceRoutes := router.Group("/devices")
	{
		deviceRoutes.Use(middleware.CORS())
		deviceRoutes.POST("/updateTemperature", deviceController.UpdateTemperature)
		deviceRoutes.POST("/updateHumidity", deviceController.UpdateHumidity)
		deviceRoutes.POST("/updateFanSpeed", deviceController.UpdateFanSpeed)
		deviceRoutes.GET("/update", deviceController.UpdateDevice)
	}
}

func (h DeviceController) UpdateTemperature(ctx *gin.Context) {
	// Read request body//
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

	// Extract "temp" values from the JSON

	temp, ok := data["temp"].(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing 'temp' value"})
		return
	}

	// Update the temperature
	if err := h.deviceService.UpdateTemperature(1, temp); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update temperature"})
		return
	}
	// Respond with success message
	ctx.JSON(http.StatusOK, gin.H{"message": "Temperature and humidity updated successfully"})

}

func (h DeviceController) UpdateHumidity(ctx *gin.Context) {
	// Read request body
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

	// Extract "humid" value from the JSON
	humid, ok := data["humid"].(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing 'humid' value"})
		return
	}

	// Update the humidity
	if err := h.deviceService.UpdateHumidity(1, humid); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update humidity"})
		return
	}

	// Respond with success message
	ctx.JSON(http.StatusOK, gin.H{"message": "Humidity updated successfully"})
}

func (h DeviceController) UpdateFanSpeed(ctx *gin.Context) {
	// Read request body
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

	// Extract "speed" value from the JSON
	speed, ok := data["speed"].(int)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing 'speed' value"})
		return
	}

	// Update the fan speed
	if err := h.deviceService.UpdateFanSpeed(1, speed); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update fan speed"})
		return
	}

	// Respond with success message
	ctx.JSON(http.StatusOK, gin.H{"message": "Fan speed updated successfully"})
}

func (h DeviceController) UpdateDevice(ctx *gin.Context) {
	// Extract the query parameters from the request
	houseID, deviceID, deviceType, data, state, err := h.NewDeviceRequest().GetDataFromDeviceRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse query parameters"})
		return
	}

	// Update the device
	if err := h.deviceService.UpdateDevice(houseID, deviceID, deviceType, data, state); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update device"})
		return
	}

	// Respond with success message
	ctx.JSON(http.StatusOK, gin.H{"message": "Device updated successfully"})
}
