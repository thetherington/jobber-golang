package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/thetherington/jobber-auth/internal/adapters/config"
	"github.com/thetherington/jobber-auth/internal/adapters/storage/postgres"
	"github.com/thetherington/jobber-auth/internal/adapters/storage/postgres/repository"
	"github.com/thetherington/jobber-auth/internal/core/port"
	"github.com/thetherington/jobber-auth/internal/util"
	token "github.com/thetherington/jobber-common/client-token"
	"github.com/thetherington/jobber-common/error-handling/svc"
	"github.com/thetherington/jobber-common/middleware"
	"github.com/thetherington/jobber-common/models/auth"
	pb "github.com/thetherington/jobber-common/protogen/go/notification"
	pbUser "github.com/thetherington/jobber-common/protogen/go/users"
	"github.com/thetherington/jobber-common/utils"
	"google.golang.org/protobuf/proto"
)

var validate *validator.Validate

/**
 * AuthService implements
 */
type AuthService struct {
	queries port.AuthRepository
	queue   port.AuthProducer
	token   token.TokenMaker
	image   port.ImageUploader
}

// NewAuthService creates a new auth service instance
func NewAuthService(queries port.AuthRepository, producer port.AuthProducer, token token.TokenMaker, image port.ImageUploader) *AuthService {
	validate = validator.New(validator.WithRequiredStructEnabled())

	return &AuthService{
		queries,
		producer,
		token,
		image,
	}
}

func (as *AuthService) SignUp(ctx context.Context, req *auth.SignUpPayload) (*auth.AuthResponse, error) {
	// Validate signup payload
	if err := req.Validate(validate); err != nil {
		slog.With("error", err).Debug("Signup Validation Failed")
		return nil, svc.NewError(svc.ErrInvalidData, err)
	}

	// Hash password
	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		slog.With("error", err).Error("Hashing Password")

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	// Cloudinary image upload
	profilePublicId := uuid.New().String()
	url, err := as.image.UploadImage(ctx, req.ProfilePicture, profilePublicId, true, true)
	if err != nil {
		slog.With("error", err).Error("Cloudinary Image Upload")

		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	// random characters for the email verify link
	randomCharacters := util.RandomString(20)

	// Save user into the database
	user, err := as.queries.CreateUser(ctx, repository.CreateUserParams{
		Username:               utils.FirstLetterUpperCase(req.Username),
		Password:               hashPassword,
		ProfilePublicID:        profilePublicId,
		Email:                  utils.LowerCase(req.Email),
		Country:                req.Country,
		ProfilePicture:         url,
		EmailVerificationToken: randomCharacters,
	})
	if err != nil {
		if errCode := postgres.ErrorCode(err); errCode == "23505" {
			slog.With("error", err).Debug("DB CreateUser")
			return nil, svc.NewError(svc.ErrInvalidData, fmt.Errorf("username or email already exists"))
		}

		slog.With("error", err).Error("DB CreateUser")
		return nil, svc.NewError(svc.ErrInternalFailure, fmt.Errorf("failed to create user"))
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
			slog.With("error", err).Error("Signup: Failed to send message to jobber-email-notification")
		}
	}

	buyer := &pbUser.BuyerPayload{
		Username:       &user.Username,
		Email:          &user.Email,
		ProfilePicture: &user.ProfilePicture,
		Country:        &user.Country,
		CreatedAt:      utils.ToDateTime(&user.CreatedAt.Time),
		UpdatedAt:      utils.ToDateTime(&user.UpdatedAt.Time),
		Action:         pbUser.Action_AUTH.Enum(),
	}

	// send the user information to the users micro service via rabbitMQ direct exchange using protobuf to create a buyer
	if data, err := proto.Marshal(buyer); err == nil {
		if err := as.queue.PublishDirectMessage("jobber-buyer-update", "user-buyer", data); err != nil {
			slog.With("error", err).Error("Signup: Failed to send message to jobber-buyer-update")
		}
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
		Message: "User created successfully",
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

func (as *AuthService) SignIn(ctx context.Context, req *auth.SignInPayload) (*auth.AuthResponse, error) {
	// Validate signin payload. username must be a string or an email
	if err := req.Validate(validate); err != nil {
		slog.With("error", err).Debug("Sigin Validation Failed")
		return nil, svc.NewError(svc.ErrUnAuthorized, fmt.Errorf("invalid username/password credentials validation"))
	}

	var (
		user    repository.Auth
		respErr error
	)

	// check whether the field is an email to fetch user by email
	if err := validate.Var(req.Username, "email"); err == nil {
		user, respErr = as.queries.GetUserByEmail(ctx, utils.LowerCase(req.Username))
	} else {
		user, respErr = as.queries.GetUserByUsername(ctx, utils.FirstLetterUpperCase(req.Username))
	}

	// user is not found
	if respErr != nil {
		slog.With("error", respErr).Debug("failed to find user")
		return nil, svc.NewError(svc.ErrUnAuthorized, fmt.Errorf("invalid username/password credentials"))
	}

	// validate password matches
	if err := util.CheckPassword(req.Password, user.Password); err != nil {
		return nil, svc.NewError(svc.ErrUnAuthorized, fmt.Errorf("invalid username/password credentials"))
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
		Message: "Login Successfully!",
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

func (as *AuthService) VerifyEmail(ctx context.Context, req *auth.VerifyEmail) (*auth.AuthResponse, error) {
	// Validate VerifyEmail payload
	if err := req.Validate(validate); err != nil {
		slog.With("error", err).Debug("VerifyEmail token validation")
		return nil, svc.NewError(svc.ErrInvalidData, fmt.Errorf("verification token is either invalid or is already used"))
	}

	// get a user document via the verification token. if null then throw bad request
	user, err := as.queries.GetUserByVerificationToken(ctx, req.Token)
	if err != nil {
		slog.With("error", err).Debug("GetUserByVerificationToken failed")
		return nil, svc.NewError(svc.ErrBadRequest, fmt.Errorf("verification token is either invalid or is already used"))
	}

	// update email verification and set to true
	err = as.queries.UpdateVerifyEmailField(ctx, repository.UpdateVerifyEmailFieldParams{
		ID:            user.ID,
		EmailVerified: true,
	})
	if err != nil {
		slog.With("error", err).Debug("UpdateVerifyEmailField failed")
		return nil, svc.NewError(svc.ErrInternalFailure, fmt.Errorf("failed to verify email"))
	}

	// create response structure and return it to the gRPC handler.
	resp := &auth.AuthResponse{
		Message: "Email verified successfully.",
		User: &auth.AuthDocument{
			Id:              int32(user.ID),
			ProfilePublicId: user.ProfilePublicID,
			Username:        user.Username,
			Email:           user.Email,
			Country:         user.Country,
			EmailVerified:   true,
			ProfilePicture:  user.ProfilePicture,
			CreatedAt:       &user.CreatedAt.Time,
			UpdatedAt:       utils.ToTimePtr(time.Now()),
		},
	}

	return resp, nil
}

func (as *AuthService) ForgotPassword(ctx context.Context, req *auth.ForgotPassword) (*auth.AuthResponse, error) {
	// Validate ForgotPassword payload
	if err := req.Validate(validate); err != nil {
		slog.With("error", err).Debug("ForgotPassword email validation")
		return nil, svc.NewError(svc.ErrInvalidData, fmt.Errorf("invalid email address"))
	}

	user, err := as.queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		slog.With("error", err).Debug("ForgotPassword GetUserByEmail failed")
		return nil, svc.NewError(svc.ErrBadRequest, fmt.Errorf("invalid credentials"))
	}

	// generate random characters that will be used for the email verification
	randomCharacters := util.RandomString(20)

	// create a new date and add 1 hour to it
	expiry := time.Now().Add(1 * time.Hour)

	// update the database
	err = as.queries.UpdatePasswordToken(ctx, repository.UpdatePasswordTokenParams{
		ID:                   user.ID,
		PasswordResetToken:   pgtype.Text{String: randomCharacters, Valid: true},
		PasswordResetExpires: pgtype.Timestamptz{Time: expiry, Valid: true},
	})
	if err != nil {
		slog.With("error", err).Debug("ForgotPassword UpdatePasswordToken failed")
		return nil, svc.NewError(svc.ErrInternalFailure, fmt.Errorf("password reset failed"))
	}

	// reset link url
	resetLink := fmt.Sprintf("%s/reset_password?token=%s", config.Config.App.ClientUrl, randomCharacters)

	// notification message for rabbitmq
	msg := &pb.EmailMessageDetails{
		ReceiverEmail: &user.Email,
		ResetLink:     &resetLink,
		Username:      &user.Username,
		Template:      utils.Ptr("forgotPassword"),
	}

	// send the user information to the notification micro service via rabbitMQ direct exchange using protobuf
	if data, err := proto.Marshal(msg); err == nil {
		if err := as.queue.PublishDirectMessage("jobber-email-notification", "auth-email", data); err != nil {
			slog.With("error", err).Error("ForgotPassword: Failed to send message to jobber-email-notification")
		}
	} else {
		slog.With("error", err).Error("ForgotPassword: Failed to marshal message to notification service")
	}

	// create response structure and return it to the gRPC handler.
	resp := &auth.AuthResponse{
		Message: "Password reset email sent.",
	}

	return resp, nil
}

func (as *AuthService) ResetPassword(ctx context.Context, req *auth.ResetPassword) (*auth.AuthResponse, error) {
	// Validate ResetPassword payload
	if err := req.Validate(validate); err != nil {
		slog.With("error", err).Debug("ResetPassword validation")
		return nil, svc.NewError(svc.ErrInvalidData, err)
	}

	if req.Password != req.ConfirmPassword {
		slog.Debug("ResetPassword passwords don't match")
		return nil, svc.NewError(svc.ErrInvalidData, fmt.Errorf("passwords don't match"))
	}

	user, err := as.queries.GetUserByPasswordToken(ctx, pgtype.Text{
		String: req.Token,
		Valid:  true,
	})
	if err != nil {
		slog.With("error", err).Debug("ResetPassword GetUserByPasswordToken failed")
		return nil, svc.NewError(svc.ErrBadRequest, fmt.Errorf("invalid or expired password token"))
	}

	// Hash password
	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		slog.With("error", err).Error("Hashing Password")
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	// Update user password in the database
	err = as.queries.UpdateUserPassword(ctx, repository.UpdateUserPasswordParams{
		ID:       user.ID,
		Password: hashPassword,
	})
	if err != nil {
		slog.With("error", err).Error("ResetPassword UpdateUserPassword failed")
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	msg := &pb.EmailMessageDetails{
		ReceiverEmail: &user.Email,
		Username:      utils.Ptr(user.Username),
		Template:      utils.Ptr("resetPasswordSuccess"),
	}

	// send the user information to the notification micro service via rabbitMQ direct exchange using protobuf
	if data, err := proto.Marshal(msg); err == nil {
		if err := as.queue.PublishDirectMessage("jobber-email-notification", "auth-email", data); err != nil {
			slog.With("error", err).Error("Signup: Failed to send message to jobber-email-notification")
		}
	} else {
		slog.With("error", err).Error("Signup: Failed to marshal message to notification service")
	}

	// create response structure and return it to the gRPC handler.
	resp := &auth.AuthResponse{
		Message: "Password successfully updated.",
	}

	return resp, nil
}

func (as *AuthService) ChangePassword(ctx context.Context, req *auth.ChangePassword) (*auth.AuthResponse, error) {
	// Validate ChangePassword payload
	if err := req.Validate(validate); err != nil {
		slog.With("error", err).Debug("ResetPassword validation")
		return nil, svc.NewError(svc.ErrInvalidData, err)
	}

	// validate that the new password is different
	if req.CurrentPassword == req.NewPassword {
		slog.Debug("ChangePassword passwords are the same")
		return nil, svc.NewError(svc.ErrBadRequest, fmt.Errorf("new password is the same"))
	}

	// Get username from the user cookie session passed down into the context.
	username := ctx.Value(middleware.CtxUsernameKey)
	if username == nil {
		slog.Debug("Username in context is nil")
		return nil, svc.NewError(svc.ErrUnAuthorized, fmt.Errorf("session invalid. please login and try again"))
	}

	// check that the username exists. if there's an error here something is very out of sync.
	user, err := as.queries.GetUserByUsername(ctx, username.(string))
	if err != nil {
		slog.With("error", err).Debug("ChangePassword GetUserByUsername failed")
		return nil, svc.NewError(svc.ErrUnAuthorized, fmt.Errorf("failed to change password. please login and try again"))
	}

	// Hash password
	hashPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		slog.With("error", err).Error("Hashing Password")
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	// Update user password in the database
	err = as.queries.UpdateUserPassword(ctx, repository.UpdateUserPasswordParams{
		ID:       user.ID,
		Password: hashPassword,
	})
	if err != nil {
		slog.With("error", err).Error("ChangePassword UpdateUserPassword failed")
		return nil, svc.NewError(svc.ErrInternalFailure, err)
	}

	msg := &pb.EmailMessageDetails{
		ReceiverEmail: &user.Email,
		Username:      utils.Ptr(user.Username),
		Template:      utils.Ptr("resetPasswordSuccess"),
	}

	// send the user information to the notification micro service via rabbitMQ direct exchange using protobuf
	if data, err := proto.Marshal(msg); err == nil {
		if err := as.queue.PublishDirectMessage("jobber-email-notification", "auth-email", data); err != nil {
			slog.With("error", err).Error("Signup: Failed to send message to jobber-email-notification")
		}
	} else {
		slog.With("error", err).Error("Signup: Failed to marshal message to notification service")
	}

	// create response structure and return it to the gRPC handler.
	resp := &auth.AuthResponse{
		Message: "Password successfully updated.",
	}

	return resp, nil
}

func (as *AuthService) Seed(ctx context.Context, count int) (string, error) {
	if count == 0 {
		return "", svc.NewError(svc.ErrInvalidData, fmt.Errorf("count must be greater than 0"))
	}

	type User struct {
		Username      string `fake:"{firstname}"`
		Password      string `fake:"-"`
		Email         string `fake:"{email}"`
		Country       string `fake:"{country}"`
		PofilePicture string `fake:"-"`
	}

	for i := 0; i < count; i++ {
		var (
			u   User
			err error
		)

		err = gofakeit.Struct(&u)
		if err != nil {
			slog.With("error", err).Error("gofakeit")
		}

		// Hash password
		u.Password, err = util.HashPassword("secret")
		if err != nil {
			slog.With("error", err).Error("Hashing Password")
			return "", svc.NewError(svc.ErrInternalFailure, err)
		}

		// random profile picture link
		u.PofilePicture = fmt.Sprintf("https://picsum.photos/seed/%s/640/480", util.RandomString(7))

		// Save user into the database
		user, err := as.queries.CreateUser(ctx, repository.CreateUserParams{
			Username:               utils.FirstLetterUpperCase(u.Username),
			Password:               u.Password,
			ProfilePublicID:        uuid.New().String(),
			Email:                  utils.LowerCase(u.Email),
			Country:                u.Country,
			ProfilePicture:         u.PofilePicture,
			EmailVerificationToken: util.RandomString(20),
		})
		if err != nil {
			if errCode := postgres.ErrorCode(err); errCode == "23505" {
				slog.With("error", err).Debug("DB CreateUser")
				return "", svc.NewError(svc.ErrInvalidData, fmt.Errorf("username or email already exists"))
			}

			slog.With("error", err).Error("DB CreateUser")
			return "", svc.NewError(svc.ErrInternalFailure, fmt.Errorf("failed to create user"))
		}

		buyer := &pbUser.BuyerPayload{
			Username:       &user.Username,
			Email:          &user.Email,
			ProfilePicture: &user.ProfilePicture,
			Country:        &user.Country,
			CreatedAt:      utils.ToDateTime(&user.CreatedAt.Time),
			UpdatedAt:      utils.ToDateTime(&user.UpdatedAt.Time),
			Action:         pbUser.Action_AUTH.Enum(),
		}

		// send the user information to the users micro service via rabbitMQ direct exchange using protobuf to create a buyer
		if data, err := proto.Marshal(buyer); err == nil {
			if err := as.queue.PublishDirectMessage("jobber-buyer-update", "user-buyer", data); err != nil {
				slog.With("error", err).Error("Seed: Failed to send message to jobber-buyer-update")
			}
		}
	}

	return fmt.Sprintf("Seeded %d Users Successfully", count), nil
}

// fmt.Println(ctx.Value(middleware.CtxUsernameKey))
// fmt.Println(ctx.Value(middleware.CtxEmailKey))
