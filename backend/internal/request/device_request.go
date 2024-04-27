package request

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NewDeviceRequest() DeviceRequest {
	return &deviceRequest{}
}

type DeviceRequest interface {
	GetDataFromDeviceRequest(c *gin.Context) (int, int, string, float64, bool, error)
}

type deviceRequest struct {
}

//api route: /devices/update?house_id=1&device_id=2&device_type=Door&data=30&state=1

func (r *deviceRequest) GetDataFromDeviceRequest(ctx *gin.Context) (int, int, string, float64, bool, error) {
	// read the request url to extract the query parameters
	HouseID, _ := ctx.GetQuery("house_id")
	DeviceID, _ := ctx.GetQuery("device_id")
	deviceType, _ := ctx.GetQuery("device_type")
	Data, _ := ctx.GetQuery("data")
	State, _ := ctx.GetQuery("state")

	// convert the query parameters to the desired data type
	houseID, _ := strconv.Atoi(HouseID)
	deviceID, _ := strconv.Atoi(DeviceID)
	data, _ := strconv.ParseFloat(Data, 64)
	state, _ := strconv.ParseBool(State)

	if data < 0 {
		return 0, 0, "", 0, false, errors.New("data cannot be negative")
	}

	if deviceType == "Door" {

	} else if deviceType == "Temperature" {
		if data > 100 {
			return 0, 0, "", 0, false, errors.New("fire detected! temperature cannot be more than 100")
		}
	} else if deviceType == "Humidity" {

	} else if deviceType == "Fan" {
		if data > 100 {
			return 0, 0, "", 0, false, errors.New("fan speed cannot be more than 100")
		}
		// } else if deviceType == "Light" {

		// }
	}
	return houseID, deviceID, deviceType, data, state, nil
}
