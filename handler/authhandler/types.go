package authhandler

type (
	UserRegistrationRequest struct {
		Password    string `json:"password" validate:"required"`
		Email       string `json:"email" validate:"required,email"`
		PhoneNumber string `json:"phoneNumber" validate:"min=9,max=12,required,numeric"`
		FirstName   string `json:"firstName" validate:"required"`
		LastName    string `json:"lastName" validate:"required"`
		RePassword  string `json:"rePassword" validate:"required"`

		//PhoneCountryCode   string `json:"phone_country_code" validate:"required"`
		//CountryCallingCode string `json:"country_calling_code" validate:"required"`
	}
	UserRegistrationResponse struct {
		NextAction  string `json:"nextAction"`
		Session     string `json:"session"`
		PhoneNumber string `json:"phoneNumber"`
		Email       string `json:"email"`
		// WaitingResendOTPSeconds uint64 `json:"waiting_resend_otp_seconds"`
		Message string `json:"message"`
	}
	UserLoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	RefreshTokenRequest struct {
	}

	VerifyOTPRequest struct {
	}
)
