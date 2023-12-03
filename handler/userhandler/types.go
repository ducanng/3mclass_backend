package userhandler

type (
	UserProfile struct {
		UserID      uint64 `json:"userID"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phoneNumber"`
		FirstName   string `json:"firstName"`
		LastName    string `json:"lastLame"`
	}

	UpdateUserProfileRequest struct {
		Email       string `json:"email"`
		FirstName   string `json:"firstName"`
		LastName    string `json:"lastName"`
		PhoneNumber string `json:"phoneNumber"`
	}
)
