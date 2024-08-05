package dto

type RegisterReq struct {
	UserId   string `json:"user_id" validate:"omitempty"`
	Email    string `json:"email" validate:"required,email,min=5,max=100"`
	FullName string `json:"full_name" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=5,max=100"`
}

type LoginWithGoogleReq struct {
	UserId       string `json:"user_id" validate:"required,min=21,max=21"`
	Email        string `json:"email" validate:"required,email,min=5,max=100"`
	FullName     string `json:"name" validate:"required,min=3,max=100"`
	PhotoProfile string `json:"picture" validate:"required,min=3,max=500"`
	RefreshToken string `json:"refresh_token" validate:"required,min=50,max=1000"`
}

type LoginReq struct {
	Email    string `json:"email" validate:"required,email,min=5,max=100"`
	Password string `json:"password" validate:"required,min=5,max=100"`
}

type SendOtpReq struct {
	Email string `json:"email" validate:"required,email,min=5,max=100"`
	Otp   string `json:"otp" validate:"required,max=6"`
}

type VerifyOtpReq struct {
	Email string `json:"email" validate:"required,email,min=5,max=100"`
	Otp   string `json:"otp" validate:"required,max=6"`
}