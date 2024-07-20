package svc

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// 400 - http.StatusBadRequest / codes.FailedPrecondition
	ErrBadRequest = errors.New("bad request")
	// 500 - http.StatusInternalServerError / codes.Internal
	ErrInternalFailure = errors.New("internal failure")
	// 404 - http.StatusNotFound / codes.NotFound
	ErrNotFound = errors.New("not found")
	// 400 - http.StatusBadRequest / codes.InvalidArgument
	ErrInvalidData = errors.New("bad data")
	// 401 - http.StatusUnauthorized
	ErrUnAuthorized = errors.New("unauthorized")
	// 500 - http.StatusInternalServerError / codes.Unavailable
	ErrUnavailable = errors.New("internal service unreachable")

	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Error struct {
	appErr error
	svcErr error
}

func (e Error) AppError() error {
	return e.appErr
}

func (e Error) SvcError() error {
	return e.svcErr
}

func NewError(svcErr, appErr error) error {
	return Error{
		svcErr: svcErr,
		appErr: appErr,
	}
}

func (e Error) Error() string {
	return errors.Join(e.svcErr, e.appErr).Error()
}

func GrpcErrorResolve(err error, caller string) error {
	st, ok := status.FromError(err)
	if ok {
		switch st.Code() {
		// 400 - http.StatusBadRequest
		case codes.FailedPrecondition:
			return NewError(ErrBadRequest, fmt.Errorf(st.Message()))

		// 500 - http.StatusInternalServerError
		case codes.Internal:
			return NewError(ErrInternalFailure, fmt.Errorf(st.Message()))

		// 404 - http.StatusNotFound
		case codes.NotFound:
			return NewError(ErrNotFound, fmt.Errorf(st.Message()))

		// 400 - http.StatusBadRequest
		case codes.InvalidArgument:
			return NewError(ErrBadRequest, fmt.Errorf(st.Message()))

		// 404 - http
		case codes.Unauthenticated:
			return NewError(ErrUnAuthorized, fmt.Errorf(st.Message()))

		// 500 - http.StatusInternalServerError
		case codes.Unavailable:
			return NewError(ErrInternalFailure, fmt.Errorf("%s: %w", caller, ErrUnavailable))
		}
	}

	// send error back because unmatched.
	return err
}
