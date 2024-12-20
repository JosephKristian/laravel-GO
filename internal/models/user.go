package models

import "time"

type User struct {
	Name             string    `json:"name" validate:"required"`
	Email            string    `json:"email" validate:"omitempty,email"`
	Phone            string    `json:"phone" validate:"required,numeric,len=10"`
	Password         string    `json:"password" validate:"omitempty,min=8"`
	ReferralCode     string    `json:"referral_code" validate:"omitempty"`
	Website          string    `json:"website" validate:"omitempty,url"`
	ConfirmationFlow string    `json:"confirmation_flow" validate:"omitempty,oneof=phone email"`
	DeviceID         string    `json:"device_id"`
	DeviceName       string    `json:"device_name"`
	IP               string    `json:"ip,omitempty"`
	UUID             string    `json:"uuid"`              // Added UUID field
	VerificationCode string    `json:"verification_code"` // Added VerificationCode field
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// UserConfirm represents the data for user confirmation
type UserConfirm struct {
	ID           string    `json:"id"`
	EmailOrPhone string    `json:"email_or_phone"`
	Code         int       `json:"code"`
	DeviceID     string    `json:"device_id"`
	Device       string    `json:"device"`
	CreatedAt    time.Time `json:"created_at"`
}

// UserActivation represents the data for user account activation
type UserActivation struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Activated   bool      `json:"activated"`
	ActivatedAt time.Time `json:"activated_at"`
}
