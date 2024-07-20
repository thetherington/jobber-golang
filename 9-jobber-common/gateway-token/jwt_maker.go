package token

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

type JWTCustomClaims struct {
	jwt.RegisteredClaims
}

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
func NewGatewayJWTMaker(secretKey string) (TokenMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{
		secretKey: secretKey,
	}, nil
}

// CreateToken creates a new token for specific username and duration
func (maker *JWTMaker) CreateToken(id string) (string, *Payload, error) {
	payload, err := NewPayload(id)
	if err != nil {
		return "", payload, err
	}

	// Create claims with multiple fields populated
	claims := JWTCustomClaims{
		jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(payload.IssuedAt),
			ID:       payload.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", payload, fmt.Errorf("failed to sign token in createJWT: %w", err)
	}

	return signedToken, payload, nil

}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &JWTCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}

		return []byte(maker.secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("problem with parsing token: %w", err)
	}

	if claims, ok := jwtToken.Claims.(*JWTCustomClaims); ok && jwtToken.Valid {
		return &Payload{
			ID:       claims.ID,
			IssuedAt: claims.IssuedAt.Local(),
		}, nil
	}

	return nil, fmt.Errorf("problem with token claims: %w", err)
}
