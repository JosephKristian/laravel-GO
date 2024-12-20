package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/JosephKristian/project-migration/internal/models"
	notification "github.com/JosephKristian/project-migration/internal/notifications" // Untuk mengirim SMS atau email
	"github.com/JosephKristian/project-migration/internal/repositories"
)

type OtpService struct {
	otpRepository repositories.OtpRepository
}

func NewOtpService(otpRepository repositories.OtpRepository) *OtpService {
	return &OtpService{otpRepository: otpRepository}
}

// Generate OTP
func (s *OtpService) GenerateOtp(length int) string {
	rand.Seed(time.Now().UnixNano()) // Inisialisasi seed menggunakan waktu saat ini
	otp := make([]byte, length)
	for i := range otp {
		otp[i] = byte('0' + rand.Intn(10)) // Konversi angka ke karakter '0'-'9'
	}
	return string(otp) // Konversi array byte menjadi string
}

// Send OTP
func (s *OtpService) SendOtp(via, processName, token, receiverName, receiver string, otpLength, expiredInSeconds int, lang string) (string, error) {
	otp := s.GenerateOtp(otpLength)
	expiredAt := time.Now().Add(time.Duration(expiredInSeconds) * time.Second) // Gunakan time.Time langsung

	// Simpan OTP ke repository (misalnya database atau cache)
	otpData := models.Otp{
		Uuid:        fmt.Sprintf("%s-%s", token, receiver), // UUID bisa lebih baik menggunakan UUID generator
		Otp:         otp,
		Token:       token,
		Destination: receiver,
		Flow:        processName,
		Channel:     via,
		ExpiredAt:   expiredAt, // Set sebagai time.Time
	}
	err := s.otpRepository.StoreOtp(&otpData)
	if err != nil {
		return "", fmt.Errorf("failed to store OTP: %v", err)
	}

	// Kirim OTP berdasarkan saluran
	var message string
	switch via {
	case "email":
		// message = fmt.Sprintf("OTP untuk %s: %s. OTP berlaku selama %d menit.", processName, otp, expiredInSeconds/60)
		err = notification.SendEmail(receiver, receiverName, processName, otp, lang)
	case "sms":
		message = fmt.Sprintf("OTP untuk %s: %s. OTP berlaku selama %d menit.", processName, otp, expiredInSeconds/60)
		err = notification.SendSms(receiver, message)
	case "whatsapp":
		err = notification.SendWhatsapp(receiver, processName, otp, expiredInSeconds/60, lang)
	default:
		return "", errors.New("unsupported channel")
	}

	if err != nil {
		return "", fmt.Errorf("failed to send OTP: %v", err)
	}

	// Simpan OTP ke cache untuk validasi nanti
	// Simulasi penyimpanan ke cache (Redis)
	// Cache.Set("otp-"+token, otp, time.Duration(expiredInSeconds)*time.Second)

	return otpData.Uuid, nil
}

// Claim OTP
func (s *OtpService) ClaimOtp(token, otpCode string) (bool, error) {
	// Ambil OTP dari database
	otpData, err := s.otpRepository.GetOtpByToken(token)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve OTP: %v", err)
	}

	// Periksa apakah OTP sudah expired
	if time.Now().After(otpData.ExpiredAt) {
		return false, errors.New("OTP expired")
	}

	// Periksa apakah OTP yang dimasukkan sesuai
	if otpData.Otp != otpCode {
		return false, errors.New("invalid OTP")
	}

	// Tandai OTP sebagai terverifikasi
	err = s.otpRepository.UpdateOtpVerificationStatus(otpData.Uuid)
	if err != nil {
		return false, fmt.Errorf("failed to verify OTP: %v", err)
	}

	return true, nil
}
