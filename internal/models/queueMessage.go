package models

type NotificationType string
type MessageType string

const (
	Email NotificationType = "email"
	Sms                    = "sms"
)

const (
	RegEmailVerification   MessageType = "reg_email_verification"
	RegPhoneVerification   MessageType = "reg_phone_verification"
	LoginEmailVerification MessageType = "login_email_verification"
	LoginPhoneVerification MessageType = "login_phone_verification"
	Welcome                MessageType = "welcome"
)

type QueueMessage struct {
	Otp              string           `json:"otp"`
	User             User             `json:"user"`
	NotificationType NotificationType `json:"notificationType"`
	Message          string           `json:"message"`
	MessageType      MessageType      `json:"messageType"`
}
