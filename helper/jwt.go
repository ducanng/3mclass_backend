package helper

import (
	"context"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

type (
	JWTClaims map[string]interface{}

	JWTHelper interface {
		IssueToken(ctx context.Context, claims JWTClaims) (string, error)
		VerifyToken(ctx context.Context, token string) (JWTClaims, error)
	}
	jwtHelper struct {
		jwtauth.JWTAuth
		expiryTime time.Duration
	}
)

func NewJWTHelper(jwtAuth jwtauth.JWTAuth, expiryTime time.Duration) JWTHelper {
	return &jwtHelper{
		JWTAuth:    jwtAuth,
		expiryTime: expiryTime,
	}
}

func (j *jwtHelper) IssueToken(ctx context.Context, claims JWTClaims) (string, error) {
	tm := time.Now().Add(time.Hour * 24 * 365)
	if j.expiryTime > 0 {
		tm = time.Now().Add(time.Second * j.expiryTime)
	}
	// Encode and sign the token
	jwtauth.SetExpiry(claims, tm)
	jwtauth.SetIssuedNow(claims)
	_, tokenString, err := j.Encode(claims)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *jwtHelper) VerifyToken(ctx context.Context, token string) (JWTClaims, error) {
	_, err := jwtauth.VerifyToken(&j.JWTAuth, token)
	if err != nil {
		return nil, err
	}
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
