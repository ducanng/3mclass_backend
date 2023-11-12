package userservice

type (
	UserRegistrationRequest struct {
		Email       string
		PhoneNumber string
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

	UserProfile struct {
		UserID      uint64
		FirstName   string
		LastName    string
		Email       string
		PhoneNumber string
		DisplayName string
	}

	UpdateUserProfileRequest struct {
		Email       string
		FirstName   string
		LastName    string
		PhoneNumber string
	}
)
