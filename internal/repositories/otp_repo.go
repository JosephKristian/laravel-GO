package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/JosephKristian/project-migration/internal/models"
)

type OtpRepository interface {
	StoreOtp(*models.Otp) error
	GetOtpByToken(string) (*models.Otp, error)
	UpdateOtpVerificationStatus(string) error
}

type OtpRepo struct {
	DB *sql.DB
}

func (r *OtpRepo) StoreOtp(otp *models.Otp) error {
	// Menyusun query SQL untuk menyimpan data OTP
	query := `INSERT INTO otps (uuid, otp, token, destination, flow, channel, expired_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.DB.Exec(query, otp.Uuid, otp.Otp, otp.Token, otp.Destination, otp.Flow, otp.Channel, otp.ExpiredAt)
	if err != nil {
		return fmt.Errorf("failed to store OTP: %v", err)
	}
	return nil
}

func (r *OtpRepo) GetOtpByToken(token string) (*models.Otp, error) {
	// Menyusun query SQL untuk mengambil OTP berdasarkan token
	query := `SELECT uuid, otp, token, destination, flow, channel, expired_at 
			  FROM otps WHERE token = $1`

	row := r.DB.QueryRow(query, token)

	var otp models.Otp
	if err := row.Scan(&otp.Uuid, &otp.Otp, &otp.Token, &otp.Destination, &otp.Flow, &otp.Channel, &otp.ExpiredAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("OTP not found")
		}
		return nil, fmt.Errorf("failed to retrieve OTP: %v", err)
	}

	return &otp, nil
}

func (r *OtpRepo) UpdateOtpVerificationStatus(uuid string) error {
	// Menyusun query SQL untuk memperbarui status verifikasi OTP
	query := `UPDATE otps SET verified_at = NOW() WHERE uuid = ?`

	_, err := r.DB.Exec(query, uuid)
	if err != nil {
		return fmt.Errorf("failed to update OTP verification status: %v", err)
	}
	return nil
}
