package userhandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"

	"github.com/ducanng/3mclass_backend/config"
	"github.com/ducanng/3mclass_backend/handler"
	"github.com/ducanng/3mclass_backend/helper"
	"github.com/ducanng/3mclass_backend/internal/service/userservice"
	"github.com/ducanng/3mclass_backend/pkg/httputil"
	"github.com/ducanng/3mclass_backend/pkg/logutil"
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
	r.Get("/profile", h.userProfile)
}

// logoutUser godoc
//
//	@Summary		Logout user
//	@Description	Logout user
//	@Tags			Public/User
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

// updateProfile godoc
//
//	@Summary		Update user profile
//	@Description	Update user profile
//	@Tags			Public/User
//	@Accept			json
//	@Produce		json
//	@Param			UserProfile	body		UpdateUserProfileRequest	true	"Request"
//	@Success		200			{string}	string						"Update profile ok!"
//	@Router			/v1/public/u/user/update_profile [post]
func (h *userHandler) updateProfile(w http.ResponseWriter, r *http.Request) {
	var (
		_, claims, _ = jwtauth.FromContext(r.Context())
		logger       = logutil.GetLogger()
		req          UpdateUserProfileRequest
		userId       = uint64(claims["user_id"].(float64))
	)
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Errorf("Failed to decode request body, err=%s", err.Error())
		httputil.WriteJSONMessage(w, http.StatusBadRequest, err.Error())
		return
	}
	err = h.validate.Struct(&req)
	if err != nil {
		logger.Errorf("Failed to validate request, err=%s", err.Error())
		httputil.WriteJSONMessage(w, http.StatusBadRequest, err.Error())
		return
	}
	err = h.UserService.UpdateUserProfile(r.Context(), userId, &userservice.UpdateUserProfileRequest{
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		logger.Errorf("Failed to update user profile, err=%s", err.Error())
		httputil.WriteJSONMessage(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.WriteJSONMessage(w, http.StatusOK, "Update profile ok!")
}

// userProfile godoc
//
//	@Summary		Get user profile
//	@Description	Get user profile
//	@Tags			Public/User
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	UserProfile
//	@Router			/v1/public/u/user/profile [get]
func (h *userHandler) userProfile(w http.ResponseWriter, r *http.Request) {
	var (
		_, claims, _ = jwtauth.FromContext(r.Context())
		logger       = logutil.GetLogger()
		userID       = uint64(claims["user_id"].(float64))
	)
	logger.Info("claims", claims)
	userProfile, err := h.UserService.GetUserInfo(r.Context(), userID)
	if err != nil {
		logger.Errorf("error while getting user profile, err: %s", err.Error())
		httputil.WriteJSONMessage(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.WriteJSONResponse(w, UserProfile{
		UserID:      userProfile.UserID,
		Email:       userProfile.Email,
		PhoneNumber: userProfile.PhoneNumber,
		FirstName:   userProfile.FirstName,
		LastName:    userProfile.LastName,
		DisplayName: userProfile.DisplayName,
	})
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
