package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-jwt/internal/middleware"
	request "go-jwt/internal/request"
	usecase "go-jwt/internal/usecase"
	"io"
	"mime/multipart"
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
		deviceRoutes.POST("/setFace", deviceController.UploadImage)
		deviceRoutes.POST("/verifyFace", deviceController.VerifyFace)
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

func (h DeviceController) UploadImage(ctx *gin.Context) {

	encodeURL := "https://face-reg-service-latest.onrender.com/img2encoding"

	// Extract the image file from the request
	file, err := ctx.FormFile("img")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read image file"})
		return
	}

	// send the file to the face recognition service in encodeURL in POST request
	client := &http.Client{}

	// Create a new form data
	formData := new(bytes.Buffer)
	writer := multipart.NewWriter(formData)
	// send in form-data with the file
	part, err := writer.CreateFormFile("img", file.Filename)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create form file"})
		return
	}

	// Open the file
	fileHandle, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open file"})
		return
	}
	defer fileHandle.Close()

	// Copy the file to the form data
	_, err = io.Copy(part, fileHandle)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to copy file to form data"})
		return
	}

	// Close the writer
	writer.Close()

	// Create a new POST request
	req, err := http.NewRequest(http.MethodPost, encodeURL, formData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create POST request"})
		return
	}

	// Set the content type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to send POST request"})
		return
	}

	// Get the string from the body in the key "face_encoding"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read response body"})
		return
	}

	// Define struct to unmarshal JSON into
	var data map[string]interface{}

	// Unmarshal JSON into the struct
	if err := json.Unmarshal(body, &data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON"})
		return
	}

	// Extract "face_encoding" value from the JSON
	faceEncode, ok := data["face_encoding"].(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing 'face_encoding' value"})
		return
	}

	// Update the face encodings
	if err := h.deviceService.UpdateFaceEncodings(1, faceEncode); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update face encodings"})
		return
	}

	// Respond with success message
	ctx.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully"})
}

func (h DeviceController) VerifyFace(ctx *gin.Context) {

	verifyURL := "https://face-reg-service-latest.onrender.com/verify"

	// Extract the image file from the request
	file, err := ctx.FormFile("img")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read image file"})
		return
	}

	// send the file to the face recognition service in encodeURL in POST request
	client := &http.Client{}

	// Create a new form data
	formData := new(bytes.Buffer)
	writer := multipart.NewWriter(formData)
	// send in form-data with the file
	part, err := writer.CreateFormFile("img", file.Filename)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create form file"})
		return
	}

	// Open the file
	fileHandle, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open file"})
		return
	}
	defer fileHandle.Close()

	// Copy the file to the form data
	_, err = io.Copy(part, fileHandle)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to copy file to form data"})
		return
	}

	// add "encoding_array" to the form data by taking the value from db
	face_encodings, err := h.deviceService.GetFaceEncoding(1)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get face encodings from database"})
		return
	}

	// Add the face encodings to "encoding_array" key as an array of strings
	for _, encoding := range face_encodings {
		if err := writer.WriteField("encoding_array", encoding); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to add face encodings to form data"})
			return
		}
	}

	// Close the writer
	writer.Close()

	// Create a new POST request
	req, err := http.NewRequest(http.MethodPost, verifyURL, formData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create POST request"})
		return
	}
	//multipart/form-data;
	// Set the content type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to send POST request"})
		return
	}

	// 	if there is not a face
	// {
	//   "error": "No face found in the image"
	// }
	// if there is a face.
	// {
	//   "is_match": true
	// }

	// Get the is_match or error from the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read response body"})
		return
	}

	// Define struct to unmarshal JSON into
	var data map[string]interface{}

	// Unmarshal JSON into the struct
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(err, data, "\n", string(body))
		//err print "invalid character '<' looking for beginning of value"
		//how to fix this know that there are only 2 kind of JSON response

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON"})
		return
	}

	// Extract "is_match" value from the JSON
	isMatch, ok := data["is_match"].(bool)
	if !ok {
		// take the error message
		errorMessage, ok := data["error"].(string)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Something wrong with the response"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Face verified successfully", "is_match": isMatch})
}
