package model

type NotificationType string
type MessageType string

const (
	Email NotificationType = "email"
	Sms                    = "sms"
)

const (
	RegEmailVerification   MessageType = "reg_email_verification"
	RegPhoneVerification               = "reg_phone_verification"
	LoginEmailVerification             = "login_email_verification"
	LoginPhoneVerification             = "login_phone_verification"
	Welcome                            = "welcome"
)

type Notification struct {
	Otp              string           `json:"otp"`
	User             User             `json:"user"`
	NotificationType NotificationType `json:"notification_type"`
	Message          string           `json:"message"`
	MessageType      MessageType      `json:"message_type"`
}
