package usecase

import (
	repository "go-jwt/internal/infrastructure/repository"
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
