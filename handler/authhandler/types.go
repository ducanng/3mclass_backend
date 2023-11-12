package authhandler

type (
	UserRegistrationRequest struct {
		Password    string `json:"password" validate:"required"`
		Email       string `json:"email" validate:"required,email"`
		PhoneNumber string `json:"phone_number" validate:"min=9,max=12,required,numeric"`
		FirstName   string `json:"first_name" validate:"required"`
		LastName    string `json:"last_name" validate:"required"`
		RePassword  string `json:"re_password" validate:"required"`

		//PhoneCountryCode   string `json:"phone_country_code" validate:"required"`
		//CountryCallingCode string `json:"country_calling_code" validate:"required"`
	}
	UserRegistrationResponse struct {
		NextAction  string `json:"next_action"`
		Session     string `json:"session"`
		PhoneNumber string `json:"phone_number"`
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
