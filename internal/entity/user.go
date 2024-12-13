package entity

// User defines the data required for user registration
type User struct {
	Name             string `json:"name" validate:"required"`
	Email            string `json:"email" validate:"omitempty,email"`
	Phone            string `json:"phone" validate:"required,numeric,len=10"`
	Password         string `json:"password" validate:"omitempty,min=8"`
	ReferralCode     string `json:"referral_code" validate:"omitempty"`
	Website          string `json:"website" validate:"omitempty,url"`
	ConfirmationFlow string `json:"confirmation_flow" validate:"omitempty,oneof=phone email"`
	DeviceID         string `json:"device_id"`
	DeviceName       string `json:"device_name"`
	IP               string `json:"ip"`
}
