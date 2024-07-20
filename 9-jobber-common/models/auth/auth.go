package auth

import (
	"time"

	"github.com/google/uuid"
)

type AuthPayload struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Iat      int    `json:"iat,omitempty"`
}

type SignUpPayload struct {
	Username       string `json:"username"       validate:"required,min=4,max=12"  errmsg:"Invalid name min(4) max(12)"`
	Password       string `json:"password"       validate:"required,min=4,max=12"  errmsg:"Bad password min(4) max(12)"`
	Email          string `json:"email"          validate:"required,email"         errmsg:"Invalid address"`
	Country        string `json:"country"        validate:"required,alpha"         errmsg:"Invalid country name"`
	ProfilePicture string `json:"profilePicture" validate:"required"               errmsg:"Invalid picture"`
}

type SignInPayload struct {
	Username string `json:"username"  validate:"required,email|required,min=4"  errmsg:"Invalid username min(4) or email"`
	Password string `json:"password"  validate:"required,min=4,max=12"          errmsg:"Bad password min(4) max(12)"`
}

type ForgotPassword struct {
	Email string `json:"email"  validate:"required,email"  errmsg:"Invalid"`
}

type ResetPassword struct {
	Password        string `json:"password"         validate:"required,min=4,max=12"  errmsg:"Bad password min(4) max(12)"`
	ConfirmPassword string `json:"confirmPassword"  validate:"required,min=4,max=12"  errmsg:"Bad password min(4) max(12)"`
	Token           string `json:"token"            validate:"required,min=20"        errmsg:"Invalid password change token"`
}

type ChangePassword struct {
	CurrentPassword string `json:"currentPassword"  validate:"required,min=4,max=12"  errmsg:"Bad password min(4) max(12)"`
	NewPassword     string `json:"newPassword"      validate:"required,min=4,max=12"  errmsg:"Bad password min(4) max(12)"`
}

type VerifyEmail struct {
	Token string `json:"token" validate:"required,min=20" errmsg:"Invalid token min(20)"`
}

type RefreshToken struct {
	Username string `json:"username" validate:"required,min=4" errmsg:"Invalid username"`
}

type ResendEmail struct {
	Id    int32  `json:"userId" validate:"required,number"  errmsg:"Invalid"`
	Email string `json:"email"  validate:"required,email"   errmsg:"Invalid"`
}

type AuthDocument struct {
	Id                     int32      `json:"id,omitempty"`
	ProfilePublicId        string     `json:"profilePublicId,omitempty"`
	Username               string     `json:"username,omitempty"`
	Email                  string     `json:"email,omitempty"`
	Password               string     `json:"password,omitempty"`
	Country                string     `json:"country,omitempty"`
	ProfilePicture         string     `json:"profilePicture,omitempty"`
	EmailVerified          bool       `json:"emailVerified"`
	EmailVerificationToken string     `json:"emailVerificationToken,omitempty"`
	CreatedAt              *time.Time `json:"createdAt,omitempty"`
	UpdatedAt              *time.Time `json:"UpdatedAt,omitempty"`
	PasswordResetToken     string     `json:"passwordResetToken,omitempty"`
	PasswordResetExpires   *time.Time `json:"passwordResetExpires,omitempty"`
}

type AuthResponse struct {
	Message string        `json:"message"`
	User    *AuthDocument `json:"user,omitempty"`
	Token   string        `json:"token,omitempty"`
}

type TokenPayload struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}
