package models

type EmailVerification struct {
	Otp   string `json:"otp"`
	Email string `json:"email"`
}
