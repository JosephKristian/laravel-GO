package service

import (
	"errors"
)

type AccountActivationService struct {
	registerService *RegisterService
}

func NewAccountActivationService() *AccountActivationService {
	return &AccountActivationService{
		registerService: &RegisterService{}, // Inject RegisterService dependency
	}
}

type ActivationResult struct {
	StatusCode int
	Message    string
}

func (s *AccountActivationService) Activate(emailOrPhone string, verificationCode int, clientIP string) (*ActivationResult, error) {
	tx, err := s.registerService.BeginTransaction()
	if err != nil {
		return &ActivationResult{StatusCode: 500, Message: "Failed to start transaction"}, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	userConfirm, err := s.registerService.GetUserConfirm(emailOrPhone)
	if err != nil || userConfirm == nil {
		tx.Rollback()
		return &ActivationResult{StatusCode: 400, Message: "Invalid verification code"}, errors.New("user confirmation not found")
	}

	isValid, err := s.registerService.VerifyActivationCode(userConfirm, verificationCode)
	if err != nil || !isValid {
		tx.Rollback()
		return &ActivationResult{StatusCode: 400, Message: "Invalid verification code"}, errors.New("verification failed")
	}

	userActivation, err := s.registerService.ActivateAccount(userConfirm)
	if err != nil || userActivation == nil {
		tx.Rollback()
		return &ActivationResult{StatusCode: 400, Message: "Account activation failed"}, err
	}

	if userConfirm.DeviceID != "" {
		err = s.registerService.AddDevice(userActivation.ID, userConfirm.DeviceID, userConfirm.Device, clientIP)
		if err != nil {
			tx.Rollback()
			return &ActivationResult{StatusCode: 500, Message: "Failed to add device"}, err
		}
	}

	err = s.registerService.DeleteUserConfirm(userConfirm.ID)
	if err != nil {
		tx.Rollback()
		return &ActivationResult{StatusCode: 500, Message: "Failed to delete confirmation data"}, err
	}

	if err := tx.Commit(); err != nil {
		return &ActivationResult{StatusCode: 500, Message: "Transaction commit failed"}, err
	}

	return &ActivationResult{StatusCode: 200, Message: "Account successfully activated"}, nil
}
