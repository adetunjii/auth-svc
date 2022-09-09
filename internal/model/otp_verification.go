package model

type OtpVerification struct {
	Otp       string `json:"otp"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	PhoneCode string `json:"phoneCode"`
}

type OtpType string

const (
	REG            OtpType = "REG"
	LOGIN          OtpType = "LOGIN"
	RESET_PASSWORD OtpType = "RESET_PASSWORD"
)
