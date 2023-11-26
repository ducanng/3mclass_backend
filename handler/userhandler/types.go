package userhandler

type (
	UserProfile struct {
		UserID      uint64 `json:"user_id"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
	}

	UpdateUserProfileRequest struct {
		Email       string `json:"email"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		PhoneNumber string `json:"phone_number"`
	}
)
