package httperror

import (
	"errors"
	"net/http"

	"github.com/thetherington/jobber-common/error-handling/svc"
)

type APIError struct {
	Status  int
	Message string
}

func FromError(err error) (APIError, bool) {
	var apiError APIError
	var svcError svc.Error

	if errors.As(err, &svcError) {
		apiError.Message = svcError.AppError().Error()
		svcError := svcError.SvcError()

		switch svcError {

		// codes.FailedPrecondition / codes.InvalidArgument
		case svc.ErrBadRequest:
			apiError.Status = http.StatusBadRequest

		// codes.NotFound
		case svc.ErrNotFound:
			apiError.Status = http.StatusNotFound

		// codes.Unauthenticated
		case svc.ErrUnAuthorized:
			apiError.Status = http.StatusUnauthorized

		// codes.Internal
		case svc.ErrInternalFailure:
			apiError.Status = http.StatusInternalServerError

		// codes.NotFound
		case svc.ErrNotFound:
			apiError.Status = http.StatusInternalServerError

		default:
			apiError.Status = http.StatusInternalServerError
		}

		return apiError, true
	}

	return APIError{}, false
}

func (e APIError) Error() string {
	return e.Message
}
