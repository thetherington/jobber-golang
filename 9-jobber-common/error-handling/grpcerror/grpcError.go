package grpcerror

import (
	"errors"

	"github.com/thetherington/jobber-common/error-handling/svc"
	"google.golang.org/grpc/codes"
)

type APIError struct {
	Status  codes.Code
	Message string
}

func FromError(err error) (APIError, bool) {
	var apiError APIError
	var svcError svc.Error

	if errors.As(err, &svcError) {
		apiError.Message = svcError.AppError().Error()
		svcError := svcError.SvcError()

		switch svcError {
		// 400 - http.StatusBadRequest
		case svc.ErrBadRequest:
			apiError.Status = codes.FailedPrecondition

		// 500 - http.StatusInternalServerError
		case svc.ErrInternalFailure:
			apiError.Status = codes.Internal

		// 401 - http.StatusNotFound
		case svc.ErrNotFound:
			apiError.Status = codes.NotFound

		// 400 - http.StatusBadRequest
		case svc.ErrInvalidData:
			apiError.Status = codes.InvalidArgument

		// 404 - http.StatusUnauthorized
		case svc.ErrUnAuthorized:
			apiError.Status = codes.Unauthenticated

		// 500 - http.StatusInternalServerError
		case svc.ErrUnavailable:
			apiError.Status = codes.Unavailable

		// 500 - http.StatusInternalServerError
		default:
			apiError.Status = codes.Internal
		}

		return apiError, true
	}

	return APIError{}, false
}

func (e APIError) Error() string {
	return e.Message
}
