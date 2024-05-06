package usecase

import (
	"bytes"
	entity "go-jwt/internal/entity"
	repository "go-jwt/internal/infrastructure/repository"
	external "go-jwt/internal/usecase/external"
)

func NewDeviceUsecase(deviceRepo repository.DeviceRepository) DeviceUsecase {
	return &deviceUsecase{
		deviceRepo: deviceRepo,
	}
}

type DeviceUsecase interface {
	UpdateTemperature(id int, temperature float64) error
	UpdateHumidity(id int, humid float64) error
	UpdateFanSpeed(id int, speed int) error
	UpdateDevice(houseID int, deviceID int, deviceType string, data float64, state bool) error
	UpdateFaceEncodings(houseID int, faceEncode string) error
	GetFaceEncoding(houseID int) ([]string, error)
	EncodeFace(houseID int, formData *bytes.Buffer, ContentType string, data *map[string]interface{}) error
	VerifyFace(houseID int, formData *bytes.Buffer, ContentType string, data *map[string]interface{}) error
	OpenDoorAfterFaceVerified(houseID int) error
	CreateActivityLog(*entity.ActivityLog) error
}

type deviceUsecase struct {
	deviceRepo repository.DeviceRepository
}

func (s *deviceUsecase) UpdateTemperature(id int, temperature float64) error {
	return s.deviceRepo.UpdateTemperature(id, temperature)
}

func (s *deviceUsecase) UpdateHumidity(id int, humid float64) error {
	return s.deviceRepo.UpdateHumidity(id, humid)
}

func (s *deviceUsecase) UpdateFanSpeed(id int, speed int) error {
	return s.deviceRepo.UpdateFanSpeed(id, speed)
}

func (s *deviceUsecase) UpdateDevice(houseID int, deviceID int, deviceType string, data float64, state bool) error {
	return s.deviceRepo.UpdateDevice(houseID, deviceID, deviceType, data, state)
}

func (s *deviceUsecase) UpdateFaceEncodings(houseID int, faceEncode string) error {
	return s.deviceRepo.UpdateFaceEncodings(houseID, faceEncode)
}

func (s *deviceUsecase) GetFaceEncoding(houseID int) ([]string, error) {
	return s.deviceRepo.GetFaceEncoding(houseID)
}

func (s *deviceUsecase) EncodeFace(houseID int, formData *bytes.Buffer, ContentType string, data *map[string]interface{}) error {
	return external.NewExternalServiceAdapter(&external.FaceRecognitionService{}).Execute(&external.EncodeFace{FormData: formData, ContentType: ContentType}, data)
}

func (s *deviceUsecase) VerifyFace(houseID int, formData *bytes.Buffer, ContentType string, data *map[string]interface{}) error {
	return external.NewExternalServiceAdapter(&external.FaceRecognitionService{}).Execute(&external.VerifyFace{FormData: formData, ContentType: ContentType}, data)
}

func (s *deviceUsecase) OpenDoorAfterFaceVerified(houseID int) error {
	return external.NewExternalServiceAdapter(&external.AdaFruitService{}).Execute(&external.DoorOpen{}, nil)
}

func (s *deviceUsecase) CreateActivityLog(activityLog *entity.ActivityLog) error {
	return s.deviceRepo.CreateActivityLog(activityLog)
}
