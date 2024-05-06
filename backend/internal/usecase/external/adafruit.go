package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type AdaFruitService struct {
	ExternalService
}

func (a *AdaFruitService) Execute(strategy Strategy, des any) error {
	return strategy.Execute(des)
}

type LightOn struct {
}

type LightOff struct {
}

type LightLevel struct {
	LightLevel float64
}

type FanSpeed struct {
	FanSpeed float64
}

type FanOn struct {
}

type FanOff struct {
}

type DoorOpen struct {
}

type DoorClose struct {
}

func (l *LightOn) Execute(des any) error {
	jsonData := map[string]string{
		"value": "Alarm On",
	}
	// Send request to Adafruit to turn on the light
	err := SendRequest("https://io.adafruit.com/api/v2/webhooks/feed/Ye9oEbz9VvPgzjLYzjz7dDC8R1dL", jsonData, des)
	if err != nil {
		return err
	}
	return nil
}

func (l *LightOff) Execute(des any) error {
	jsonData := map[string]string{
		"value": "Alarm Off",
	}
	// Send request to Adafruit to turn off the light
	err := SendRequest("https://io.adafruit.com/api/v2/webhooks/feed/Ye9oEbz9VvPgzjLYzjz7dDC8R1dL", jsonData, des)
	if err != nil {
		return err
	}
	return nil
}

func (l *LightLevel) Execute(des any) error {
	jsonData := map[string]string{
		"value": strconv.FormatFloat(l.LightLevel, 'f', -1, 64),
	}

	// Send request to Adafruit to set the light level
	err := SendRequest("https://io.adafruit.com/api/v2/webhooks/feed/YUgssBNR6j1J24jF6RDYG71QqH4c", jsonData, des)
	if err != nil {
		return err
	}
	return nil
}

func (f *FanSpeed) Execute(des any) error {
	jsonData := map[string]string{
		"value": strconv.FormatFloat(f.FanSpeed, 'f', -1, 64),
	}

	// Send request to Adafruit to set the fan speed
	err := SendRequest("https://io.adafruit.com/api/v2/webhooks/feed/GDfmkBYDyWBUV6A6M17stLHytSEM", jsonData, des)
	if err != nil {
		return err
	}
	return nil
}

func (f *FanOn) Execute(des any) error {
	jsonData := map[string]string{
		"value": "Fan On",
	}

	// Send request to Adafruit to turn on the fan
	err := SendRequest("https://io.adafruit.com/api/v2/webhooks/feed/9xJ4R9ZM7A9tKEeJcaJh9rS7t6L5", jsonData, des)
	if err != nil {
		return err
	}
	return nil
}

func (f *FanOff) Execute(des any) error {
	jsonData := map[string]string{
		"value": "Fan Off",
	}

	// Send request to Adafruit to turn off the fan
	err := SendRequest("https://io.adafruit.com/api/v2/webhooks/feed/9xJ4R9ZM7A9tKEeJcaJh9rS7t6L5", jsonData, des)
	if err != nil {
		return err
	}
	return nil
}

func (d *DoorOpen) Execute(des any) error {
	jsonData := map[string]string{
		"value": "Open Door",
	}

	// Send request to Adafruit to open the door
	err := SendRequest("https://io.adafruit.com/api/v2/webhooks/feed/iMQFZUbNRJPM5ZCzRN4ped6GbL4W", jsonData, des)
	if err != nil {
		return err
	}
	return nil
}

func (d *DoorClose) Execute(des any) error {
	jsonData := map[string]string{
		"value": "Close Door",
	}

	// Send request to Adafruit to close the door
	err := SendRequest("https://io.adafruit.com/api/v2/webhooks/feed/iMQFZUbNRJPM5ZCzRN4ped6GbL4W", jsonData, des)
	if err != nil {
		return err
	}
	return nil
}

func SendRequest(baseURL string, jsonData map[string]string, des any) error {
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

	if des != nil {
		err = json.Unmarshal(body, &des)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return err
		}
	}

	return nil
}
