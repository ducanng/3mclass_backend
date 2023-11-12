package userservice

type (
	UserRegistrationRequest struct {
		Email       string
		PhoneNumber string
		DisplayName string
		FirstName   string
		LastName    string
		Password    string
		RePassword  string
		// PhoneCountryCode   string
		// CountryCallingCode string
	}

	UserRegistrationResponse struct {
		NextAction  string
		Session     string
		PhoneNumber string
		Email       string
		// WaitingResendOTPSeconds uint64
		Message string
	}

	UserLogin struct {
		Email    string
		Password string
	}
	UserInfo struct {
		UserName string
		UserID   uint64
	}
)
