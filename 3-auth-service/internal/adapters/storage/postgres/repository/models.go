// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package repository

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Auth struct {
	ID                     int64
	Username               string
	Password               string
	ProfilePublicID        string
	Email                  string
	Country                string
	ProfilePicture         string
	EmailVerificationToken string
	EmailVerified          bool
	CreatedAt              pgtype.Timestamptz
	UpdatedAt              pgtype.Timestamptz
	PasswordResetToken     pgtype.Text
	PasswordResetExpires   pgtype.Timestamptz
}
