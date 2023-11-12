package userhandler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"

	"github.com/ducanng/no-name/config"
	"github.com/ducanng/no-name/handler"
	"github.com/ducanng/no-name/helper"
	"github.com/ducanng/no-name/internal/service/userservice"
	"github.com/ducanng/no-name/pkg/httputil"
)

type userHandler struct {
	userservice.UserService
	cfg       *config.Config
	validate  *validator.Validate
	jwtHelper helper.JWTHelper
}

func NewUserHandler(cfg *config.Config, service userservice.UserService, jwtHelper helper.JWTHelper) handler.Handler {
	return &userHandler{
		UserService: service,
		cfg:         cfg,
		validate:    validator.New(),
		jwtHelper:   jwtHelper,
	}
}

func (h *userHandler) Register(r chi.Router) {
	r.Post("/logout", h.logoutUser)
	r.Post("/change_password", h.changePassword)
	r.Post("/update_profile", h.updateProfile)
	r.Get("/profile", h.updateProfile)
}

// logoutUser godoc
//
//	@Summary		Logout user
//	@Description	Logout user
//	@Tags			Public/Auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	string	"Logout ok!"
//	@Router			/v1/public/u/user/logout [post]
func (h *userHandler) logoutUser(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	// handle logout
	var domain = h.cfg.AccessTokenCookie.Domain
	// TODO: check domain
	deleteAccessToken(w, []string{h.cfg.AccessTokenCookie.CookieName}, domain, false)
	httputil.WriteJSONMessage(w, http.StatusOK, fmt.Sprintf("logout OK! %v", claims["user_id"]))
}

func (h *userHandler) changePassword(w http.ResponseWriter, r *http.Request) {

}

func (h *userHandler) updateProfile(w http.ResponseWriter, r *http.Request) {

}

func deleteAccessToken(w http.ResponseWriter, cookieNames []string, domain string, skipAuthenticated bool) {
	for _, cookieName := range cookieNames {
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			MaxAge:   -1,
			Domain:   domain,
		})
	}
	if skipAuthenticated {
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "authenticated",
		Value:    "0",
		HttpOnly: false,
		Path:     "/",
		Domain:   domain,
		MaxAge:   -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "authenticated",
		Value:    "0",
		HttpOnly: false,
		Path:     "/",
		MaxAge:   -1,
	})
}
