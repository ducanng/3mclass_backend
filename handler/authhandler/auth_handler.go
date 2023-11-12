package authhandler

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"

	"github.com/ducanng/no-name/config"
	"github.com/ducanng/no-name/handler"
	"github.com/ducanng/no-name/helper"
	"github.com/ducanng/no-name/internal/service/userservice"
	"github.com/ducanng/no-name/pkg/httputil"
	"github.com/ducanng/no-name/pkg/logutil"
)

type authHandler struct {
	userservice.UserService
	cfg       *config.Config
	validate  *validator.Validate
	jwtHelper helper.JWTHelper
}

func NewAuthHandler(cfg *config.Config, service userservice.UserService, jwtHelper helper.JWTHelper) handler.Handler {
	return &authHandler{
		UserService: service,
		cfg:         cfg,
		validate:    validator.New(),
		jwtHelper:   jwtHelper,
	}
}

func (a *authHandler) Register(r chi.Router) {
	r.Post("/registration", a.userRegistration)
	r.Post("/login", a.userLogin)
	r.Get("/refresh_token", a.refreshToken)
	r.Post("/verify_otp", a.verifyOTP)
}

// userRegistration godoc
//
//	@Summary		User registration by phone, email, display name
//	@Description	User registration by phone, email, display name
//	@Tags			Public/Auth
//	@Accept			json
//	@Produce		json
//	@Param			UserRegistrationRequest	body		UserRegistrationRequest	true	"Request"
//	@Success		200						{object}	UserRegistrationResponse
//	@Router			/v1/public/auth/registration [post]
func (a *authHandler) userRegistration(w http.ResponseWriter, r *http.Request) {
	var (
		logger = logutil.GetLogger()
		req    UserRegistrationRequest
		// phoneNumber *phonenumbers.PhoneNumber
	)
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Errorf("Failed to decode request body, err=%s", err.Error())
		httputil.WriteJSONMessage(w, http.StatusBadRequest, err.Error())
		return
	}
	err = a.validate.Struct(&req)
	if err != nil {
		logger.Errorf("Failed to validate request, err=%s", err.Error())
		httputil.WriteJSONMessage(w, http.StatusBadRequest, err.Error())
		return
	}
	if _, err = mail.ParseAddress(req.Email); err != nil {
		logger.Errorf("Invalid input email %s, err:%s ", req.Email, err)
		httputil.WriteJSONMessage(w, http.StatusBadRequest, err.Error())
	}
	//if phoneNumber, err = phonenumbers.Parse(req.PhoneNumber, req.PhoneCountryCode); err != nil {
	//	logger.Errorf("Invalid input phone number %s, calling call: %s, country code: %s, err:%s ", req.PhoneNumber, req.CountryCallingCode, req.PhoneCountryCode, err)
	//	httputil.WriteJSONMessage(w, http.StatusBadRequest, err.Error())
	//	return
	//}
	//parsedPhoneNumber := fmt.Sprintf("%d%d", *phoneNumber.CountryCode, *phoneNumber.NationalNumber)

	respRegistration, err := a.UserService.RegisterUser(r.Context(), &userservice.UserRegistrationRequest{
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		DisplayName: req.DisplayName,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Password:    req.Password,
		RePassword:  req.RePassword,
	})
	if err != nil {
		logger.Errorf("Failed to register user, err=%s", err.Error())
		httputil.WriteJSONMessage(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := UserRegistrationResponse{
		NextAction:  respRegistration.NextAction,
		Session:     respRegistration.Session,
		PhoneNumber: respRegistration.PhoneNumber,
		Email:       respRegistration.Email,
		Message:     respRegistration.Message,
	}
	httputil.WriteJSONResponse(w, resp)
}

// userLogin godoc
//
//	@Summary		User login by phone or email
//	@Description	User login by phone or email
//	@Tags			Public/Auth
//	@Accept			json
//	@Produce		json
//	@Param			UserLoginRequest	body		UserLoginRequest	true	"Request"
//	@Success		200					{object}	UserLoginResponse
//	@Router			/v1/public/auth/login [post]
func (a *authHandler) userLogin(w http.ResponseWriter, r *http.Request) {
	var (
		loginReq UserLoginRequest
	)
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		deleteAccessToken(w, "")
		httputil.WriteJSONMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	err = a.validate.Struct(loginReq)
	if err != nil {
		deleteAccessToken(w, "")
		httputil.WriteJSONMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := a.UserService.UserLogin(r.Context(), userservice.UserLogin{
		Email:    strings.ToLower(strings.TrimSpace(loginReq.Email)),
		Password: loginReq.Password,
	})
	if err != nil {
		deleteAccessToken(w, "")
		httputil.WriteJSONMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	// write to cookie
	token, err := a.jwtHelper.IssueToken(r.Context(), map[string]interface{}{
		"user_id": userInfo.UserID,
	})
	if err != nil {
		deleteAccessToken(w, "")
		httputil.WriteJSONMessage(w, http.StatusBadRequest, err.Error())
		return
	}
	var (
		expireAt = time.Now().Add(time.Hour * 24 * 365)
	)
	if a.cfg.JWT.ExpiryTime > 0 {
		expireAt = time.Now().Add(time.Second * time.Duration(a.cfg.JWT.ExpiryTime))
	}
	var domain = a.cfg.AccessTokenCookie.Domain
	w.Header().Set("Authorization", "BEARER"+token)
	http.SetCookie(w, &http.Cookie{
		Name:     a.cfg.AccessTokenCookie.CookieName,
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Expires:  expireAt,
		Domain:   domain,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "authenticated",
		Value:    "1",
		HttpOnly: false,
		Path:     "/",
		Expires:  expireAt,
		Domain:   domain,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	if domain != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "authenticated",
			Value:    "1",
			HttpOnly: false,
			Path:     "/",
			MaxAge:   -1,
		})
	}
	httputil.WriteJSONMessage(w, http.StatusOK, "Ok!")
}
func deleteAccessToken(w http.ResponseWriter, domain string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1,
		Domain:   domain,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "authenticated",
		Value:    "0",
		HttpOnly: false,
		Path:     "/",
		MaxAge:   -1,
		Domain:   domain,
	})
}

// refreshToken godoc
//
//	@Summary		Refresh Token
//	@Description	Refresh Token
//	@Tags			Public/Auth
//	@Accept			json
//	@Produce		json
//	@Param			RefreshTokenRequest	query		RefreshTokenRequest	true	"Request"
//	@Success		200					{object}	RefreshTokenResponse
//	@Router			/v1/public/auth/refresh_token [get]
func (a *authHandler) refreshToken(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSONMessage(w, http.StatusNotImplemented, "Ok!")
}

// verifyOTP godoc
//
//	@Summary		Verify OTP
//	@Description	Verify OTP
//	@Tags			Public/Auth
//	@Accept			json
//	@Produce		json
//	@Param			VerifyOTPRequest	body		VerifyOTPRequest	true	"Request"
//	@Success		200					{object}	VerifyOTPResponse
//	@Router			/v1/public/auth/verify_otp [post]
func (a *authHandler) verifyOTP(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSONMessage(w, http.StatusNotImplemented, "Ok!")
}
