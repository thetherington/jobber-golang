package auth

import (
	"github.com/go-playground/validator/v10"
	"github.com/thetherington/jobber-common/utils"
)

func (m *SignUpPayload) Validate(validate *validator.Validate) error {
	return utils.ValidateFunc[SignUpPayload](*m, validate)
}

func (m *SignInPayload) Validate(validate *validator.Validate) error {
	return utils.ValidateFunc[SignInPayload](*m, validate)
}

func (m *ForgotPassword) Validate(validate *validator.Validate) error {
	return utils.ValidateFunc[ForgotPassword](*m, validate)
}

func (m *VerifyEmail) Validate(validate *validator.Validate) error {
	return utils.ValidateFunc[VerifyEmail](*m, validate)
}

func (m *ResetPassword) Validate(validate *validator.Validate) error {
	return utils.ValidateFunc[ResetPassword](*m, validate)
}

func (m *ChangePassword) Validate(validate *validator.Validate) error {
	return utils.ValidateFunc[ChangePassword](*m, validate)
}

func (m *RefreshToken) Validate(validate *validator.Validate) error {
	return utils.ValidateFunc[RefreshToken](*m, validate)
}

func (m *ResendEmail) Validate(validate *validator.Validate) error {
	return utils.ValidateFunc[ResendEmail](*m, validate)
}
