package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/alexedwards/scs/v2"
	token "github.com/thetherington/jobber-common/client-token"
)

const BEARER_SCHEMA = "Bearer "

// MiddlewareHandler represents the HTTP handler for middleware functions
type MiddlewareHandler struct {
	session *scs.SessionManager
	token   token.TokenMaker
}

// NewMiddlewareHandler creates a new MiddlewareHandler instance
func NewMiddlewareHandler(session *scs.SessionManager, token token.TokenMaker) *MiddlewareHandler {
	return &MiddlewareHandler{
		session,
		token,
	}
}

func (handler *MiddlewareHandler) VerifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		// try to get the token from the auth header
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, BEARER_SCHEMA) {
			token = authHeader[len(BEARER_SCHEMA):]
		}

		// if no Authorization header Bearer token, try the session store from the request cookie
		if token == "" {
			// check if token exists in cookie
			token = handler.session.GetString(r.Context(), "token")
			if token == "" {
				ErrorJSON(w, errors.New("token is not available. Please login again, GatewayService verifyUser() method error"), http.StatusUnauthorized)
				return
			}
		}

		// verify token is signed correctly
		payload, err := handler.token.VerifyToken(token)
		if err != nil {
			ErrorJSON(w, errors.New("token is not available. Please login again, GatewayService verifyUser() method error"), http.StatusUnauthorized)
			return
		}

		// put the payload into the session to be accessed via CheckAuthentication & grpc clients metadata insert interceptor
		handler.session.Put(r.Context(), "currentUser", payload)

		next.ServeHTTP(w, r)
	})
}

func (handler *MiddlewareHandler) CheckAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// cast currentUser to validate information is correct
		payload, ok := handler.session.Get(r.Context(), "currentUser").(*token.Payload)
		if !ok {
			ErrorJSON(w, errors.New("token is not available. Please login again, GatewayService CheckAuthentication() method error"), http.StatusBadRequest)
			return
		}

		// validate token expiry (probably redundant)
		if err := payload.Valid(); err != nil {
			ErrorJSON(w, errors.New("token is expired. Please login again, GatewayService CheckAuthentication() valid() error"), http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
