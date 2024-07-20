package http

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	token "github.com/thetherington/jobber-common/client-token"
	"github.com/thetherington/jobber-common/models/auth"
)

// /api/gateway/v1/auth
func (ah AuthHandler) CurrentUserRoutes(router chi.Router) {
	router.Get("/current-user", ah.CurrentUser)
	router.Get("/refresh-token/{username}", ah.RefreshToken)
	router.Get("/logged-in-user", ah.GetLoggedInUsers)

	router.Put("/change-password", ah.ChangePassword)
	router.Post("/resend-email", ah.ResendEmail)

	router.Delete("/logged-in-user/{username}", ah.RemoveLoggedInUser)
}

func (ah *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var changePassword *auth.ChangePassword

	// unmarshal the request body into the verifyEmail
	if err := ReadJSON(w, r, &changePassword); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// send a request to the auth microservice via grpc client
	resp, err := ah.svc.ChangePassword(r.Context(), changePassword)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("ChangePassword: failed to write http response")
	}
}

func (ah *AuthHandler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	// send a request to the auth microservice via grpc client
	resp, err := ah.svc.CurrentUser(r.Context())
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("CurrentUser: failed to write http response")
	}
}

func (ah *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// from sesssion cast currentUser to get username information
	payload, ok := ah.session.Get(r.Context(), "currentUser").(*token.Payload)
	if !ok {
		ErrorJSON(w, errors.New("token is not available. Please login again, GatewayService RefreshToken() method error"), http.StatusUnauthorized)
		return
	}

	// send a request to the auth microservice via grpc client
	resp, err := ah.svc.RefreshToken(r.Context(), &auth.RefreshToken{Username: payload.Username})
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// refresh the cookie expiry and update the stored token
	ah.session.RenewToken(r.Context())
	ah.session.Put(r.Context(), "token", resp.Token)

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("RefreshToken: failed to write http response")
	}
}

func (ah *AuthHandler) GetLoggedInUsers(w http.ResponseWriter, r *http.Request) {
	resp, err := ah.svc.GetLoggedInUsers(r.Context())
	if err != nil {
		ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("GetLoggedInUsers: failed to write http response")
	}
}

func (ah *AuthHandler) ResendEmail(w http.ResponseWriter, r *http.Request) {
	var resendEmail *auth.ResendEmail

	// unmarshal the request body into the verifyEmail
	if err := ReadJSON(w, r, &resendEmail); err != nil {
		ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// send a request to the auth microservice via grpc client
	resp, err := ah.svc.ResendEmail(r.Context(), resendEmail)
	if err != nil {
		ServiceErrorResolve(w, err)
		return
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("ResendEmail: failed to write http response")
	}
}

func (ah *AuthHandler) RemoveLoggedInUser(w http.ResponseWriter, r *http.Request) {
	var username string

	// get the token from the url
	if username = chi.URLParam(r, "username"); username == "" {
		ErrorJSON(w, fmt.Errorf("missing username"), http.StatusBadRequest)
		return
	}

	if err := ah.svc.RemoveLoggedInUser(r.Context(), username); err != nil {
		ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	resp := &struct {
		Message string `json:"message"`
	}{
		Message: "User is offline",
	}

	// response back to the front-end application
	if err := WriteJSON(w, http.StatusOK, &resp); err != nil {
		slog.With("error", err).Error("RemoveLoggedInUser: failed to write http response")
	}
}
