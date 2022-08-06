package models

type OtpVerification struct {
	Otp       string `json:"otp"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	PhoneCode string `json:"phoneCode"`
}
