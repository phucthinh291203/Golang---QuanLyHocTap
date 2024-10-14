package Models

import "time"

type OTP struct {
	Email     string    `bson:"email"`
	Role      string    `bson:"role"`
	Code      string    `bson:"code"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
	Verified  bool      `bson:"verified"`
}

type SMS struct {
	PhoneNumber string    `bson:"phone_number"`
	Role        string    `bson:"role"`
	Code        string    `bson:"code"`
	CreatedAt   time.Time `bson:"created_at"`
	ExpiresAt   time.Time `bson:"expires_at"`
	Verified    bool      `bson:"verified"`
}

type RequestSendSMS struct {
	PhoneNumber string `bson:"phone_number"`
	Role        string `bson:"string"`
}

type RequestVerifySMS struct {
	PhoneNumber string `bson:"phone_number"`
	OTPInput string `bson:"otp_input"`
}

type RequestResetPasswordSMS struct {
	PhoneNumber        string `bson:"phone_number"`
	NewPassword        string `bson:"new_password"`
	ConfirmNewPassword string `bson:"confirm_password"`
}

type RequestSendOTP struct {
	Email string `bson:"email"`
	Role  string `bson:"string"`
}

type RequestVerifyOTP struct {
	Email    string `bson:"email"`
	OTPInput string `bson:"otp_input"`
}

type RequestResetPassword struct {
	Email              string `bson:"email"`
	NewPassword        string `bson:"new_password"`
	ConfirmNewPassword string `bson:"confirm_password"`
}
