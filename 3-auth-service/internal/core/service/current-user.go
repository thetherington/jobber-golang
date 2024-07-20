package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/thetherington/jobber-auth/internal/adapters/config"
	"github.com/thetherington/jobber-auth/internal/adapters/storage/postgres/repository"
	"github.com/thetherington/jobber-auth/internal/util"
	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/middleware"
	"github.com/thetherington/jobber-common/models/auth"
	pb "github.com/thetherington/jobber-common/protogen/go/notification"
	"github.com/thetherington/jobber-common/utils"
	"google.golang.org/protobuf/proto"
)

func (as *AuthService) CurrentUser(ctx context.Context) (*auth.AuthResponse, error) {
	// Get username from the user cookie session passed down into the context.
	username := ctx.Value(middleware.CtxUsernameKey)
	if username == nil {
		slog.Debug("Username in context is nil")
		return nil, svc.NewError(svc.ErrUnAuthorized, fmt.Errorf("session invalid. please login and try again"))
	}

	// check that the username exists. if there's an error here something is very out of sync.
	user, err := as.queries.GetUserByUsername(ctx, username.(string))
	if err != nil {
		slog.With("error", err).Debug("CurrentUser GetUserByUsername failed")
		return nil, svc.NewError(svc.ErrUnAuthorized, fmt.Errorf("failed to get current user. please login and try again"))
	}

	// create response structure and return it to the gRPC handler.
	resp := &auth.AuthResponse{
		Message: "Authenticated user",
		User: &auth.AuthDocument{
			Id:              int32(user.ID),
			ProfilePublicId: user.ProfilePublicID,
			Username:        user.Username,
			Email:           user.Email,
			Country:         user.Country,
			EmailVerified:   user.EmailVerified,
			ProfilePicture:  user.ProfilePicture,
			CreatedAt:       &user.CreatedAt.Time,
			UpdatedAt:       &user.UpdatedAt.Time,
		},
	}

	return resp, nil
}

func (as *AuthService) RefreshToken(ctx context.Context, req *auth.RefreshToken) (*auth.AuthResponse, error) {
	// Validate RefreshToken payload
	if err := req.Validate(validate); err != nil {
		slog.With("error", err).Debug("RefreshToken Validation Failed")
		return nil, svc.NewError(svc.ErrInvalidData, err)
	}

	// check that the username exists. if there's an error here something is very out of sync.
	user, err := as.queries.GetUserByUsername(ctx, req.Username)
	if err != nil {
		slog.With("error", err).Debug("RefreshToken GetUserByUsername failed")
		return nil, svc.NewError(svc.ErrUnAuthorized, fmt.Errorf("failed to get current user. please login and try again"))
	}

	// Create JWT Token
	// TODO token duration via environment variable
	jwt, _, err := as.token.CreateToken(user.Username, user.Email, 24*time.Hour)
	if err != nil {
		slog.With("error", err).Error("Signup: Failed to create token")
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	// create response structure and return it to the gRPC handler.
	resp := &auth.AuthResponse{
		Message: "Refresh token",
		User: &auth.AuthDocument{
			Id:              int32(user.ID),
			ProfilePublicId: user.ProfilePublicID,
			Username:        user.Username,
			Email:           user.Email,
			Country:         user.Country,
			EmailVerified:   user.EmailVerified,
			ProfilePicture:  user.ProfilePicture,
			CreatedAt:       &user.CreatedAt.Time,
			UpdatedAt:       &user.UpdatedAt.Time,
		},
		Token: jwt,
	}

	return resp, nil
}

func (as *AuthService) ResendEmail(ctx context.Context, req *auth.ResendEmail) (*auth.AuthResponse, error) {
	// Validate ResendEmail payload
	if err := req.Validate(validate); err != nil {
		slog.With("error", err).Debug("ResendEmail Validation Failed")
		return nil, svc.NewError(svc.ErrInvalidData, err)
	}

	user, err := as.queries.GetUserByEmail(ctx, utils.LowerCase(req.Email))
	if err != nil {
		slog.With("error", err).Debug("failed to find user")
		return nil, svc.NewError(svc.ErrBadRequest, fmt.Errorf("invalid email"))
	}

	// random characters for the email verify link
	randomCharacters := util.RandomString(20)

	rows, err := as.queries.UpdateEmailVerificationToken(ctx, repository.UpdateEmailVerificationTokenParams{
		ID:                     int64(req.Id),
		EmailVerificationToken: randomCharacters,
	})
	if err != nil || rows < 1 {
		slog.With("error", err).Debug("ResendEmail UpdateEmailVerificationToken failed")
		return nil, svc.NewError(svc.ErrInternalFailure, fmt.Errorf("updating email verification token failed"))
	}

	// Send user info to email notification via protobuf message
	verificationLink := fmt.Sprintf("%s/confirm_email?v_token=%s", config.Config.App.ClientUrl, randomCharacters)

	msg := &pb.EmailMessageDetails{
		ReceiverEmail: &user.Email,
		Username:      utils.Ptr(user.Username),
		VerifyLink:    utils.Ptr(verificationLink),
		Template:      utils.Ptr("verifyEmail"),
	}

	// send the user information to the notification micro service via rabbitMQ direct exchange using protobuf
	if data, err := proto.Marshal(msg); err == nil {
		if err := as.queue.PublishDirectMessage("jobber-email-notification", "auth-email", data); err != nil {
			slog.With("error", err).Error("ResendEmail: Failed to send message to jobber-email-notification")
		}
	} else {
		slog.With("error", err).Error("ResendEmail: Failed to marshal message to notification service")
	}

	// create response structure and return it to the gRPC handler.
	resp := &auth.AuthResponse{
		Message: "Email verification sent",
		User: &auth.AuthDocument{
			Id:              int32(user.ID),
			ProfilePublicId: user.ProfilePublicID,
			Username:        user.Username,
			Email:           user.Email,
			Country:         user.Country,
			EmailVerified:   user.EmailVerified,
			ProfilePicture:  user.ProfilePicture,
			CreatedAt:       &user.CreatedAt.Time,
			UpdatedAt:       &user.UpdatedAt.Time,
		},
	}

	return resp, nil
}
