package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/thetherington/jobber-common/error-handling/httperror"
)

type ErrorJSONResponse struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1024 * 1024 * 20 // 10 megabyte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{}) // decode a second object if it exists
	if err != io.EOF {            // err should be io.EOF if only one object exists
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload ErrorJSONResponse

	payload.Status = "error"
	payload.Message = err.Error()
	payload.StatusCode = statusCode

	return WriteJSON(w, statusCode, payload)
}

// Parses error into httperror and creates ErrorJSON response
func ServiceErrorResolve(w http.ResponseWriter, err error) {
	// try to cast the error to a httperror lookup
	if apiError, ok := httperror.FromError(err); ok {
		ErrorJSON(w, apiError, apiError.Status)
		return
	}

	// generic response
	ErrorJSON(w, err, http.StatusInternalServerError)
}
