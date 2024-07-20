package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/thetherington/jobber-common/models/auth"
	"github.com/thetherington/jobber-gateway/internal/core/port"
)

// AuthHandler represents the HTTP handler for auth service requests
type AuthHandler struct {
	svc     port.AuthService
	session *scs.SessionManager
}

// /api/gateway/v1/auth
func (ah AuthHandler) Routes(router chi.Router) {
	router.Post("/signup", ah.SignUp)
	router.Post("/signin", ah.SignIn)
	router.Post("/signout", ah.SignOut)

	router.Put("/verify-email", ah.VerifyEmail)
	router.Put("/forgot-password", ah.ForgotPassword)
	router.Put("/reset-password/{token}", ah.ResetPassword)

	router.Put("/seed/{count}", ah.Seed)
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(svc port.AuthService, session *scs.SessionManager) *AuthHandler {
	return &AuthHandler{
		svc,
		session,
	}
}

func (ah *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var signupPayload *auth.SignUpPayload

	// unmarshal the request body into the signup
	if err := ReadJSON(w, r, &signupPayload); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	ah.session.RenewToken(r.Context())

	// send a request to the auth microservice via grpc client
	resp, err := ah.svc.SignUp(r.Context(), signupPayload)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// take the token from the response from auth microservice and set the session cookie token value
	ah.session.Put(r.Context(), "token", resp.Token)

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusCreated, &resp); err != nil {
		slog.With("error", err).Error("Signup: failed to write http response")
	}
}

func (ah *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var signinPayload *auth.SignInPayload

	// unmarshal the request body into the signup
	if err := ReadJSON(w, r, &signinPayload); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	ah.session.RenewToken(r.Context())

	// send a request to the auth microservice via grpc client
	resp, err := ah.svc.SignIn(r.Context(), signinPayload)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// take the token from the response from auth microservice and set the session cookie token value
	ah.session.Put(r.Context(), "token", resp.Token)

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("Signin: failed to write http response")
	}
}

func (ah *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	// temp
	ah.svc.SignOut(r.Context())

	if err := ah.session.Destroy(r.Context()); err != nil {
		ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "Logout successful",
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, resp); err != nil {
		slog.With("error", err).Error("SignOut: failed to write http response")
	}
}

func (ah *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var verifyEmailPayload *auth.VerifyEmail

	// unmarshal the request body into the verifyEmail
	if err := ReadJSON(w, r, &verifyEmailPayload); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// send a request to the auth microservice via grpc client
	resp, err := ah.svc.VerifyEmail(r.Context(), verifyEmailPayload)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("VerifyEmail: failed to write http response")
	}
}

func (ah *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var forgotPassword *auth.ForgotPassword

	// unmarshal the request body into the verifyEmail
	if err := ReadJSON(w, r, &forgotPassword); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// send a request to the auth microservice via grpc client
	resp, err := ah.svc.ForgotPassword(r.Context(), forgotPassword)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("ForgotPassword: failed to write http response")
	}
}

func (ah *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var resetPassword *auth.ResetPassword

	// unmarshal the request body into the verifyEmail
	if err := ReadJSON(w, r, &resetPassword); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// get the token from the url
	if resetPassword.Token = chi.URLParam(r, "token"); resetPassword.Token == "" {
		ErrorJSON(w, fmt.Errorf("missing token from url"), http.StatusBadRequest)
		return
	}

	// send a request to the auth microservice via grpc client
	resp, err := ah.svc.ResetPassword(r.Context(), resetPassword)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("ResetPassword: failed to write http response")
	}
}

func (ah *AuthHandler) Seed(w http.ResponseWriter, r *http.Request) {
	var (
		param string
		count int
		err   error
	)

	// get the count from the url
	if param = chi.URLParam(r, "count"); param == "" {
		ErrorJSON(w, fmt.Errorf("missing count"), http.StatusBadRequest)
		return
	}

	count, err = strconv.Atoi(param)
	if err != nil {
		ErrorJSON(w, fmt.Errorf("count is not a number"), http.StatusBadRequest)
		return
	}

	resp, err := ah.svc.Seed(r.Context(), count)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("Seed: failed to write http response")
	}
}
