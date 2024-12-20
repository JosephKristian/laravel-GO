package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/JosephKristian/project-migration/internal/models"
	"github.com/JosephKristian/project-migration/internal/repositories"
	"github.com/JosephKristian/project-migration/internal/usecase/helpers"
	"github.com/google/uuid"
)

type RegisterService struct {
	DB         *sql.DB
	userRepo   repositories.UserRepo
	otpService *OtpService
}

func NewRegisterService(db *sql.DB, userRepo repositories.UserRepo, otpService *OtpService) *RegisterService {
	return &RegisterService{
		DB:         db,
		userRepo:   userRepo,
		otpService: otpService,
	}
}

func (r *RegisterService) Register(data *models.RegisterInput) (*models.User, error) {
	// Validasi input
	if data.Name == "" {
		return nil, errors.New("name is required")
	}

	if data.Email != "" && !helpers.IsValidEmail(data.Email) {
		return nil, errors.New("invalid email format")
	}

	isEmailExist, err := r.userRepo.CheckEmailExist(data.Email)
	if err != nil {
		return nil, errors.New("failed to check email existence: " + err.Error())
	}
	if isEmailExist {
		return nil, errors.New("email already registered")
	}

	if data.Phone == "" {
		return nil, errors.New("phone is required")
	}
	if len(data.Phone) < 8 || len(data.Phone) > 15 {
		return nil, errors.New("phone number must be between 8 and 15 digits")
	}

	if data.Password == "" {
		return nil, errors.New("password is required")
	} else if len(data.Password) < 8 || !helpers.ContainsRequiredChars(data.Password) {
		return nil, errors.New("password must be at least 8 characters, include uppercase, lowercase, number, and special character")
	}

	if data.VerificationChannel == "" {
		return nil, errors.New("verification_channel is required")
	}

	// Hashing password
	hashedPassword, err := helpers.HashPassword(data.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// UUID dan OTP
	userUUID := uuid.NewString()

	// Simpan data ke database
	user := &models.User{
		Name:             data.Name,
		Email:            strings.ToLower(strings.TrimSpace(data.Email)),
		Phone:            helpers.FormatPhoneNumber(data.Phone),
		Password:         hashedPassword,
		UUID:             userUUID,
		DeviceID:         data.DeviceID,
		DeviceName:       data.DeviceName,
		ConfirmationFlow: data.VerificationChannel,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := r.userRepo.Save(user); err != nil {
		return nil, fmt.Errorf("failed to save user: %v", err)
	}

	// Kirim OTP melalui OtpService
	processName := "User Registration"
	token := userUUID
	receiverName := data.Name
	receiver := data.Email
	if data.VerificationChannel == "whatsapp" {
		receiver = data.Phone
	}

	_, err = r.otpService.SendOtp(data.VerificationChannel, processName, token, receiverName, receiver, 6, 300, "en")
	if err != nil {
		return nil, fmt.Errorf("failed to send verification: %v", err)
	}

	return user, nil
}

// BeginTransaction memulai transaksi database
func (s *RegisterService) BeginTransaction() (*sql.Tx, error) {
	if s.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}
	return s.DB.Begin()
}

func (s *RegisterService) GetUserConfirm(emailOrPhone string) (*models.UserConfirm, error) {
	// Simulasikan pencarian data UserConfirm berdasarkan email atau phone
	if emailOrPhone == "" {
		return nil, errors.New("email or phone is required")
	}

	// Contoh mock data UserConfirm
	userConfirm := &models.UserConfirm{
		ID:           "123",
		EmailOrPhone: emailOrPhone,
		Code:         123456,
		DeviceID:     "device123",
		CreatedAt:    time.Now(),
	}
	return userConfirm, nil
}

func (s *RegisterService) VerifyActivationCode(userConfirm *models.UserConfirm, code int) (bool, error) {
	if userConfirm == nil {
		return false, errors.New("user confirm data is required")
	}
	if userConfirm.Code != code {
		return false, nil
	}
	return true, nil
}

func (s *RegisterService) ActivateAccount(userConfirm *models.UserConfirm) (*models.UserActivation, error) {
	if userConfirm == nil {
		return nil, errors.New("user confirm data is required")
	}

	// Simulasikan proses aktivasi akun
	userActivation := &models.UserActivation{
		ID:          "activation123",
		UserID:      "user123",
		Activated:   true,
		ActivatedAt: time.Now(),
	}
	return userActivation, nil
}

func (s *RegisterService) AddDevice(userID, deviceID, device, clientIP string) error {
	if userID == "" || deviceID == "" || device == "" || clientIP == "" {
		return errors.New("all fields are required")
	}

	// Simulasikan penambahan device sukses
	return nil
}

func (s *RegisterService) DeleteUserConfirm(confirmID string) error {
	if confirmID == "" {
		return errors.New("confirm ID is required")
	}

	// Simulasikan penghapusan UserConfirm
	return nil
}

func (s *RegisterService) SendActivationCode(userConfirm *models.UserConfirm, channel string) error {
	if userConfirm == nil {
		return errors.New("user confirm data is required")
	}
	if channel != "whatsapp" && channel != "email" && channel != "sms" {
		return errors.New("invalid channel")
	}

	// Simulasikan pengiriman kode aktivasi sukses
	return nil
}
