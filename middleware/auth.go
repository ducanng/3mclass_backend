package middleware

import (
	"context"
	"net/http"

	"github.com/ducanng/no-name/config"
	"github.com/ducanng/no-name/helper"
)

type (
	AuthMiddleware interface {
		Handle(h http.Handler) http.Handler
	}
	authMiddleware struct {
		cfg       *config.Config
		jwtHelper helper.JWTHelper
	}
	authHeaderInfo struct {
		// Only 1 of below fields should be present
		userID     uint64
		actorEmail string
	}
)

func NewAuthMiddleware(cfg *config.Config, jwtHelper helper.JWTHelper) AuthMiddleware {
	return &authMiddleware{
		cfg:       cfg,
		jwtHelper: jwtHelper,
	}
}

func (a *authMiddleware) Handle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		claims, err := a.jwtHelper.VerifyToken(r.Context(), r.Header.Get("Authorization"))
		if err != nil {
			// Invalid access token
			deleteAccessToken(w, a.cfg.AccessTokenCookie.Domain, a.cfg.AccessTokenCookie.PreviousDomain)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Kiểm tra quyền truy cập nếu cần
		// ...

		ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func deleteAccessToken(w http.ResponseWriter, domain, previousDomain string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1,
		Domain:   domain,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "authenticated",
		Value:    "0",
		HttpOnly: false,
		Path:     "/",
		Domain:   domain,
		MaxAge:   -1,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	if previousDomain != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			MaxAge:   -1,
			Domain:   previousDomain,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "authenticated",
			Value:    "0",
			HttpOnly: false,
			Path:     "/",
			Domain:   previousDomain,
			MaxAge:   -1,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "authenticated",
		Value:    "0",
		HttpOnly: false,
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}
