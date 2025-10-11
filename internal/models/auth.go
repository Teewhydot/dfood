package models

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type LogoutRequest struct {
	Token string `json:"token" binding:"required"`
}

type DeleteAccountRequest struct {
	Email string `json:"email" binding:"required,email"`
	Token string `json:"token" binding:"required"`
}

type VerifyOtpRequest struct {
	Email string `json:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type EmailVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}
